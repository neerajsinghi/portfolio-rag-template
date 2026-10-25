// Package kb loads the agent's knowledge base from a folder of Markdown/text
// documents (the "corpus") and ranks chunks with TF-IDF for retrieval (RAG).
//
// Corpus format (see backend/corpus/README.md):
//
//   - Every *.md / *.markdown / *.txt file in the folder (recursively) is read.
//     README.md files are skipped so the folder can document itself.
//   - Each "## Heading" section becomes one retrievable chunk. Text before the
//     first "##" heading becomes its own chunk, titled by the "# H1" heading or
//     the file name.
//   - An optional line "Tags: a, b, c" inside a section adds search keywords
//     that are indexed but not shown to the model.
package kb

import (
	"bufio"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Doc is one retrievable chunk of the owner's profile.
type Doc struct {
	ID    string
	Title string
	Tags  string
	Text  string
}

// Index is an in-memory TF-IDF index over a set of Docs.
type Index struct {
	docs    []Doc
	docFreq map[string]int
}

var wordRe = regexp.MustCompile(`[a-z0-9+#.]+`)

func tokenize(s string) []string { return wordRe.FindAllString(strings.ToLower(s), -1) }

// New builds an index over docs.
func New(docs []Doc) *Index {
	df := map[string]int{}
	for _, d := range docs {
		seen := map[string]bool{}
		for _, t := range tokenize(d.Title + " " + d.Tags + " " + d.Text) {
			if !seen[t] {
				seen[t] = true
				df[t]++
			}
		}
	}
	return &Index{docs: docs, docFreq: df}
}

// Load reads every document in dir and builds an index over its chunks.
func Load(dir string) (*Index, error) {
	var docs []Doc
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".md" && ext != ".markdown" && ext != ".txt" {
			return nil
		}
		if strings.EqualFold(d.Name(), "README.md") {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		docs = append(docs, parse(rel, string(raw))...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("load corpus %q: %w", dir, err)
	}
	return New(docs), nil
}

// parse splits one document into chunks on "## " headings.
func parse(relPath, content string) []Doc {
	base := strings.TrimSuffix(filepath.ToSlash(relPath), filepath.Ext(relPath))
	fileTitle := base
	var out []Doc
	var title, tags string
	var body strings.Builder

	flush := func() {
		text := strings.TrimSpace(body.String())
		if text != "" {
			t := title
			if t == "" {
				t = fileTitle
			}
			id := base
			if title != "" {
				id = base + "#" + slug(title)
			}
			out = append(out, Doc{ID: id, Title: t, Tags: tags, Text: text})
		}
		body.Reset()
		tags = ""
	}

	sc := bufio.NewScanner(strings.NewReader(content))
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "## "):
			flush()
			title = strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))
		case strings.HasPrefix(trimmed, "# ") && title == "" && body.Len() == 0:
			fileTitle = strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
		case strings.HasPrefix(strings.ToLower(trimmed), "tags:"):
			tags = strings.TrimSpace(trimmed[len("tags:"):])
		default:
			body.WriteString(line)
			body.WriteString("\n")
		}
	}
	flush()
	return out
}

func slug(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	dash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			dash = false
		} else if !dash && b.Len() > 0 {
			b.WriteRune('-')
			dash = true
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}

// Len returns the number of indexed chunks.
func (ix *Index) Len() int {
	if ix == nil {
		return 0
	}
	return len(ix.docs)
}

func (ix *Index) score(qTokens []string, d Doc) float64 {
	tf := map[string]int{}
	for _, t := range tokenize(d.Tags + " " + d.Text) {
		tf[t]++
	}
	// Boost title matches — a query matching the title is strongly relevant.
	for _, t := range tokenize(d.Title) {
		tf[t] += 3
	}
	var s float64
	for _, q := range qTokens {
		if c, ok := tf[q]; ok {
			idf := math.Log(1 + float64(len(ix.docs))/float64(max(ix.docFreq[q], 1)))
			s += (1 + math.Log(float64(c))) * idf
		}
	}
	return s
}

// Retrieve returns the top-k most relevant docs for a query (the RAG step).
func (ix *Index) Retrieve(query string, k int) []Doc {
	if ix == nil {
		return nil
	}
	q := tokenize(query)
	type sd struct {
		d Doc
		s float64
	}
	scored := make([]sd, 0, len(ix.docs))
	for _, d := range ix.docs {
		if s := ix.score(q, d); s > 0 {
			scored = append(scored, sd{d, s})
		}
	}
	sort.SliceStable(scored, func(i, j int) bool { return scored[i].s > scored[j].s })
	out := []Doc{}
	for i := 0; i < len(scored) && i < k; i++ {
		out = append(out, scored[i].d)
	}
	return out
}
