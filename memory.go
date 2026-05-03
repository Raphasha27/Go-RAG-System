package main

import "fmt"

// -------------------------------------------------------------------------
// 🧠 AGENT MEMORY SYSTEM (SHORT-TERM & LONG-TERM)
// -------------------------------------------------------------------------

// AgentMemory handles both transient context (short-term) and persistent facts (long-term).
type AgentMemory struct {
	ShortTerm []string
	LongTerm  map[string]string
}

func NewAgentMemory() *AgentMemory {
	return &AgentMemory{
		ShortTerm: make([]string, 0),
		LongTerm:  make(map[string]string),
	}
}

func (m *AgentMemory) AddShortTerm(data string) {
	m.ShortTerm = append(m.ShortTerm, data)
}

func (m *AgentMemory) ClearShortTerm() {
	m.ShortTerm = nil
}

func (m *AgentMemory) AddLongTerm(key, value string) {
	m.LongTerm[key] = value
}

func (m *AgentMemory) GetLongTerm(key string) string {
	if val, ok := m.LongTerm[key]; ok {
		return val
	}
	return "Not Found"
}

func DemonstrateMemory(traceID string) {
	fmt.Println("\n==================================================")
	fmt.Println("             AGENT MEMORY MANAGEMENT")
	fmt.Println("==================================================")

	mem := NewAgentMemory()
	LogEvent(traceID, "MEMORY", "Initializing Agent Memory System...")

	// Simulate conversation turns
	mem.AddShortTerm("User asked about Bangalore weather.")
	mem.AddShortTerm("User asked for top 5 restaurants.")
	LogEvent(traceID, "MEMORY_SHORT", fmt.Sprintf("Stored %d items in short-term context.", len(mem.ShortTerm)))

	// Simulate extracting long-term facts
	mem.AddLongTerm("user_preference_city", "Bangalore")
	mem.AddLongTerm("user_name", "CodePathIndia Fan")
	LogEvent(traceID, "MEMORY_LONG", fmt.Sprintf("Saved persistent facts. City: %s", mem.GetLongTerm("user_preference_city")))

	// Session Ends
	mem.ClearShortTerm()
	LogEvent(traceID, "MEMORY_CLEAR", "Short-Term Memory Cleared for next session. Long-term memory persists.")
}
