// Package server exposes the agent over HTTP: /api/v1/chat (SSE),
// /api/v1/github and /api/v1/health.
package server

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/your-username/portfolio-rag-template/backend/internal/agent"
	"github.com/your-username/portfolio-rag-template/backend/internal/github"
	"github.com/your-username/portfolio-rag-template/backend/internal/kb"
	"github.com/your-username/portfolio-rag-template/backend/internal/llm"
	"github.com/your-username/portfolio-rag-template/backend/internal/profile"
	"github.com/your-username/portfolio-rag-template/backend/internal/tools"
)

// currentAPIVersion is the canonical API version. Routes are mounted under
// /api/{currentAPIVersion}/... and aliased at the unversioned /api/... paths.
const currentAPIVersion = "v1"

// maxChatBody caps the chat request size to protect the LLM budget.
const maxChatBody = 64 << 10

// NewHandler builds the backend HTTP handler with CORS middleware.
func NewHandler(prov llm.Provider, index *kb.Index, p profile.Profile) http.Handler {
	executor := tools.NewExecutor(index, p)
	ag := agent.New(prov, executor, agent.SystemPrompt(p.Name, p.Headline), tools.Definitions(p.Name))

	healthHandler := func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"ok": true, "model": prov.ModelName(), "api_version": currentAPIVersion,
			"owner": p.Name, "corpus_chunks": index.Len(),
		})
	}

	mux := http.NewServeMux()
	registerVersioned(mux, "/api/health", healthHandler)
	registerVersioned(mux, "/api/chat", chatHandler(ag))
	registerVersioned(mux, "/api/github", reposHandler)
	return withCORS(withAPIVersionHeader(mux))
}

func registerVersioned(mux *http.ServeMux, path string, handler http.HandlerFunc) {
	mux.HandleFunc(path, handler)
	mux.HandleFunc("/api/"+currentAPIVersion+strings.TrimPrefix(path, "/api"), handler)
}

func withAPIVersionHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-API-Version", currentAPIVersion)
		next.ServeHTTP(w, r)
	})
}

// withCORS allows configured comma-separated frontend origins (ALLOWED_ORIGIN).
func withCORS(next http.Handler) http.Handler {
	configured := strings.Split(os.Getenv("ALLOWED_ORIGIN"), ",")
	allowed := make(map[string]bool, len(configured))
	for _, origin := range configured {
		if origin = strings.TrimSpace(origin); origin != "" {
			allowed[origin] = true
		}
	}
	if len(allowed) == 0 {
		allowed["*"] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowed["*"] {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "600")
		if r.Method == http.MethodOptions {
			if !allowed["*"] && !allowed[origin] {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type chatRequest struct {
	Messages []agent.Turn `json:"messages"`
}

// chatHandler streams the agent run as Server-Sent Events.
func chatHandler(ag *agent.Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req chatRequest
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxChatBody)).Decode(&req); err != nil || len(req.Messages) == 0 {
			http.Error(w, "expected JSON {messages:[{role,content}]}", http.StatusBadRequest)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		emit := func(event string, data any) {
			payload, _ := json.Marshal(data)
			_, _ = w.Write([]byte("event: " + event + "\ndata: "))
			_, _ = w.Write(payload)
			_, _ = w.Write([]byte("\n\n"))
			flusher.Flush()
		}

		if err := ag.Run(r.Context(), req.Messages, emit); err != nil {
			emit("error", map[string]string{"message": err.Error()})
		}
	}
}

// reposHandler returns the cached public repo list for the live "from GitHub" strip.
func reposHandler(w http.ResponseWriter, _ *http.Request) {
	repos, err := github.FetchRepos()
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"user": github.User(), "repos": []github.Repo{}, "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": github.User(), "repos": repos})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// LoadDotEnv loads KEY=VALUE lines from a .env file if present, without
// overriding variables already set in the environment.
func LoadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k, v = strings.TrimSpace(k), strings.Trim(strings.TrimSpace(v), `"'`)
		if _, exists := os.LookupEnv(k); !exists {
			_ = os.Setenv(k, v)
		}
	}
}
