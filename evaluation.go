package main

import "fmt"

// -------------------------------------------------------------------------
// 📊 AGENT EVALUATION LAYER
// -------------------------------------------------------------------------

// EvaluationMetrics holds the performance data for a single agent execution run
type EvaluationMetrics struct {
	GoalAchievement  float64 // 0.0 to 1.0
	TaskSuccessRate  float64 // 0.0 to 1.0
	Efficiency       float64 // 0.0 to 1.0
	SafetyScore      float64 // 0.0 to 1.0
}

// EvaluateAgent calculates the overall performance score of the agent
func EvaluateAgent(traceID string, metrics EvaluationMetrics) {
	// Weighted scoring model
	overallScore := (metrics.GoalAchievement * 0.35) +
		(metrics.TaskSuccessRate * 0.25) +
		(metrics.Efficiency * 0.20) +
		(metrics.SafetyScore * 0.20)

	fmt.Println("\n==================================================")
	fmt.Println("             AGENT EVALUATION REPORT")
	fmt.Println("==================================================")
	fmt.Printf("Trace ID         : %s\n", traceID)
	fmt.Printf("Goal Achievement : %.2f%%\n", metrics.GoalAchievement*100)
	fmt.Printf("Task Success Rate: %.2f%%\n", metrics.TaskSuccessRate*100)
	fmt.Printf("Efficiency       : %.2f%%\n", metrics.Efficiency*100)
	fmt.Printf("Safety Score     : %.2f%%\n", metrics.SafetyScore*100)
	fmt.Printf("OVERALL SCORE    : %.2f%%\n", overallScore*100)
	
	if overallScore >= 0.8 {
		fmt.Println("Performance      : GOOD ✅")
	} else {
		fmt.Println("Performance      : NEEDS IMPROVEMENT ⚠️")
	}
	fmt.Println("==================================================\n")
}
