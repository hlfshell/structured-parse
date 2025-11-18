//go:build js && wasm

package structuredparse

import (
	"encoding/json"
	"syscall/js"
)

// wasmParse is the exported function for parsing text.
// It accepts a JSON string representing ParseRequest and returns a JSON string with WasmResponse.
func wasmParse(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return createErrorResponse("expected 1 argument: request JSON string")
	}

	requestJSON := args[0].String()
	var req ParseRequest
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		return createErrorResponse("failed to parse request JSON: " + err.Error())
	}

	labels := convertLabelsFromJSON(req.Labels)
	opts := convertOptionsFromJSON(req.Options)

	parser, err := NewParser(labels, opts)
	if err != nil {
		return createErrorResponse("failed to create parser: " + err.Error())
	}

	result, errors := parser.Parse(req.Text)

	response := WasmResponse{
		Ok:     len(errors) == 0,
		Result: result,
		Errors: errors,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return createErrorResponse("failed to marshal response: " + err.Error())
	}

	return string(responseJSON)
}

// wasmParseBlocks is the exported function for parsing text into blocks.
// It accepts a JSON string representing ParseBlocksRequest and returns a JSON string with WasmResponse.
func wasmParseBlocks(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return createErrorResponse("expected 1 argument: request JSON string")
	}

	requestJSON := args[0].String()
	var req ParseBlocksRequest
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		return createErrorResponse("failed to parse request JSON: " + err.Error())
	}

	labels := convertLabelsFromJSON(req.Labels)
	opts := convertOptionsFromJSON(req.Options)

	parser, err := NewParser(labels, opts)
	if err != nil {
		return createErrorResponse("failed to create parser: " + err.Error())
	}

	blocks, errors := parser.ParseBlocks(req.Text)

	response := WasmResponse{
		Ok:     len(errors) == 0,
		Result: blocks,
		Errors: errors,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return createErrorResponse("failed to marshal response: " + err.Error())
	}

	return string(responseJSON)
}

// wasmVersion returns the version of the WASM module.
func wasmVersion(this js.Value, args []js.Value) interface{} {
	return "1.0.0"
}

// RegisterWasmFunctions registers all WASM functions to be exported.
func RegisterWasmFunctions() {
	js.Global().Set("wasmParse", js.FuncOf(wasmParse))
	js.Global().Set("wasmParseBlocks", js.FuncOf(wasmParseBlocks))
	js.Global().Set("wasmVersion", js.FuncOf(wasmVersion))
}
