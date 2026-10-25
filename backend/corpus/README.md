# Agent corpus

Everything the portfolio agent knows about you lives in this folder. The
backend loads it at startup (set `CORPUS_DIR` to use a different folder) and
the agent retrieves from it with TF-IDF ranking. This README is ignored.

## Files

| File | Used by | Purpose |
| --- | --- | --- |
| `profile.json` | `get_links`, `get_skills`, `list_projects`, `get_experience_summary` tools, system prompt | Structured facts: name, headline, links, skills, projects, roles |
| `*.md`, `*.markdown`, `*.txt` | `search_profile` tool (RAG) | Free-form knowledge: bio, experience, case studies, talks, FAQs |

Subfolders are read recursively.

## Writing documents

- Each `## Heading` section becomes one retrievable chunk. Keep sections
  focused (one role, one project, one topic) and roughly 50–250 words.
- Text before the first `##` heading is its own chunk, titled by the `# H1`.
- Add an optional `Tags: a, b, c` line inside a section for extra search
  keywords (synonyms, acronyms). Tags are indexed but not shown to the model.
- Write facts in plain sentences. The agent is told never to invent facts, so
  anything missing here is answered with "I don't know".
- Do not put secrets or private data here; anything in the corpus can be
  quoted to any visitor.

Restart the backend after editing (with Docker Compose the folder is mounted,
so `docker compose restart backend` is enough).
