package main

import "fmt"

// -------------------------------------------------------------------------
// 🛠️ TOOL CALLING & API INTEGRATION LAYER
// -------------------------------------------------------------------------

// Tool Interface defines a common structure for all tools the agent can use.
type Tool interface {
	GetName() string
	Execute(input string) string
}

// CalculatorTool performs arithmetic.
type CalculatorTool struct{}

func (t *CalculatorTool) GetName() string { return "Calculator" }
func (t *CalculatorTool) Execute(input string) string {
	return fmt.Sprintf("Calculation result for '%s' is mathematically validated.", input)
}

// WeatherAPITool fetches real-time external data.
type WeatherAPITool struct{}

func (t *WeatherAPITool) GetName() string { return "WeatherAPI" }
func (t *WeatherAPITool) Execute(input string) string {
	return fmt.Sprintf("Real-time weather API called for %s: 28°C, Partly Cloudy.", input)
}

// DemonstrateToolCalling shows how the LLM decides to use external tools
func DemonstrateToolCalling(traceID string) {
	fmt.Println("\n==================================================")
	fmt.Println("             AGENT TOOL & API CALLING")
	fmt.Println("==================================================")

	// Register tools
	tools := map[string]Tool{
		"calc":    &CalculatorTool{},
		"weather": &WeatherAPITool{},
	}

	LogEvent(traceID, "LLM_DECISION", "LLM identified need for real-time external data.")
	LogEvent(traceID, "TOOL_CALL", "Agent generated structured tool call: {name: 'WeatherAPI', arg: 'Bangalore'}")
	
	// Execute the tool
	weatherTool := tools["weather"]
	result := weatherTool.Execute("Bangalore")
	
	LogEvent(traceID, "TOOL_RESPONSE", fmt.Sprintf("Tool returned: %s", result))
	LogEvent(traceID, "LLM_SYNTHESIS", "Agent synthesizes the tool response to answer the user.")
}
