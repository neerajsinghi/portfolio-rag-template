# ── Stage 1: build ──────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# ── Stage 2: runtime ─────────────────────────────────────────────────────────
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /server ./server
# The agent's knowledge base. Mount a volume over /app/corpus to update it
# without rebuilding (docker-compose.yml does this).
COPY --from=builder /app/corpus ./corpus

ENV CORPUS_DIR=/app/corpus
EXPOSE 8080

ENTRYPOINT ["./server"]
