package model

import (
	"crypto/md5"
	"fmt"
)

// TaskType represents the type of evaluation task
type TaskType string

const (
	TaskMedicalQA TaskType = "medical_qa"
	TaskReasoning TaskType = "reasoning"
	TaskWorkflow  TaskType = "workflow"
	TaskRAG       TaskType = "rag"
)

// EvaluationCase represents a single evaluation case
type EvaluationCase struct {
	ID       string
	TaskType TaskType
	Question string
	Options  []string
	Answer   string
	Context  string
	Metadata map[string]string
}

// GenerateID creates a unique ID from case content
func GenerateID(content string) string {
	hash := md5.Sum([]byte(content))
	return fmt.Sprintf("%x", hash)[:12]
}

// Dataset represents a collection of evaluation cases
type Dataset struct {
	Name        string
	Source      string
	TaskType    TaskType
	TotalCases  int
	Cases       []EvaluationCase
	Description string
}

// Filter returns a new dataset with at most n cases
func (d *Dataset) Filter(n int) *Dataset {
	if n >= len(d.Cases) {
		return d
	}
	return &Dataset{
		Name:        d.Name,
		Source:      d.Source,
		TaskType:    d.TaskType,
		TotalCases:  n,
		Cases:       d.Cases[:n],
		Description: d.Description,
	}
}

// NewEvaluationCase creates a new evaluation case with generated ID
func NewEvaluationCase(question, answer string, options []string, taskType TaskType, metadata map[string]string) EvaluationCase {
	content := question + answer
	return EvaluationCase{
		ID:       GenerateID(content),
		Question: question,
		Answer:   answer,
		Options:  options,
		TaskType: taskType,
		Metadata: metadata,
	}
}
