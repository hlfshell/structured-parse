package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	sp "github.com/hlfshell/structured-parse/go/structuredparse"
)

// WasmResponse represents the standard response structure for all WASM functions.
type WasmResponse struct {
	Ok     bool                   `json:"ok"`
	Result interface{}            `json:"result,omitempty"`
	Errors []string               `json:"errors,omitempty"`
	Error  string                 `json:"error,omitempty"`
}

// Request represents a unified request structure.
type Request struct {
	Command string                 `json:"command"` // "parse", "parseBlocks", or "version"
	Labels  []LabelJSON            `json:"labels,omitempty"`
	Options *ParserOptionsJSON     `json:"options,omitempty"`
	Text    string                 `json:"text,omitempty"`
}

// LabelJSON represents a label in JSON format.
type LabelJSON struct {
	Name         string   `json:"name"`
	Required     bool     `json:"required,omitempty"`
	RequiredWith []string `json:"requiredWith,omitempty"`
	IsJSON       bool     `json:"isJson,omitempty"`
	IsBlockStart bool     `json:"isBlockStart,omitempty"`
}

// ParserOptionsJSON represents parser options in JSON format.
type ParserOptionsJSON struct {
	Separators string `json:"separators,omitempty"`
}

func main() {
	// Read JSON from stdin
	inputJSON, err := io.ReadAll(os.Stdin)
	if err != nil {
		writeError("failed to read input: " + err.Error())
		return
	}

	var req Request
	if err := json.Unmarshal(inputJSON, &req); err != nil {
		writeError("failed to parse request JSON: " + err.Error())
		return
	}

	switch req.Command {
	case "parse":
		handleParse(req)
	case "parseBlocks":
		handleParseBlocks(req)
	case "version":
		handleVersion()
	default:
		writeError("unknown command: " + req.Command)
	}
}

func handleParse(req Request) {
	labels := convertLabelsFromJSON(req.Labels)
	opts := convertOptionsFromJSON(req.Options)

	parser, err := sp.NewParser(labels, opts)
	if err != nil {
		writeError("failed to create parser: " + err.Error())
		return
	}

	result, errors := parser.Parse(req.Text)

	response := WasmResponse{
		Ok:     len(errors) == 0,
		Result: result,
		Errors: errors,
	}

	writeResponse(response)
}

func handleParseBlocks(req Request) {
	labels := convertLabelsFromJSON(req.Labels)
	opts := convertOptionsFromJSON(req.Options)

	parser, err := sp.NewParser(labels, opts)
	if err != nil {
		writeError("failed to create parser: " + err.Error())
		return
	}

	blocks, errors := parser.ParseBlocks(req.Text)

	response := WasmResponse{
		Ok:     len(errors) == 0,
		Result: blocks,
		Errors: errors,
	}

	writeResponse(response)
}

func handleVersion() {
	response := WasmResponse{
		Ok:     true,
		Result: "1.0.0",
	}
	writeResponse(response)
}

func convertLabelsFromJSON(jsonLabels []LabelJSON) []sp.Label {
	labels := make([]sp.Label, len(jsonLabels))
	for i, jl := range jsonLabels {
		labels[i] = sp.Label{
			Name:         jl.Name,
			Required:     jl.Required,
			RequiredWith: jl.RequiredWith,
			IsJSON:       jl.IsJSON,
			IsBlockStart: jl.IsBlockStart,
		}
	}
	return labels
}

func convertOptionsFromJSON(jsonOpts *ParserOptionsJSON) *sp.ParserOptions {
	if jsonOpts == nil {
		return nil
	}
	return &sp.ParserOptions{
		Separators: jsonOpts.Separators,
	}
}

func writeResponse(response WasmResponse) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		writeError("failed to marshal response: " + err.Error())
		return
	}
	fmt.Println(string(responseJSON))
}

func writeError(errMsg string) {
	response := WasmResponse{
		Ok:    false,
		Error: errMsg,
	}
	responseJSON, _ := json.Marshal(response)
	fmt.Println(string(responseJSON))
}

