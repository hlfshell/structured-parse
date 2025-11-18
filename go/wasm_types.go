package structuredparse

import "encoding/json"

// WasmResponse represents the standard response structure for all WASM functions.
// It contains either a result or an error, along with any parsing errors.
type WasmResponse struct {
	Ok     bool                   `json:"ok"`
	Result interface{}            `json:"result,omitempty"`
	Errors []string               `json:"errors,omitempty"`
	Error  string                 `json:"error,omitempty"` // For system errors
}

// LabelJSON represents a label in JSON format for WASM consumption.
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

// NewParserRequest represents the request to create a new parser.
type NewParserRequest struct {
	Labels  []LabelJSON        `json:"labels"`
	Options *ParserOptionsJSON `json:"options,omitempty"`
}

// ParseRequest represents a request to parse text.
type ParseRequest struct {
	Labels  []LabelJSON        `json:"labels"`
	Options *ParserOptionsJSON `json:"options,omitempty"`
	Text    string             `json:"text"`
}

// ParseBlocksRequest represents a request to parse text into blocks.
type ParseBlocksRequest struct {
	Labels  []LabelJSON        `json:"labels"`
	Options *ParserOptionsJSON `json:"options,omitempty"`
	Text    string             `json:"text"`
}

// convertLabelsFromJSON converts JSON labels to internal Label structs.
func convertLabelsFromJSON(jsonLabels []LabelJSON) []Label {
	labels := make([]Label, len(jsonLabels))
	for i, jl := range jsonLabels {
		labels[i] = Label{
			Name:         jl.Name,
			Required:     jl.Required,
			RequiredWith: jl.RequiredWith,
			IsJSON:       jl.IsJSON,
			IsBlockStart: jl.IsBlockStart,
		}
	}
	return labels
}

// convertOptionsFromJSON converts JSON options to internal ParserOptions.
func convertOptionsFromJSON(jsonOpts *ParserOptionsJSON) *ParserOptions {
	if jsonOpts == nil {
		return nil
	}
	return &ParserOptions{
		Separators: jsonOpts.Separators,
	}
}

// createErrorResponse creates a JSON error response string.
func createErrorResponse(errMsg string) string {
	response := WasmResponse{
		Ok:    false,
		Error: errMsg,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON)
}

