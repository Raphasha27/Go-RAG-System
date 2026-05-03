# 🧠 Go-RAG-System (Retrieval Augmented Generation)

Welcome to the **Go-RAG-System**. This repository contains a complete architecture for building a Retrieval-Augmented Generation (RAG) system utilizing **Golang**, **PostgreSQL (with pgvector)**, **OpenAI Embeddings**, and **Tool Calling**.

This architecture combines the reasoning power of Large Language Models (LLMs) with real-time data retrieval from your own knowledge base to generate accurate, highly contextual, and actionable responses without fine-tuning.

---

## 🏗️ System Architecture

The workflow consists of 6 core phases:

1.  **User Query:** The user asks a question or issues a command.
2.  **Query Embedding:** The system converts the text query into a high-dimensional vector using an embedding model (e.g., `text-embedding-3-small`).
3.  **Vector DB Retrieval:** A similarity search (Cosine Distance / HNSW) is executed against a PostgreSQL database running the `pgvector` extension to fetch the `Top K` most relevant document chunks.
4.  **Augmentation:** The retrieved textual context is merged with the original query to construct an enriched prompt.
5.  **LLM Generation:** The language model (e.g., GPT-4o) processes the augmented prompt to formulate a precise answer.
6.  **Tool Calling (If needed):** If the LLM determines it needs real-time, external data (e.g., Web Search, DB Lookup, API calls), it triggers a tool execution and feeds the result back into the generation loop.

---

## 🗄️ Database Schema (PostgreSQL + pgvector)

Before running the application, ensure your PostgreSQL instance has `pgvector` installed and execute the following SQL:

```sql
-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Create documents table to store chunks and embeddings
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    embedding VECTOR(1536), -- 1536 dimensions for text-embedding-3-small
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create index for fast similarity search
CREATE INDEX ON documents 
USING ivfflat (embedding vector_cosine_ops) 
WITH (lists = 100);
```

### Why PostgreSQL + pgvector?
*   **ACID Compliant & Reliable**: Keeps your relational data and vectors in the same robust database.
*   **Supports Metadata (JSONB)**: Perfect for hybrid search (Vector + SQL filtering).
*   **Powerful Indexing**: Supports IVFFlat and HNSW for rapid similarity lookups.

## 🛡️ Agent Security & Authorization

Securing AI agents is essential to ensure they act safely, reliably, and in alignment with user goals and ethical boundaries. This repository implements **Agent Security Best Practices**:

*   **Input Validation:** Sanitization of all user inputs before vectorization.
*   **Privilege Limitation (RBAC):** Applies the Principle of Least Privilege. The LLM agent is restricted to specific actions (`read`, `execute`) and blocked from destructive actions (`delete`, `write`) unless explicitly elevated.
*   **Audit & Logging:** Every tool call decision made by the agent passes through an audit trace.
*   **Fail-Safe Mechanisms:** If the agent hallucinates a tool call it shouldn't execute, the authorization layer dynamically blocks the execution and returns an error to the LLM.

```go
// Example: Agent tries to execute 'web_search'
role := "agent"
action := "execute" // Approved!

// Example: Agent maliciously tries to execute 'delete_record'
action := "delete" // Denied! Fallback triggered.
```

## 📊 Advanced Agentic Capabilities Implemented
This blueprint doesn't just stop at generating text. It incorporates deep **Agentic AI** principles for autonomous safety and learning:

1.  **Agent Observability (Trace, Log, Metrics):** Every action generates a pseudo-UUID `TraceID` and pushes structured logs across execution boundaries to ensure deep accountability.
2.  **Agent Evaluation:** At the end of execution, the agent evaluates its `Goal Achievement`, `Task Success Rate`, `Efficiency`, and `Safety Score` to produce an automated performance review.
3.  **Reflection & Self-Correction:** Instead of failing silently, the agent analyzes *why* an action failed (e.g., "Irrelevant Results"), learns from the mistake ("Hallucination risk detected"), and dynamically self-corrects its strategy ("Use metadata filters") in a self-healing loop.
4.  **Agent Planning (Task Decomposition):** High-level user goals are autonomously broken down into executable, step-by-step sub-tasks before execution begins.
5.  **ReAct Framework (Reason + Act):** The agent actively "thinks" about its situation, decides on a tool to use, observes the output of that tool, and then reasons again in a continuous loop until the goal is met.
6.  **Human-in-the-loop (HITL):** For critical or destructive tasks, the AI pauses execution and defers to a human operator for final Approval, Modification, or Rejection.

---

## 🚀 Getting Started

### Prerequisites
*   [Go 1.22+](https://go.dev/)
*   PostgreSQL Database with `pgvector`
*   OpenAI API Key

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/Raphasha27/Go-RAG-System.git
    cd Go-RAG-System
    ```

2.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Environment Variables:**
    Copy `.env.example` to `.env` and fill in your credentials:
    ```bash
    DATABASE_URL="postgres://user:password@localhost:5432/ragdb?sslmode=disable"
    OPENAI_API_KEY="sk-your-openai-key"
    ```

4.  **Run the Application:**
    ```bash
    go run main.go
    ```

---

## 🧠 Benefits of RAG

*   **Reduces Hallucinations:** Answers are grounded in actual retrieved documents.
*   **Uses up-to-date information:** Instantly reflects changes in your database.
*   **Domain-Specific:** Understands your private corporate data.
*   **Cost-Effective:** Eliminates the need for expensive, continuous model fine-tuning.
*   **Explainable:** You can trace exactly which document chunks were used to generate the answer.

## 💼 Use Cases
*   Chat with your specific documents / PDFs
*   Enterprise Knowledge Base search
*   Customer Support Bots
*   Research Assistants
*   Automated Data Analysis & Reporting
