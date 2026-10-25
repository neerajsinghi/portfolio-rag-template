// Package tools defines the tools the agent can call and executes them
// against the corpus (RAG), the structured profile, and the GitHub API.
package tools

import (
	"fmt"
	"sort"
	"strings"

	"github.com/your-username/portfolio-rag-template/backend/internal/github"
	"github.com/your-username/portfolio-rag-template/backend/internal/kb"
	"github.com/your-username/portfolio-rag-template/backend/internal/profile"
)

// Tool is the JSON-serializable tool definition sent to the model.
// This is the same shape an MCP server advertises — name, description, schema.
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"input_schema"`
}

func obj(props map[string]any, required ...string) map[string]any {
	m := map[string]any{"type": "object", "properties": props}
	if len(required) > 0 {
		m["required"] = required
	}
	return m
}

// Definitions returns the tool registry, with descriptions personalised to the owner.
func Definitions(owner string) []Tool {
	return []Tool{
		{
			"search_profile",
			"Full-text search over " + owner + "'s knowledge base (résumé, bio, roles, skills, projects, and any other documents in the corpus). Use this first for any factual question about their background, experience, or fit for a role. Returns the top matching passages.",
			obj(map[string]any{
				"query": map[string]any{"type": "string", "description": "natural-language search query, e.g. 'cloud infrastructure experience'"},
				"top_n": map[string]any{"type": "integer", "description": "number of passages to return (1-5, default 3)", "default": 3},
			}, "query"),
		},
		{
			"get_experience_summary",
			"Return a chronological list of " + owner + "'s work roles (title, company, dates, one-line summary). Use for career overview or timeline questions.",
			obj(map[string]any{}),
		},
		{
			"get_skills",
			"Return " + owner + "'s technical skills grouped by category. Use when asked about tech stack or whether they know a technology.",
			obj(map[string]any{}),
		},
		{
			"list_projects",
			"List " + owner + "'s featured projects with descriptions, tech stacks, and links.",
			obj(map[string]any{}),
		},
		{
			"get_links",
			"Return " + owner + "'s contact details, social links, location, and availability. Use when asked how to reach them or what they are looking for.",
			obj(map[string]any{}),
		},
		{
			"get_github_repos",
			"Fetch " + owner + "'s public GitHub repositories live (name, description, stars, language). Use when asked about open-source work.",
			obj(map[string]any{}),
		},
	}
}

// Executor runs tools against the loaded corpus and profile.
type Executor struct {
	index   *kb.Index
	profile profile.Profile
}

// NewExecutor returns an Executor. index may be nil (empty corpus).
func NewExecutor(index *kb.Index, p profile.Profile) *Executor {
	return &Executor{index: index, profile: p}
}

// Execute runs a tool and returns (resultText, sourceIDs).
func (e *Executor) Execute(name string, input map[string]any) (string, []string) {
	p := e.profile
	switch name {
	case "search_profile":
		q, _ := input["query"].(string)
		k := 3
		if n, ok := input["top_n"].(float64); ok && n >= 1 && n <= 5 {
			k = int(n)
		}
		docs := e.index.Retrieve(q, k)
		if len(docs) == 0 {
			return "No matching passage found in the knowledge base.", nil
		}
		var b strings.Builder
		srcs := make([]string, 0, len(docs))
		for i, d := range docs {
			if i > 0 {
				b.WriteString("\n\n")
			}
			b.WriteString("## " + d.Title + "\n" + d.Text)
			srcs = append(srcs, d.ID)
		}
		return b.String(), srcs

	case "get_experience_summary":
		if len(p.Experience) == 0 {
			return "No work history is configured in profile.json.", nil
		}
		lines := make([]string, 0, len(p.Experience))
		for _, r := range p.Experience {
			lines = append(lines, fmt.Sprintf("• %s | %s | %s\n  %s", r.Title, r.Company, r.Period, r.Summary))
		}
		return strings.Join(lines, "\n\n"), []string{"profile:experience"}

	case "get_skills":
		if len(p.Skills) == 0 {
			return "No skills are configured in profile.json.", nil
		}
		cats := make([]string, 0, len(p.Skills))
		for c := range p.Skills {
			cats = append(cats, c)
		}
		sort.Strings(cats)
		lines := make([]string, 0, len(cats))
		for _, c := range cats {
			lines = append(lines, c+": "+strings.Join(p.Skills[c], ", "))
		}
		return strings.Join(lines, "\n"), []string{"profile:skills"}

	case "list_projects":
		if len(p.Projects) == 0 {
			return "No projects are configured in profile.json.", nil
		}
		lines := make([]string, 0, len(p.Projects))
		for _, pr := range p.Projects {
			line := fmt.Sprintf("• %s\n  %s\n  Stack: %s", pr.Name, pr.Description, strings.Join(pr.Stack, ", "))
			if pr.URL != "" {
				line += "\n  URL: " + pr.URL
			}
			lines = append(lines, line)
		}
		return strings.Join(lines, "\n\n"), []string{"profile:projects"}

	case "get_links":
		var b strings.Builder
		b.WriteString("Contact & availability:\n")
		if p.Email != "" {
			b.WriteString("Email: " + p.Email + "\n")
		}
		for _, l := range p.Links {
			b.WriteString(l.Label + ": " + l.URL + "\n")
		}
		if p.Location != "" {
			b.WriteString("Location: " + p.Location + "\n")
		}
		if p.Availability != "" {
			b.WriteString("Availability: " + p.Availability + "\n")
		}
		return strings.TrimSpace(b.String()), []string{"profile:links"}

	case "get_github_repos":
		repos, err := github.FetchRepos()
		if err != nil || len(repos) == 0 {
			return "No public GitHub repositories were returned (GITHUB_USER may be unset, the account may be empty, or the API is rate-limited).", []string{"github"}
		}
		lines := make([]string, 0, 8)
		for i, r := range repos {
			if i >= 8 {
				break
			}
			line := "• " + r.Name
			if r.Language != "" {
				line += " (" + r.Language + ")"
			}
			if r.Stars > 0 {
				line += fmt.Sprintf(" ★%d", r.Stars)
			}
			if r.Description != "" {
				line += " — " + r.Description
			}
			line += " " + r.HTMLURL
			lines = append(lines, line)
		}
		return strings.Join(lines, "\n"), []string{"github"}
	}
	return "Unknown tool.", nil
}
