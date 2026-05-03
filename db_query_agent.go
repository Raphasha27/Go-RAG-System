package main

import (
	"database/sql"
	"fmt"
)

// -------------------------------------------------------------------------
// 🗄️ DATABASE QUERY AGENT (NL2SQL)
// -------------------------------------------------------------------------

// DatabaseQueryAgent simulates an agent that understands natural language,
// converts it to SQL, executes it securely, and returns the response.
func DatabaseQueryAgent(traceID string, db *sql.DB, naturalLanguageQuery string) {
	fmt.Println("\n==================================================")
	fmt.Println("             DATABASE QUERY AGENT")
	fmt.Println("==================================================")

	LogEvent(traceID, "USER_QUERY", naturalLanguageQuery)

	// 1. UNDERSTAND & TRANSLATE (NL2SQL)
	// The LLM parses intent and generates the appropriate SQL query.
	LogEvent(traceID, "TRANSLATE", "LLM understanding intent and converting Natural Language to SQL...")
	
	// Simulated Generated SQL for safety
	generatedSQL := "SELECT id, content FROM documents LIMIT 2;"
	LogEvent(traceID, "GENERATED_SQL", generatedSQL)

	// 2. EXECUTE
	LogEvent(traceID, "EXECUTE", "Executing SQL query securely on PostgreSQL...")
	
	if db == nil {
		LogEvent(traceID, "ERROR", "Database connection is nil. Ensure PostgreSQL is running.")
		return
	}

	// In a real application, prepared statements should be used to prevent SQL Injection
	rows, err := db.Query(generatedSQL)
	if err != nil {
		LogEvent(traceID, "ERROR", fmt.Sprintf("SQL Execution failed (Mock Mode Active): %v", err))
		return
	}
	defer rows.Close()

	// 3. RETURN RESPONSE
	LogEvent(traceID, "RESPONSE", "Results fetched successfully:")
	for rows.Next() {
		var id int
		var content string
		if err := rows.Scan(&id, &content); err == nil {
			fmt.Printf("   -> [Row ID: %d] %s\n", id, content[:min(len(content), 50)]+"...")
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
