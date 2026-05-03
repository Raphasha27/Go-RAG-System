package main

import (
	"fmt"
	"time"
)

// -------------------------------------------------------------------------
// 🤖🤖 MULTI-AGENT SYSTEM & EXECUTOR AGENT
// -------------------------------------------------------------------------

// AgentWorker represents a discrete autonomous entity in a Multi-Agent System.
type AgentWorker struct {
	Name string
	Role string
}

func (a *AgentWorker) ExecuteTask(traceID string, task string) {
	LogEvent(traceID, "MULTI_AGENT_EXECUTE", fmt.Sprintf("[%s - %s] Executing task: '%s'", a.Name, a.Role, task))
	time.Sleep(300 * time.Millisecond) // Simulate real-world execution
	LogEvent(traceID, "MULTI_AGENT_DONE", fmt.Sprintf("[%s] Task completed successfully.", a.Name))
}

// DemonstrateMultiAgentSystem shows interaction between a Planner and an Executor.
func DemonstrateMultiAgentSystem(traceID string) {
	fmt.Println("\n==================================================")
	fmt.Println("             MULTI-AGENT SYSTEM & EXECUTOR")
	fmt.Println("==================================================")

	planner := &AgentWorker{Name: "Agent-Alpha", Role: "Planning Agent"}
	executor := &AgentWorker{Name: "Agent-Beta", Role: "Executor Agent"}

	LogEvent(traceID, "MULTI_AGENT_PLAN", fmt.Sprintf("%s is decomposing goal into discrete tasks...", planner.Name))
	
	tasks := []string{
		"Search Vector DB for contextual parameters",
		"Call Calculator Tool to aggregate metrics",
		"Format Final Markdown Report",
	}
	
	for _, task := range tasks {
		LogEvent(traceID, "MULTI_AGENT_DELEGATE", fmt.Sprintf("%s delegating task -> %s", planner.Name, executor.Name))
		executor.ExecuteTask(traceID, task)
	}

	LogEvent(traceID, "MULTI_AGENT_SUCCESS", "All agents collaborated successfully to achieve the global goal.")
}
