package tools

import (
	"strings"
	"testing"

	"github.com/your-username/portfolio-rag-template/backend/internal/kb"
	"github.com/your-username/portfolio-rag-template/backend/internal/profile"
)

func TestSearchProfileReturnsPassagesAndSources(t *testing.T) {
	ix := kb.New([]kb.Doc{
		{ID: "bio", Title: "Bio", Text: "Builds distributed systems in Go."},
		{ID: "hobbies", Title: "Hobbies", Text: "Enjoys hiking."},
	})
	ex := NewExecutor(ix, profile.Profile{Name: "Alex Doe"})
	text, srcs := ex.Execute("search_profile", map[string]any{"query": "distributed systems", "top_n": float64(1)})
	if !strings.Contains(text, "Builds distributed systems") || len(srcs) != 1 || srcs[0] != "bio" {
		t.Fatalf("text=%q srcs=%v", text, srcs)
	}
}

func TestStructuredToolsUseProfile(t *testing.T) {
	p := profile.Profile{
		Name:       "Alex Doe",
		Email:      "alex@example.com",
		Links:      []profile.Link{{Label: "GitHub", URL: "https://github.com/your-username"}},
		Skills:     map[string][]string{"Languages": {"Go", "TypeScript"}},
		Projects:   []profile.Project{{Name: "Demo", Description: "A demo.", Stack: []string{"Go"}}},
		Experience: []profile.Role{{Title: "Engineer", Company: "Example Corp", Period: "2020 – Present", Summary: "Builds things."}},
	}
	ex := NewExecutor(nil, p)
	for tool, want := range map[string]string{
		"get_links":              "alex@example.com",
		"get_skills":             "Languages: Go, TypeScript",
		"list_projects":          "Demo",
		"get_experience_summary": "Example Corp",
	} {
		if got, _ := ex.Execute(tool, map[string]any{}); !strings.Contains(got, want) {
			t.Errorf("%s = %q, want substring %q", tool, got, want)
		}
	}
	if got, _ := ex.Execute("search_profile", map[string]any{"query": "x"}); !strings.Contains(got, "No matching") {
		t.Errorf("search on empty index = %q", got)
	}
}

func TestDefinitionsMentionOwner(t *testing.T) {
	for _, d := range Definitions("Alex Doe") {
		if !strings.Contains(d.Description, "Alex Doe") {
			t.Errorf("%s description does not mention owner", d.Name)
		}
	}
}
