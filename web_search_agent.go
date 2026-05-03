package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// -------------------------------------------------------------------------
// 🌐 WEB SEARCH AGENT
// -------------------------------------------------------------------------

// WebSearchAgent simulates an AI agent that autonomously searches the web,
// extracts useful information, and provides accurate answers.
func WebSearchAgent(traceID string, query string) {
	fmt.Println("\n==================================================")
	fmt.Println("             WEB SEARCH AGENT")
	fmt.Println("==================================================")
	LogEvent(traceID, "WEB_SEARCH_START", fmt.Sprintf("Query: '%s'", query))

	// 1. User Query (Received)
	// 2. Web Search Execution
	apiKey := os.Getenv("SEARCH_API_KEY")
	engineID := os.Getenv("SEARCH_ENGINE_ID")

	if apiKey == "" || engineID == "" {
		LogEvent(traceID, "WARN", "Search API Keys not set. Running simulated web search extraction.")
		SimulatedWebSearch(traceID, query)
		return
	}

	// Real execution using Google Custom Search API
	encodedQuery := url.QueryEscape(query)
	searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&key=%s&cx=%s", encodedQuery, apiKey, engineID)

	resp, err := http.Get(searchURL)
	if err != nil {
		LogEvent(traceID, "ERROR", fmt.Sprintf("HTTP Get failed: %v", err))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		LogEvent(traceID, "ERROR", "Failed to read response body.")
		return
	}

	// 3. Extract & Analyze
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	items, ok := result["items"].([]interface{})
	if !ok {
		LogEvent(traceID, "ERROR", "No items found in search response.")
		return
	}

	// 4. Generate Response
	LogEvent(traceID, "EXTRACT", fmt.Sprintf("Extracted %d relevant links.", len(items)))
	for i := 0; i < len(items) && i < 3; i++ {
		item := items[i].(map[string]interface{})
		title := item["title"].(string)
		link := item["link"].(string)
		snippet := item["snippet"].(string)
		fmt.Printf("\n[Result %d]\nTitle: %s\nLink: %s\nSnippet: %s\n", i+1, title, link, snippet)
	}
	
	// 5. Final Output
	LogEvent(traceID, "WEB_SEARCH_COMPLETE", "Successfully provided accurate answers to the user.")
}

// SimulatedWebSearch provides a fallback mock execution if API keys are missing.
func SimulatedWebSearch(traceID string, query string) {
	LogEvent(traceID, "SIMULATE_SEARCH", "Executing simulated HTTP request for extraction...")
	fmt.Printf("\n[Result 1]\nTitle: Latest Advancements in AI\nLink: https://example.com/ai-news\nSnippet: Researchers have developed a new ReAct agent framework combining web search and semantic reasoning.\n")
	fmt.Printf("\n[Result 2]\nTitle: Go-RAG System Released\nLink: https://example.com/go-rag\nSnippet: A complete Golang implementation of 10 different Agentic AI patterns.\n")
	LogEvent(traceID, "EXTRACT", "Simulated data extraction and analysis complete.")
	LogEvent(traceID, "WEB_SEARCH_COMPLETE", "Successfully generated helpful response from simulated search data.")
}
