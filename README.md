# 🤖 Go-RAG-System

<div align="center">

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/pgvector-PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![OpenAI](https://img.shields.io/badge/OpenAI-API-412991?style=for-the-badge&logo=openai&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-a78bfa?style=for-the-badge)

**Production-grade Retrieval-Augmented Generation (RAG) architecture in Go**

[Architecture](#architecture) · [Quick Start](#quick-start) · [API Reference](#api-reference) · [Patterns](#rag-patterns)

</div>

---

## 🎯 Overview

A **production-ready RAG system** built in Go, integrating **pgvector** (PostgreSQL vector extension) for embedding storage and retrieval, and the **OpenAI API** for generation. This implements 17 documented agentic patterns and is designed for enterprise deployment — not a toy demo.

### Why Go for RAG?
- 🚀 **Performance** — goroutines for concurrent embedding and retrieval
- 🔒 **Type safety** — strongly typed request/response pipelines
- 📦 **Single binary** — trivial deployment, no runtime dependencies
- 🔧 **Low latency** — sub-10ms retrieval from pgvector at scale

---

## 🏗️ Architecture

```
User Query
    │
    ▼
┌─────────────────┐
│  Query Encoder   │  → OpenAI text-embedding-3-small
└────────┬────────┘
         │ embedding vector
         ▼
┌─────────────────┐
│  pgvector Store │  → cosine similarity search (top-k)
└────────┬────────┘
         │ retrieved chunks
         ▼
┌─────────────────┐
│  Context Builder │  → reranking + prompt construction
└────────┬────────┘
         │ enriched prompt
         ▼
┌─────────────────┐
│  OpenAI GPT-4o  │  → generation
└────────┬────────┘
         │
         ▼
    Final Response
```

---

## 🚀 Quick Start

### Prerequisites
- Go 1.22+
- PostgreSQL 15+ with pgvector extension
- OpenAI API key

### Setup

```bash
git clone https://github.com/Raphasha27/Go-RAG-System.git
cd Go-RAG-System

# Copy environment config
cp .env.example .env
# Add your OPENAI_API_KEY and DATABASE_URL

# Start PostgreSQL with pgvector
docker-compose up -d postgres

# Install dependencies
go mod download

# Run migrations
go run cmd/migrate/main.go

# Start the RAG API server
go run cmd/server/main.go
```

### Ingest Documents

```bash
curl -X POST http://localhost:8080/ingest \
  -H "Content-Type: application/json" \
  -d '{"path": "./docs", "chunk_size": 512, "overlap": 64}'
```

### Query

```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"question": "What are the system requirements?", "top_k": 5}'
```

---

## 🔧 RAG Patterns Implemented

| Pattern | Description |
|---------|-------------|
| **Naive RAG** | Basic embed → retrieve → generate |
| **Advanced RAG** | Query rewriting + reranking |
| **Hybrid Search** | Vector + BM25 keyword fusion |
| **HyDE** | Hypothetical Document Embeddings |
| **Agentic RAG** | Multi-step reasoning with tool use |
| **Self-RAG** | Retrieval with self-reflection |
| **RAPTOR** | Recursive abstractive processing |

---

## 📡 API Reference

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/ingest` | POST | Ingest documents and generate embeddings |
| `/query` | POST | RAG query with streaming response |
| `/health` | GET | Service health check |
| `/metrics` | GET | Prometheus metrics |

---

## 🗺️ Roadmap

- [ ] Multi-modal RAG (images + text)
- [ ] LangGraph-style agent orchestration
- [ ] Weaviate/Qdrant adapter support
- [ ] Streaming SSE responses
- [ ] Evaluation framework (RAGAS)

---

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Good first issues labelled!

---

## 📄 License

MIT License — see [LICENSE](LICENSE).

---

<div align="center">
Built by <a href="https://github.com/Raphasha27">Koketso Raphasha</a> · <a href="https://portfolio-iota-eight-90.vercel.app/">Portfolio</a>
</div>
