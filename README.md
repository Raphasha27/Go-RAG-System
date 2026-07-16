# 🧠 Go-RAG-System
### reference implementation Retrieval Augmented Generation Architecture in Golang

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-pgvector-336791?logo=postgresql&logoColor=white)](https://github.com/pgvector/pgvector)
[![OpenAI](https://img.shields.io/badge/OpenAI-Embeddings-412991?logo=openai&logoColor=white)](https://platform.openai.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Active%20Development-brightgreen)]()

---

## 🧬 Project Overview

The **Go-RAG-System** is a complete, production-ready implementation of a **Retrieval-Augmented Generation (RAG)** pipeline, written entirely in **Go**. It combines the raw retrieval power of a **PostgreSQL + pgvector** vector database with the language generation capabilities of **OpenAI GPT-4o** and an advanced **17-pattern Agentic Architecture**.

Unlike toy implementations, this system is designed to handle real-world complexity:
- It doesn't just generate text — it **reasons, plans, executes tools, and self-corrects**.
- It doesn't just store data — it indexes, retrieves, and augments with **semantic similarity search** using 1536-dimensional embeddings.
- It doesn't just run one AI call — it orchestrates **multi-agent pipelines** with full observability, security authorization, and human approval gates.

---

## 🏗️ System Architecture

The RAG workflow runs through **6 distinct phases** before a response is returned:

```
 User Query
     │
     ▼
┌─────────────────────────────────────┐
│  Phase 1: Embedding                 │
│  text-embedding-3-small (1536d)     │
└─────────────────┬───────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Phase 2: Vector DB Retrieval       │
│  PostgreSQL + pgvector              │
│  IVFFlat Cosine Similarity Search   │
│  → Returns Top K Document Chunks    │
└─────────────────┬───────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Phase 3: Prompt Augmentation       │
│  Retrieved chunks + Original query  │
│  → Enriched context prompt          │
└─────────────────┬───────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Phase 4: LLM Generation            │
│  GPT-4o processes enriched prompt   │
│  → Generates precise answer         │
└─────────────────┬───────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Phase 5: Tool Calling (if needed)  │
│  LLM triggers Web Search / DB Query │
│  → External data fetched & merged   │
└─────────────────┬───────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│  Phase 6: Agentic Post-Processing   │
│  Observability → Evaluation         │
│  → Reflection → Final Response      │
└─────────────────────────────────────┘
```

---

## 🗄️ Database Schema (PostgreSQL + pgvector)

Before running the application, ensure your PostgreSQL instance has `pgvector` installed and execute the following initialization SQL:

```sql
-- Step 1: Enable the pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Step 2: Create the documents table to store text chunks and their embeddings
CREATE TABLE documents (
    id          SERIAL PRIMARY KEY,
    content     TEXT NOT NULL,
    embedding   VECTOR(1536),       -- 1536 dimensions matches text-embedding-3-small
    metadata    JSONB,              -- Supports hybrid filtering (e.g., source, date)
    created_at  TIMESTAMP DEFAULT NOW()
);

-- Step 3: Create an IVFFlat index for fast approximate nearest-neighbor search
-- 'lists' controls the number of cluster centroids (tune based on dataset size)
CREATE INDEX ON documents
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);

-- Optional: HNSW index for higher recall at slightly more memory cost
-- CREATE INDEX ON documents USING hnsw (embedding vector_cosine_ops);
```

### Index Selection Guide

| Index Type | Recall | Speed | Memory | Best For |
|:---|:---|:---|:---|:---|
| **IVFFlat** | ~95% | Very Fast | Low | Datasets < 1M rows |
| **HNSW** | ~99% | Fast | Higher | High-precision requirements |
| No Index (exact) | 100% | Slow | Minimal | Dev / < 100k rows |

### Why PostgreSQL + pgvector?

*   **ACID Compliance:** Your document store and relational data live in one reliable, transactional database — no eventual consistency issues.
*   **Hybrid Search:** Combine vector similarity with SQL `WHERE` clauses. Example: "Find semantically similar documents, but only from the last 30 days, tagged as `compliance`".
*   **Familiar Tooling:** Use standard PostgreSQL tooling (pgAdmin, psql, pg_dump) for backups, monitoring, and schema migrations.
*   **No Additional Infrastructure:** Eliminates the need for a separate dedicated vector database (Pinecone, Weaviate, etc.) reducing cost and ops overhead.

---

## 📁 Repository Structure

```text
Go-RAG-System/
│
├── main.go                 # Application entry point — wires all components
├── security.go             # Agent Security & Authorization (RBAC checks)
├── observability.go        # TraceID generation, structured logging, metrics
├── evaluation.go           # Response quality & safety scoring
├── agent_logic.go          # Reflection & Self-Correction loop
├── react_planner.go        # ReAct (Reason + Act) Planning Framework
├── human_in_the_loop.go    # HITL Approval Gates
├── memory.go               # Short-term & Long-term Agent Memory
├── single_agent.go         # Perceive → Reason → Act → Learn loop
├── db_query_agent.go       # NL2SQL — Natural Language to SQL
├── web_search_agent.go     # Autonomous Web Research (SerpAPI)
├── ml_models.go            # Dynamic AI Model Selection Engine
├── multi_agent.go          # Multi-Agent Coordinator (Planner + Executor)
├── tool_caller.go          # Tool & External API Calling Interface
├── rest_api.go             # REST API CRUD + Bulk Operations mapping
├── goal_based_agent.go     # BFS Goal-Based Pathfinding Agent
├── concurrency.go          # Goroutine-based Concurrent Task Execution
├── custom_errors.go        # AgentHallucinationError & structured exceptions
│
├── go.mod                  # Go module definition
├── go.sum                  # Dependency lock file
├── .env.example            # Environment variable template
├── .github/
│   └── workflows/
│       └── daily-contribution-sync.yml  # Automated daily maintenance
└── README.md
```

---

## 🤖 The 17-Pattern Agentic Architecture

This repository implements every major pattern from modern agentic AI research and engineering:

### Core Intelligence Patterns

**1. Agent Observability (`observability.go`)**
Every agent execution generates a pseudo-UUID `TraceID`. Structured log events are emitted at every stage — tool call attempts, LLM responses, evaluation scores — providing a complete, reproducible audit trail.

**2. Agent Evaluation (`evaluation.go`)**
After each response, scores are computed across four dimensions:
- Goal Achievement: Did the agent actually answer the question?
- Task Success Rate: How many sub-steps completed without error?
- Efficiency: Response latency relative to complexity.
- Safety Score: Did the response stay within ethical bounds?

**3. Reflection & Self-Correction (`agent_logic.go`)**
When evaluation fails, instead of raising an error, the agent enters a self-correction loop. It analyzes *why* it failed (e.g., "Irrelevant search results", "Hallucination detected"), updates its strategy ("Add metadata filter", "Use a different tool"), and retries.

```go
// The agent catches its own failure and corrects:
// Attempt 1: Failed — "Irrelevant results"
// Correction: "Narrow query to last 7 days using metadata filter"
// Attempt 2: Success
```

### Planning & Reasoning Patterns

**4. ReAct Framework (`react_planner.go`)**
The agent actively reasons in a loop: **Thought → Action → Observation → Thought**. Each iteration brings it closer to the goal by incorporating real-world feedback from tool executions.

**5. Agent Planning / Task Decomposition (`react_planner.go`)**
High-level user goals (e.g., "Prepare a Q1 compliance report") are autonomously broken into executable sub-tasks before any action is taken, preventing cascading failures from ambiguous prompts.

**6. Goal-Based Agent (`goal_based_agent.go`)**
Uses BFS (Breadth-First Search) pathfinding to evaluate multiple action paths and select the most optimal sequence to reach a desired end state — critical for multi-step workflows.

### Safety & Control Patterns

**7. Agent Security & Authorization (`security.go`)**
Enforces **Principle of Least Privilege** at the agent layer. Every tool call is checked against an RBAC table before execution.

```go
// Role: "agent" → Action: "execute" → APPROVED ✅
// Role: "agent" → Action: "delete"  → DENIED ❌ (fallback triggered)
```

**8. Human-in-the-Loop (`human_in_the_loop.go`)**
For destructive or irreversible actions, the agent halts and routes to a human approval gate. In production, this triggers a webhook notification to an admin interface and awaits a response.

### Data & Tool Patterns

**9. Agent Memory (`memory.go`)**
Dual-layer memory:
- **Short-Term:** An in-memory key-value store cleared at session end.
- **Long-Term:** PostgreSQL-backed persistent memory for cross-session recall.

**10. Tool Calling (`tool_caller.go`)**
Exposes a structured JSON interface enabling the LLM to invoke real-world tools — Calculators, Weather APIs, Custom REST endpoints — extending its capabilities far beyond its static training data.

**11. Database Query Agent / NL2SQL (`db_query_agent.go`)**
Translates natural language questions into optimized SQL queries, executes them against the PostgreSQL store, and synthesizes a human-readable response from the result set.

**12. Web Search Agent (`web_search_agent.go`)**
Connects to SerpAPI to retrieve real-time web results. Parses the JSON payload, extracts the top results, and synthesizes a grounded, citation-backed response.

### Scalability Patterns

**13. Single Agent Workflow (`single_agent.go`)**
A formalized `Perceive → Reason → Act → Learn` loop — the fundamental building block for all agentic systems.

**14. Multi-Agent System (`multi_agent.go`)**
A **Planner Agent** decomposes a complex goal and dispatches sub-tasks to specialized **Executor Agents** running concurrently via goroutines.

**15. Concurrent Task Execution (`concurrency.go`)**
Leverages native Go goroutines and `sync.WaitGroup` to execute multiple agent tasks in parallel — no blocking, no deadlocks.

**16. Dynamic AI Model Selection (`ml_models.go`)**
Selects the appropriate algorithmic approach (NLP, Neural Network, Linear Regression) at runtime based on the semantic classification of the task.

**17. Custom Exceptions (`custom_errors.go`)**
Implements `AgentHallucinationError` and structured exception types that carry rich context (trace ID, tool name, expected vs. actual values) enabling precise error recovery.

---

## 🚀 Getting Started

### Prerequisites

| Dependency | Version | Notes |
|:---|:---|:---|
| Go | 1.22+ | [Install from go.dev](https://go.dev/dl/) |
| PostgreSQL | 15+ | With `pgvector` extension installed |
| OpenAI API Key | — | Billing account required |

### Installation

**1. Clone the repository:**
```bash
git clone https://github.com/Raphasha27/Go-RAG-System.git
cd Go-RAG-System
```

**2. Install Go dependencies:**
```bash
go mod tidy
```

**3. Configure environment variables:**
```bash
cp .env.example .env
```

Edit `.env`:
```env
DATABASE_URL="postgres://user:password@localhost:5432/ragdb?sslmode=disable"
OPENAI_API_KEY="sk-your-openai-api-key"
```

**4. Initialize the database:**
```bash
psql -U your_user -d ragdb -f schema.sql
```

**5. Run the application:**
```bash
go run main.go
```

Expected output:
```
[TraceID: a1b2c3d4] RAG System Starting...
[TraceID: a1b2c3d4] Connecting to PostgreSQL...
[TraceID: a1b2c3d4] Embedding query: "What is our refund policy?" (1536 dims)
[TraceID: a1b2c3d4] Retrieved 3 context chunks from pgvector
[TraceID: a1b2c3d4] 🔧 Tool Call Triggered: web_search
[TraceID: a1b2c3d4] Agent Evaluation: Quality=95 Safety=100 ✅ PASSED
```

---

## 🛡️ Security Best Practices

The system implements the following security controls at the agent layer:

| Control | Implementation | File |
|:---|:---|:---|
| Privilege Limitation | RBAC — agents cannot execute `delete` or `write` | `security.go` |
| Audit Logging | Every tool call logged with TraceID | `observability.go` |
| Hallucination Detection | `AgentHallucinationError` thrown + caught | `custom_errors.go` |
| Human Override Gate | High-risk actions paused for approval | `human_in_the_loop.go` |
| Input Validation | All queries validated before vectorization | `main.go` |

---

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run with race condition detector (recommended before production)
go test ./... -race
```

---

## 🔧 Troubleshooting

### Issue: `pgvector` extension not found
**Fix:** Ensure the extension is installed in your PostgreSQL instance:
```sql
CREATE EXTENSION IF NOT EXISTS vector;
-- If this fails, install via: sudo apt install postgresql-15-pgvector
```

### Issue: Goroutine leak detected under `-race`
**Cause:** A concurrent task may be holding a resource without releasing it.
**Fix:** Review `concurrency.go`. Ensure all goroutines are wrapped with a `defer wg.Done()` call.

### Issue: OpenAI 429 Rate Limit Error
**Fix:** Add exponential backoff retry logic around the `client.chat.completions.create()` call. The current implementation will return `nil` on failure — a future PR will add automatic retry.

### Issue: IVFFlat index returning poor results
**Cause:** The `lists` parameter may be too small for your dataset size.
**Fix:** PostgreSQL recommends `lists = sqrt(rows)`. For 1M rows, use `lists = 1000`.

---

## 🗺️ Roadmap

- [ ] Add a REST HTTP server (using `net/http` or `gin`) to expose as a microservice
- [ ] Implement HNSW index as a configuration option alongside IVFFlat
- [ ] Add streaming response support (Server-Sent Events)
- [ ] Integrate Weaviate as an alternative vector backend
- [ ] Build a CLI tool for document ingestion and index management
- [ ] Add Prometheus metrics endpoint for production monitoring

---

## 💼 Real-World Use Cases

*   **Chat with your documents / PDFs:** Index company policy documents, technical manuals, or legal contracts and query them in plain English.
*   **Enterprise Knowledge Base:** Build a Confluence / SharePoint replacement that actually understands what you're asking.
*   **Customer Support Bots:** Ground AI responses in your actual product documentation to eliminate hallucinations.
*   **Research Assistants:** Index academic papers and ask cross-paper synthesis questions.
*   **Automated Compliance Reporting:** Query regulatory documents and generate gap analysis reports.

---

## 🧠 Why RAG Over Fine-Tuning?

| Aspect | Fine-Tuning | RAG |
|:---|:---|:---|
| **Cost** | Thousands of dollars per run | One-time embedding cost |
| **Freshness** | Stale after training cutoff | Real-time via DB updates |
| **Explainability** | Black box | Traceable to source chunks |
| **Domain Adaptation** | Requires large labeled dataset | Works with any documents |
| **Deployment** | New model version per update | Index update only |

---

## 🤝 Contributing

1.  Fork the repository.
2.  Create a feature branch: `git checkout -b feat/your-feature`
3.  Commit your changes with a conventional commit message: `git commit -m 'feat: add HNSW index support'`
4.  Push to the branch: `git push origin feat/your-feature`
5.  Open a Pull Request describing the change and its motivation.

---

## 📜 License

MIT License. See [LICENSE](LICENSE) for details.

---

## 🔗 Ecosystem

Part of the **Kirov Dynamics Technology** ecosystem:

[![Portfolio](https://img.shields.io/badge/Portfolio-⭐29-00ffcc?style=flat-square)](https://github.com/Raphasha27/Portfolio)
[![AI-Agent](https://img.shields.io/badge/AI--Agent-⭐3-004a99?style=flat-square)](https://github.com/Raphasha27/AI-Agent)
[![Github-Harden](https://img.shields.io/badge/Github--Harden-Security-00ffcc?style=flat-square)](https://github.com/Raphasha27/Github-Harden)
[![Go-RAG-System](https://img.shields.io/badge/Go--RAG--System-00ADD8?style=flat-square)](https://github.com/Raphasha27/Go-RAG-System)
[![Nexus-Quant](https://img.shields.io/badge/Nexus--Quant-Quant-00ffcc?style=flat-square)](https://github.com/Raphasha27/Nexus-Quant)
[![CyberShield SOC](https://img.shields.io/badge/CyberShield--SOC-Security-004a99?style=flat-square)](https://github.com/Raphasha27/cybershield_soc)
[![Dev Factory](https://img.shields.io/badge/Dev--Factory-v7-005571?style=flat-square)](https://github.com/Raphasha27/autonomous-dev-factory-v7)
[![SaaS Backend](https://img.shields.io/badge/SaaS--Backend-Multi--tenant-004a99?style=flat-square)](https://github.com/Raphasha27/saas-multitenant-backend)
[![FastAPI Starter](https://img.shields.io/badge/FastAPI--Starter-Enterprise-005571?style=flat-square)](https://github.com/Raphasha27/enterprise-fastapi-starter)
[![Repo Audit](https://img.shields.io/badge/Repo--Audit--Bot-CLI-00ffcc?style=flat-square)](https://github.com/Raphasha27/repo-audit-bot)

*Building the infrastructure of autonomous systems.*

---

*Part of the **Future AGI** ecosystem — built by **Koketso Raphasha (Raphasha27)**.
*Kirov Dynamics Technology | Building the Infrastructure of Autonomous Systems.*

