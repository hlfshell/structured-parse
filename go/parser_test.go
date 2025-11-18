package structuredparse

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

// Test scaffolding for parser, will load test cases from assets.

// TestBasicFunctionality verifies that the parser correctly parses a typical input with all fields present.
func TestBasicFunctionality(t *testing.T) {
	// Load input text from asset file
	input, err := os.ReadFile("../test-assets/basic_functionality_input.txt")
	if err != nil {
		t.Fatalf("failed to read input asset: %v", err)
	}

	// Load expected output from asset file
	expectedBytes, err := os.ReadFile("../test-assets/basic_functionality_output.json")
	if err != nil {
		t.Fatalf("failed to read output asset: %v", err)
	}
	var expected map[string]interface{}
	if err := json.Unmarshal(expectedBytes, &expected); err != nil {
		t.Fatalf("failed to unmarshal expected output: %v", err)
	}

	// Define parser labels as in the Python test
	labels := []Label{
		{Name: "Action Input", RequiredWith: []string{"Action"}, IsJSON: true},
		{Name: "Action", RequiredWith: []string{"Action Input"}},
		{Name: "Thought"},
		{Name: "Result", Required: true},
	}
	parser, err := NewParser(labels, nil)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	// Parse the input
	result, errors := parser.Parse(string(input))
	if len(errors) > 0 {
		t.Errorf("unexpected errors: %v", errors)
	}

	// Compare result to expected output
	if !deepEqual(t, result, expected) {
		t.Errorf("result does not match expected.\nGot: %#v\nExpected: %#v", result, expected)
	}
}

// deepEqual is a helper for comparing parser outputs in tests.
// It recursively compares maps, slices, and primitive types for equality.
func deepEqual(t *testing.T, a, b interface{}) bool {
	t.Helper()
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	switch aVal := a.(type) {
	case map[string]interface{}:
		bVal := b.(map[string]interface{})
		if len(aVal) != len(bVal) {
			return false
		}
		for k, v := range aVal {
			if !deepEqual(t, v, bVal[k]) {
				return false
			}
		}
		return true
	case []interface{}:
		bVal, ok := b.([]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for i := range aVal {
			if !deepEqual(t, aVal[i], bVal[i]) {
				return false
			}
		}
		return true
	// Handle slices of maps by converting both to []interface{} and comparing recursively
	case []map[string]interface{}:
		// Convert both a and b to []interface{}
		ai := make([]interface{}, len(aVal))
		for i := range aVal {
			ai[i] = aVal[i]
		}
		bi := make([]interface{}, len(b.([]map[string]interface{})))
		for i := range b.([]map[string]interface{}) {
			bi[i] = b.([]map[string]interface{})[i]
		}
		return deepEqual(t, ai, bi)
	default:
		return reflect.DeepEqual(a, b)
	}
}

// TestMixedCaseMultiline checks handling of mixed case labels and multiline values.
func TestMixedCaseMultiline(t *testing.T) {
	input, err := os.ReadFile("../test-assets/mixed_case_multiline_input.txt")
	if err != nil {
		t.Fatalf("failed to read input asset: %v", err)
	}

	expectedBytes, err := os.ReadFile("../test-assets/mixed_case_multiline_output.json")
	if err != nil {
		t.Fatalf("failed to read output asset: %v", err)
	}
	var expected map[string]interface{}
	if err := json.Unmarshal(expectedBytes, &expected); err != nil {
		t.Fatalf("failed to unmarshal expected output: %v", err)
	}

	labels := []Label{
		{Name: "Context"}, {Name: "Intention"}, {Name: "Role"}, {Name: "Action"}, {Name: "Outcome"}, {Name: "Notes"},
	}
	parser, err := NewParser(labels, nil)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	result, errors := parser.Parse(string(input))
	if len(errors) > 0 {
		t.Errorf("unexpected errors: %v", errors)
	}
	if !deepEqual(t, result, expected) {
		t.Errorf("result mismatch.\nGot: %#v\nExpected: %#v", result, expected)
	}
}

// TestJSONAndMalformed checks JSON and malformed JSON parsing and error reporting.
func TestJSONAndMalformed(t *testing.T) {
	input, err := os.ReadFile("../test-assets/json_and_malformed_input.txt")
	if err != nil {
		t.Fatalf("failed to read input asset: %v", err)
	}

	expectedBytes, err := os.ReadFile("../test-assets/json_and_malformed_output.json")
	if err != nil {
		t.Fatalf("failed to read output asset: %v", err)
	}

	errorsBytes, err := os.ReadFile("../test-assets/json_and_malformed_errors.json")
	if err != nil {
		t.Fatalf("failed to read errors asset: %v", err)
	}

	var expected map[string]interface{}
	if err := json.Unmarshal(expectedBytes, &expected); err != nil {
		t.Fatalf("failed to unmarshal expected output: %v", err)
	}

	var expectedErrors []string
	if err := json.Unmarshal(errorsBytes, &expectedErrors); err != nil {
		t.Fatalf("failed to unmarshal expected errors: %v", err)
	}

	labels := []Label{
		{Name: "Config", IsJSON: true}, {Name: "Data", IsJSON: true}, {Name: "Description"},
	}
	parser, err := NewParser(labels, nil)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	result, errors := parser.Parse(string(input))
	if !deepEqual(t, result, expected) {
		t.Errorf("result mismatch.\nGot: %#v\nExpected: %#v", result, expected)
	}
	if len(errors) != len(expectedErrors) || (len(errors) > 0 && errors[0] != expectedErrors[0]) {
		t.Errorf("error mismatch.\nGot: %#v\nExpected: %#v", errors, expectedErrors)
	}
}

// TestRequiredDependency checks required and dependency validation.
func TestRequiredDependency(t *testing.T) {
	input, err := os.ReadFile("../test-assets/required_dependency_input.txt")
	if err != nil {
		t.Fatalf("failed to read input asset: %v", err)
	}

	expectedBytes, err := os.ReadFile("../test-assets/required_dependency_output.json")
	if err != nil {
		t.Fatalf("failed to read output asset: %v", err)
	}

	errorsBytes, err := os.ReadFile("../test-assets/required_dependency_errors.json")
	if err != nil {
		t.Fatalf("failed to read errors asset: %v", err)
	}

	var expected map[string]interface{}
	if err := json.Unmarshal(expectedBytes, &expected); err != nil {
		t.Fatalf("failed to unmarshal expected output: %v", err)
	}

	var expectedErrors []string
	if err := json.Unmarshal(errorsBytes, &expectedErrors); err != nil {
		t.Fatalf("failed to unmarshal expected errors: %v", err)
	}

	labels := []Label{
		{Name: "FieldA"}, {Name: "FieldB", RequiredWith: []string{"FieldA"}},
	}
	parser, err := NewParser(labels, nil)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	result, errors := parser.Parse(string(input))
	if !deepEqual(t, result, expected) {
		t.Errorf("result mismatch.\nGot: %#v\nExpected: %#v", result, expected)
	}
	if len(errors) != len(expectedErrors) || (len(errors) > 0 && errors[0] != expectedErrors[0]) {
		t.Errorf("error mismatch.\nGot: %#v\nExpected: %#v", errors, expectedErrors)
	}
}

// TestBlockParsing checks block parsing with multiple blocks.
func TestBlockParsing(t *testing.T) {
	input, err := os.ReadFile("../test-assets/block_parsing_input.txt")
	if err != nil {
		t.Fatalf("failed to read input asset: %v", err)
	}

	expectedBytes, err := os.ReadFile("../test-assets/block_parsing_output.json")
	if err != nil {
		t.Fatalf("failed to read output asset: %v", err)
	}

	var expected []map[string]interface{}
	if err := json.Unmarshal(expectedBytes, &expected); err != nil {
		t.Fatalf("failed to unmarshal expected output: %v", err)
	}

	labels := []Label{
		{Name: "Task", IsBlockStart: true}, {Name: "Input", IsJSON: true}, {Name: "Result"},
	}
	parser, err := NewParser(labels, nil)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	blocks, errors := parser.ParseBlocks(string(input))
	if len(errors) > 0 {
		t.Errorf("unexpected errors: %v", errors)
	}
	if !deepEqual(t, blocks, expected) {
		t.Errorf("block result mismatch.\nGot: %#v\nExpected: %#v", blocks, expected)
	}
}

// TestNewParserNonMutation verifies that NewParser does not modify the input labels slice.
func TestNewParserNonMutation(t *testing.T) {
	originalLabels := []Label{
		{Name: "Reason", Required: true},
		{Name: "Function", IsJSON: true},
		{Name: "Parameters"},
	}

	// Make a copy to compare against
	labelsCopy := make([]Label, len(originalLabels))
	copy(labelsCopy, originalLabels)

	parser, err := NewParser(originalLabels, nil)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	// Verify original labels are unchanged
	if len(originalLabels) != len(labelsCopy) {
		t.Fatalf("label slice length changed: got %d, want %d", len(originalLabels), len(labelsCopy))
	}

	for i := range originalLabels {
		if originalLabels[i].Name != labelsCopy[i].Name {
			t.Errorf("label[%d].Name was mutated: got %q, want %q", i, originalLabels[i].Name, labelsCopy[i].Name)
		}
		if originalLabels[i].Required != labelsCopy[i].Required {
			t.Errorf("label[%d].Required was mutated: got %v, want %v", i, originalLabels[i].Required, labelsCopy[i].Required)
		}
		if originalLabels[i].IsJSON != labelsCopy[i].IsJSON {
			t.Errorf("label[%d].IsJSON was mutated: got %v, want %v", i, originalLabels[i].IsJSON, labelsCopy[i].IsJSON)
		}
	}

	// Verify parser works correctly with case-preserving keys
	result, errs := parser.Parse("Reason: test\nFunction: {\"key\": \"value\"}\nParameters: param1")
	if len(errs) > 0 {
		t.Errorf("unexpected errors: %v", errs)
	}

	// Check that result keys use original casing
	if _, ok := result["Reason"]; !ok {
		t.Errorf("expected result key 'Reason', got keys: %v", result)
	}
	if _, ok := result["Function"]; !ok {
		t.Errorf("expected result key 'Function', got keys: %v", result)
	}
	if _, ok := result["Parameters"]; !ok {
		t.Errorf("expected result key 'Parameters', got keys: %v", result)
	}

	// Verify lowercase keys are NOT present
	if _, ok := result["reason"]; ok {
		t.Errorf("unexpected lowercase key 'reason' in result: %v", result)
	}
}

// TestCustomSeparators verifies that custom separators work correctly.
func TestCustomSeparators(t *testing.T) {
	labels := []Label{
		{Name: "Key"},
		{Name: "Value"},
	}

	// Test with custom separators (only equals sign)
	parser, err := NewParser(labels, &ParserOptions{Separators: "="})
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	text := "Key=test value\nValue=another value"
	result, errs := parser.Parse(text)
	if len(errs) > 0 {
		t.Errorf("unexpected errors: %v", errs)
	}

	if result["Key"] != "test value" {
		t.Errorf("expected Key='test value', got %v", result["Key"])
	}
	if result["Value"] != "another value" {
		t.Errorf("expected Value='another value', got %v", result["Value"])
	}

	// Test that colon doesn't work with custom separator
	text2 := "Key: should not match"
	result2, errs2 := parser.Parse(text2)
	if len(errs2) > 0 {
		t.Errorf("unexpected errors: %v", errs2)
	}
	// Key should be empty since colon doesn't match
	if result2["Key"] != "" {
		t.Errorf("expected Key to be empty (colon doesn't match), got %v", result2["Key"])
	}
}

// TestRequiredWithFix verifies that RequiredWith only fires when the label is actually present.
func TestRequiredWithFix(t *testing.T) {
	labels := []Label{
		{Name: "FieldA"},
		{Name: "FieldB", RequiredWith: []string{"FieldA"}},
	}

	parser, err := NewParser(labels, nil)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	// Test case: FieldB is missing entirely - should NOT trigger RequiredWith error
	text := "FieldA: valueA"
	result, errs := parser.Parse(text)
	if len(errs) > 0 {
		t.Errorf("unexpected errors when FieldB is missing: %v", errs)
	}
	if result["FieldA"] != "valueA" {
		t.Errorf("expected FieldA='valueA', got %v", result["FieldA"])
	}

	// Test case: FieldB is present but FieldA is missing - SHOULD trigger error
	text2 := "FieldB: valueB"
	_, errs2 := parser.Parse(text2)
	if len(errs2) == 0 {
		t.Error("expected error when FieldB is present but FieldA is missing")
	}
	found := false
	for _, e := range errs2 {
		if e == "'FieldB' requires 'FieldA'" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected error message about FieldB requiring FieldA, got: %v", errs2)
	}
}

// TestEqualsSeparator verifies that equals sign works as a separator.
func TestEqualsSeparator(t *testing.T) {
	labels := []Label{
		{Name: "Name"},
		{Name: "Age"},
	}

	parser, err := NewParser(labels, nil) // Default includes "="
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	text := "Name=John\nAge=30"
	result, errs := parser.Parse(text)
	if len(errs) > 0 {
		t.Errorf("unexpected errors: %v", errs)
	}

	if result["Name"] != "John" {
		t.Errorf("expected Name='John', got %v", result["Name"])
	}
	if result["Age"] != "30" {
		t.Errorf("expected Age='30', got %v", result["Age"])
	}
}
