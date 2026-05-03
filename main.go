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
		fmt.Printf("\n[🔧 Tool Call Triggered]: %s\n", message.ToolCalls[0].Function.Name)
		// In a real app, you would execute the tool here and pass the result back to the LLM
		return "[Tool execution deferred in this snippet]", nil
	}

	return message.Content, nil
}

// -------------------------------------------------------------------------
// MAIN APPLICATION
// -------------------------------------------------------------------------
func main() {
	fmt.Println("🚀 Initializing RAG System (Go + PostgreSQL + pgvector + OpenAI)")

	// 1. Connect to PostgreSQL
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fallback for demonstration if env var is missing
		connStr = "postgres://user:password@localhost:5432/ragdb?sslmode=disable"
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Printf("⚠️ Warning: Could not ping database. Make sure PostgreSQL is running. Error: %v\n", err)
		log.Println("Continuing in mock mode for demonstration purposes...")
		// For the sake of the structural blueprint, we don't fatal crash if DB isn't up
	}

	// 2. Initialize OpenAI client
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("⚠️ Warning: OPENAI_API_KEY not set. Generation will fail.")
	}
	client := openai.NewClient(apiKey)

	// User Query
	userQuery := "How does tool calling work in RAG systems?"
	fmt.Printf("\n[👤 User Query]: %s\n", userQuery)

	// 3. Get embedding for query
	fmt.Println("[🧠 Action]: Embedding Query...")
	// NOTE: Bypassed for pure demo compilation without API key
	/*
	queryEmbedding, err := getEmbedding(client, userQuery)
	if err != nil {
		log.Fatalf("Embedding error: %v", err)
	}
	*/
	
	// Mock embedding for structural completeness
	queryEmbedding := make([]float32, 1536) 
	queryEmbedding[0] = 0.12
	queryEmbedding[1] = -0.85

	// 4. Search similar documents
	fmt.Println("[🗄️ Action]: Retrieving Context from Vector DB...")
	/*
	docs, err := searchSimilarDocs(db, queryEmbedding, 3)
	if err != nil {
		log.Fatalf("Search error: %v", err)
	}
	*/
	
	// Mock Retrieved Documents
	docs := []Document{
		{Content: "Tool calling allows LLMs to execute external functions like searching the web or querying a database. The result is returned to the LLM to augment its final answer."},
		{Content: "RAG stands for Retrieval-Augmented Generation. It combines semantic search with language models."},
	}

	// 5. Build prompt with context
	fmt.Println("[🧩 Action]: Augmenting Prompt...")
	prompt := buildPrompt(userQuery, docs)

	// 6. Call LLM
	fmt.Println("[🤖 Action]: Generating Final Answer...")
	if apiKey != "" {
		answer, err := callLLM(client, prompt)
		if err != nil {
			log.Fatalf("LLM Error: %v", err)
		}
		fmt.Printf("\n[✨ Answer]:\n%s\n", answer)
	} else {
		fmt.Println("\n[✨ Answer]:\n(Mocked) Tool calling works by the LLM recognizing when it lacks information, pausing generation, returning a specific JSON schema requesting a function execution, waiting for the external system to run the tool, and then consuming the tool's result to generate the final response.")
	}
}
