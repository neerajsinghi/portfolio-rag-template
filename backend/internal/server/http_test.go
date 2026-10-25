package server

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/your-username/portfolio-rag-template/backend/internal/kb"
	"github.com/your-username/portfolio-rag-template/backend/internal/llm"
	"github.com/your-username/portfolio-rag-template/backend/internal/profile"
)

// scriptedProvider calls search_profile once, then answers.
type scriptedProvider struct{ calls int }

func (s *scriptedProvider) Complete(_ context.Context, req llm.Request) (llm.Response, error) {
	s.calls++
	if s.calls == 1 {
		return llm.Response{
			StopReason: llm.StopReasonToolUse,
			Content:    []llm.Block{{Type: "tool_use", ID: "t1", Name: "search_profile", Input: json.RawMessage(`{"query":"go"}`)}},
		}, nil
	}
	return llm.Response{StopReason: llm.StopReasonEndTurn, Content: []llm.Block{{Type: "text", Text: "Alex writes Go."}}}, nil
}
func (*scriptedProvider) ModelName() string { return "fake-model" }

func newTestHandler() *scriptedProvider { return &scriptedProvider{} }

func TestVersionedRoutesAliasLegacyRoutes(t *testing.T) {
	handler := NewHandler(newTestHandler(), nil, profile.Profile{Name: "Alex Doe"})
	for _, path := range []string{"/api/health", "/api/v1/health"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		if response.Code != 200 {
			t.Fatalf("%s status = %d, want 200", path, response.Code)
		}
		if version := response.Header().Get("X-API-Version"); version != "v1" {
			t.Fatalf("%s X-API-Version = %q, want v1", path, version)
		}
	}
}

func TestChatStreamsToolTraceAndAnswer(t *testing.T) {
	index := kb.New([]kb.Doc{{ID: "bio", Title: "Bio", Text: "Writes Go services."}})
	handler := NewHandler(newTestHandler(), index, profile.Profile{Name: "Alex Doe"})
	body := strings.NewReader(`{"messages":[{"role":"user","content":"What languages?"}]}`)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("POST", "/api/v1/chat", body))
	out := response.Body.String()
	for _, want := range []string{"event: tool", `"search_profile"`, "event: sources", `"bio"`, "event: text", "Alex writes Go.", "event: done"} {
		if !strings.Contains(out, want) {
			t.Fatalf("stream missing %q:\n%s", want, out)
		}
	}
}

func TestChatRejectsEmptyBody(t *testing.T) {
	handler := NewHandler(newTestHandler(), nil, profile.Profile{})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("POST", "/api/v1/chat", strings.NewReader(`{}`)))
	if response.Code != 400 {
		t.Fatalf("status = %d, want 400", response.Code)
	}
}
