package structuredparse

import (
	"strconv"
	"strings"
	"testing"
)

// BenchmarkParse_SmallInput benchmarks Parse with a small input (3-5 labels, ~200-500 bytes).
func BenchmarkParse_SmallInput(b *testing.B) {
	labels := []Label{
		{Name: "Reason"},
		{Name: "Function"},
		{Name: "Parameters", IsJSON: true},
	}

	parser, err := NewParser(labels, nil)
	if err != nil {
		b.Fatalf("failed to create parser: %v", err)
	}

	text := `Reason: I need to process some files.
Function: process_data
Parameters: {"input_files": ["a.txt", "b.txt"], "output_dir": "out/"}
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(text)
	}
}

// BenchmarkParse_MediumInput benchmarks Parse with a medium input (~10-15 labels, ~1-5 KB).
func BenchmarkParse_MediumInput(b *testing.B) {
	labels := []Label{
		{Name: "Context"},
		{Name: "Intention"},
		{Name: "Role"},
		{Name: "Action"},
		{Name: "Outcome"},
		{Name: "Notes"},
		{Name: "Config", IsJSON: true},
		{Name: "Metadata", IsJSON: true},
		{Name: "Summary"},
		{Name: "Details"},
		{Name: "Status"},
		{Name: "Timestamp"},
	}

	parser, err := NewParser(labels, nil)
	if err != nil {
		b.Fatalf("failed to create parser: %v", err)
	}

	// Build a medium-sized input
	var textBuilder strings.Builder
	textBuilder.WriteString("Context: This is a test context that provides background information.\n")
	textBuilder.WriteString("Intention: To CHECK the system functionality.\n")
	textBuilder.WriteString("Role: AGENT\n")
	textBuilder.WriteString("Action: process_data\n")
	textBuilder.WriteString("Outcome: Success\n")
	textBuilder.WriteString("Notes: This is a multiline note that continues on the next line.\n")
	textBuilder.WriteString("  It has multiple lines of content to simulate real-world usage.\n")
	textBuilder.WriteString("Config: {\"threshold\": 0.8, \"enabled\": true, \"timeout\": 30}\n")
	textBuilder.WriteString("Metadata: {\"version\": \"1.0\", \"author\": \"test\", \"tags\": [\"test\", \"benchmark\"]}\n")
	textBuilder.WriteString("Summary: A comprehensive summary of the operation.\n")
	textBuilder.WriteString("Details: Additional details about the process and its results.\n")
	textBuilder.WriteString("Status: completed\n")
	textBuilder.WriteString("Timestamp: 2024-01-01T00:00:00Z\n")

	text := textBuilder.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.Parse(text)
	}
}

// BenchmarkParseBlocks_MultipleBlocks benchmarks ParseBlocks with multiple blocks (5-10 blocks).
func BenchmarkParseBlocks_MultipleBlocks(b *testing.B) {
	labels := []Label{
		{Name: "Task", IsBlockStart: true},
		{Name: "Input", IsJSON: true},
		{Name: "Result"},
		{Name: "Status"},
	}

	parser, err := NewParser(labels, nil)
	if err != nil {
		b.Fatalf("failed to create parser: %v", err)
	}

	// Build input with multiple blocks
	var textBuilder strings.Builder
	for i := 1; i <= 8; i++ {
		iStr := strconv.Itoa(i)
		textBuilder.WriteString("Task: Task ")
		textBuilder.WriteString(iStr)
		textBuilder.WriteString("\n")
		textBuilder.WriteString("Input: {\"id\": ")
		textBuilder.WriteString(iStr)
		textBuilder.WriteString(", \"data\": \"block ")
		textBuilder.WriteString(iStr)
		textBuilder.WriteString(" data\"}\n")
		textBuilder.WriteString("Result: Result for task ")
		textBuilder.WriteString(iStr)
		textBuilder.WriteString("\n")
		textBuilder.WriteString("Status: completed\n")
		if i < 8 {
			textBuilder.WriteString("\n")
		}
	}

	text := textBuilder.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parser.ParseBlocks(text)
	}
}
