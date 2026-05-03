package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pgvector/pgvector-go"
	"github.com/sashabaranov/go-openai"
)

// Document represents a chunk of text retrieved from the Vector DB
type Document struct {
	ID        int             `json:"id"`
	Content   string          `json:"content"`
	Metadata  json.RawMessage `json:"metadata"`
	Embedding pgvector.Vector `json:"embedding"`
}

// -------------------------------------------------------------------------
// Helper: Get Embedding from OpenAI
// -------------------------------------------------------------------------
func getEmbedding(client *openai.Client, text string) ([]float32, error) {
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.SmallEmbedding3, // text-embedding-3-small (1536 dims)
	}

	resp, err := client.CreateEmbeddings(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}

// -------------------------------------------------------------------------
// Helper: Search Similar Documents in pgvector
// -------------------------------------------------------------------------
func searchSimilarDocs(db *sql.DB, queryEmbedding []float32, topK int) ([]Document, error) {
	// pgvector uses the `<->` operator for Euclidean distance, `<#>` for inner product, and `<=>` for cosine distance.
	// We use cosine distance here.
	query := `
		SELECT id, content, metadata
		FROM documents
		ORDER BY embedding <=> $1
		LIMIT $2;
	`
	
	// Convert standard float slice to pgvector format
	vec := pgvector.NewVector(queryEmbedding)

	rows, err := db.Query(query, vec, topK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []Document
	for rows.Next() {
		var doc Document
		// Metadata is JSONB
		if err := rows.Scan(&doc.ID, &doc.Content, &doc.Metadata); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

// -------------------------------------------------------------------------
// Helper: Build Augmented Prompt
// -------------------------------------------------------------------------
func buildPrompt(userQuery string, docs []Document) string {
	var contextBuilder strings.Builder

	contextBuilder.WriteString("Use the following pieces of retrieved context to answer the question. If you don't know the answer, just say that you don't know.\n\n")
	
	for i, doc := range docs {
		contextBuilder.WriteString(fmt.Sprintf("--- Context %d ---\n%s\n", i+1, doc.Content))
	}

	contextBuilder.WriteString(fmt.Sprintf("\nQuestion: %s\nAnswer:", userQuery))
	return contextBuilder.String()
}

// -------------------------------------------------------------------------
// Helper: Call LLM (with basic Tool Calling setup)
// -------------------------------------------------------------------------
func callLLM(client *openai.Client, prompt string) (string, error) {
	// Define a mock tool (e.g., Internet Search)
	tools := []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "web_search",
				Description: "Search the internet for real-time information",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"query": {
							Type: jsonschema.String,
							Description: "The search query to look up online",
						},
					},
					Required: []string{"query"},
				},
			},
		},
	}

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful and intelligent AI assistant. Use provided context and tools to give accurate answers.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Tools: tools, // Attach tools for LLM generation phase
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}

	message := resp.Choices[0].Message

	// If the LLM decided to call a tool
	if len(message.ToolCalls) > 0 {
		toolName := message.ToolCalls[0].Function.Name
		fmt.Printf("\n[🔧 Tool Call Triggered]: %s\n", toolName)
		
		// 🛡️ AGENT SECURITY LAYER (Authorization Check)
		// Simulating an agent trying to execute a tool
		agentRole := "agent" 
		
		if EnforceSecurity(agentRole, toolName) {
			// In a real app, you would execute the tool here and pass the result back to the LLM
			return "[Tool execution simulated and approved by security layer]", nil
		} else {
			// Fail-safe mechanism: Block execution
			return "Error: Agent is not authorized to execute this tool.", nil
		}
	}

	return message.Content, nil
}

// -------------------------------------------------------------------------
// MAIN APPLICATION
// -------------------------------------------------------------------------
func main() {
	traceID := generateTraceID()
	LogEvent(traceID, "INFO", "🚀 Initializing RAG System (Go + PostgreSQL + pgvector + OpenAI)")

	// 1. Connect to PostgreSQL
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://user:password@localhost:5432/ragdb?sslmode=disable"
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		LogEvent(traceID, "FATAL", fmt.Sprintf("Failed to open DB connection: %v", err))
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		LogEvent(traceID, "WARN", "Could not ping database. Continuing in mock mode.")
	}

	// 2. Initialize OpenAI client
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	userQuery := "How does tool calling work in RAG systems?"
	LogEvent(traceID, "INPUT", fmt.Sprintf("User Query: %s", userQuery))

	// 3. Get embedding for query
	LogEvent(traceID, "INFO", "Action: Embedding Query")
	queryEmbedding := make([]float32, 1536) 
	queryEmbedding[0] = 0.12

	// 4. Search similar documents
	LogEvent(traceID, "INFO", "Action: Retrieving Context from Vector DB")
	docs := []Document{
		{Content: "Tool calling allows LLMs to execute external functions like searching the web."},
	}

	// 5. Build prompt
	prompt := buildPrompt(userQuery, docs)

	// 6. Call LLM
	LogEvent(traceID, "INFO", "Action: Generating Final Answer")
	
	// Simulate Evaluation Metrics collection during run
	metrics := EvaluationMetrics{
		GoalAchievement: 0.95,
		TaskSuccessRate: 0.90,
		Efficiency:      0.85,
		SafetyScore:     1.00, // Agent didn't execute malicious tools
	}

	if apiKey != "" {
		answer, err := callLLM(client, prompt)
		if err != nil {
			LogEvent(traceID, "ERROR", fmt.Sprintf("LLM Error: %v", err))
		}
		LogEvent(traceID, "OUTPUT", "\n" + answer)
	} else {
		LogEvent(traceID, "OUTPUT", "(Mocked) Tool calling works by the LLM recognizing when it lacks information, pausing generation, requesting a tool, and then consuming the tool's result.")
	}

	// 7. Evaluate the Agent Run
	EvaluateAgent(traceID, metrics)

	// 8. Trigger Agentic Self-Correction & Reflection Loop Example
	fmt.Println("\n==================================================")
	fmt.Println("       AGENT REFLECTION & SELF-CORRECTION")
	fmt.Println("==================================================")
	SelfCorrectingExecution(traceID, "Find highly specific architectural context")

	// 9. Trigger Agent Planning & ReAct Framework
	fmt.Println("\n==================================================")
	fmt.Println("       AGENT TASK PLANNING & DECOMPOSITION")
	fmt.Println("==================================================")
	AgentPlanner(traceID, "Deploy highly secure, multi-agent RAG cluster")
	ReActExecution(traceID, "How does ReAct architecture work?")

	// 10. Trigger Human-in-the-loop (HITL) Execution
	HumanApproval(traceID, "Execute SQL query to clear 'temp_embeddings' table")

	LogEvent(traceID, "INFO", "Agent run entirely completed.")
}
