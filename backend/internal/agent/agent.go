// Package agent runs the provider-agnostic tool-use loop.
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/your-username/portfolio-rag-template/backend/internal/llm"
	"github.com/your-username/portfolio-rag-template/backend/internal/tools"
)

// SystemPrompt builds the agent's instructions for a given owner. Edit this
// to change the agent's voice or rules.
func SystemPrompt(owner, headline string) string {
	who := owner
	if headline != "" {
		who += ", " + headline
	}
	return fmt.Sprintf("You are the portfolio agent for %s. "+
		"Answer a visitor's questions about %s — their experience, skills, projects, and fit for roles. "+
		"Rules: (1) ALWAYS ground answers in the tools; call search_profile (or the other tools) before answering "+
		"factual questions, and never invent facts. If the tools return nothing relevant, say you don't know. "+
		"(2) Include clickable full URLs from tool results when relevant. (3) Be concise and specific — 2 to 4 short "+
		"paragraphs max, in a warm, professional voice. (4) Refer to them by name. (5) If asked whether they fit a role, "+
		"search first, then give an honest, evidence-based take. (6) If a question is unrelated to %s or their career, "+
		"briefly say so and steer back. You are speaking to recruiters and engineers, so be substantive, not salesy.",
		who, owner, owner)
}

// Turn is a simple text message in the visitor-facing conversation history.
type Turn struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Agent runs the tool-use loop against an LLM provider.
type Agent struct {
	provider llm.Provider
	executor *tools.Executor
	system   string
	toolDefs []llm.ToolDef
}

// New returns an Agent backed by the given provider and tool executor.
func New(p llm.Provider, executor *tools.Executor, systemPrompt string, defs []tools.Tool) *Agent {
	toolDefs := make([]llm.ToolDef, len(defs))
	for i, t := range defs {
		toolDefs[i] = llm.ToolDef{Name: t.Name, Description: t.Description, InputSchema: t.InputSchema}
	}
	return &Agent{provider: p, executor: executor, system: systemPrompt, toolDefs: toolDefs}
}

// Run executes the tool-use loop, emitting events as it goes:
//
//	tool    {name, input}   – the model decided to call a tool (MCP-style)
//	sources {name, sources} – which corpus chunks the tool returned
//	text    {text}          – an answer paragraph from the model
//	done    {}              – finished
func (a *Agent) Run(ctx context.Context, history []Turn, emit func(event string, data any)) error {
	msgs := make([]llm.Message, 0, len(history)+6)
	for _, t := range history {
		role := t.Role
		if role != "assistant" {
			role = "user"
		}
		msgs = append(msgs, llm.Message{Role: role, Content: t.Content})
	}

	for step := 0; step < 5; step++ {
		resp, err := a.provider.Complete(ctx, llm.Request{
			System:    a.system,
			Messages:  msgs,
			Tools:     a.toolDefs,
			MaxTokens: 1000,
		})
		if err != nil {
			return err
		}

		msgs = append(msgs, llm.Message{Role: "assistant", Content: resp.Content})

		var toolUses []llm.Block
		for _, b := range resp.Content {
			switch b.Type {
			case "text":
				if strings.TrimSpace(b.Text) != "" {
					emit("text", map[string]string{"text": b.Text})
				}
			case "tool_use":
				toolUses = append(toolUses, b)
			}
		}

		if resp.StopReason != llm.StopReasonToolUse || len(toolUses) == 0 {
			break
		}

		results := make([]llm.Block, 0, len(toolUses))
		for _, tu := range toolUses {
			var input map[string]any
			_ = json.Unmarshal(tu.Input, &input)
			emit("tool", map[string]any{"name": tu.Name, "input": input})
			text, srcs := a.executor.Execute(tu.Name, input)
			emit("sources", map[string]any{"name": tu.Name, "sources": srcs})
			results = append(results, llm.Block{Type: "tool_result", ToolUseID: tu.ID, Content: text})
		}
		msgs = append(msgs, llm.Message{Role: "user", Content: results})
	}

	emit("done", map[string]any{})
	return nil
}
