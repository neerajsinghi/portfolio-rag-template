package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/your-username/portfolio-rag-template/backend/internal/kb"
	"github.com/your-username/portfolio-rag-template/backend/internal/llm"
	anthropicprovider "github.com/your-username/portfolio-rag-template/backend/internal/llm/anthropic"
	geminiprovider "github.com/your-username/portfolio-rag-template/backend/internal/llm/gemini"
	grokprovider "github.com/your-username/portfolio-rag-template/backend/internal/llm/grok"
	openaiprovider "github.com/your-username/portfolio-rag-template/backend/internal/llm/openai"
	"github.com/your-username/portfolio-rag-template/backend/internal/profile"
	"github.com/your-username/portfolio-rag-template/backend/internal/server"
)

// newProvider picks the LLM backend from LLM_PROVIDER (default: anthropic).
func newProvider() llm.Provider {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LLM_PROVIDER"))) {
	case "openai":
		return openaiprovider.New()
	case "grok", "xai":
		return grokprovider.New()
	case "gemini", "google":
		return geminiprovider.New()
	case "", "anthropic", "claude":
		return anthropicprovider.New()
	default:
		log.Fatalf("unknown LLM_PROVIDER %q (want anthropic, openai, grok or gemini)", os.Getenv("LLM_PROVIDER"))
		return nil
	}
}

func main() {
	server.LoadDotEnv(".env")

	corpusDir := os.Getenv("CORPUS_DIR")
	if corpusDir == "" {
		corpusDir = "corpus"
	}
	index, err := kb.Load(corpusDir)
	if err != nil {
		log.Printf("warning: %v — the agent will run with an empty knowledge base", err)
	}
	p, err := profile.Load(filepath.Join(corpusDir, "profile.json"))
	if err != nil {
		log.Printf("warning: %v — using placeholder profile", err)
	}

	prov := newProvider()
	h := server.NewHandler(prov, index, p)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("portfolio agent for %q listening on :%s (model=%s, corpus=%s, chunks=%d)",
		p.Name, port, prov.ModelName(), corpusDir, index.Len())
	log.Fatal(http.ListenAndServe(":"+port, h))
}
