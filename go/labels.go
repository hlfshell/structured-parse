// Package structuredparse provides a small Go library for parsing labeled/structured LLM output
// into a map of fields. It supports block parsing, JSON fields, required field validation,
// and robust handling of LLM output quirks.
package structuredparse

import (
	"errors"
	"regexp"
	"strings"
)

// Label defines a label for parsing with options for required, dependencies, JSON, and block start.
type Label struct {
	Name         string   // Name of the label (case-insensitive matching, but original casing preserved in results)
	Required     bool     // Whether this label is required
	RequiredWith []string // List of other label names required with this one
	IsJSON       bool     // Whether this label should be parsed as JSON
	IsBlockStart bool     // Whether this label starts a new block
}

type labelPattern struct {
	// Name of the label (lowercase for matching)
	Name string
	// Regex pattern for the label
	Pattern *regexp.Regexp
}

// ParserOptions allows customization of parser behavior.
type ParserOptions struct {
	// Separators is a string containing the allowed separator characters.
	// Default is ":~-=" (colon, tilde, dash, equals).
	// Each character in the string is treated as a valid separator.
	Separators string
}

// NewParser creates a new Parser with the given labels and optional options.
// If opts is nil, default options are used (separators: ":~-=").
func NewParser(labels []Label, opts *ParserOptions) (*Parser, error) {
	internalLabels := make([]Label, len(labels))
	copy(internalLabels, labels)

	labelMap := make(map[string]Label)
	originalNames := make(map[string]string)
	blockStartCount := 0

	for i := range internalLabels {
		originalName := internalLabels[i].Name
		lowerName := strings.ToLower(originalName)

		internalLabels[i].Name = lowerName
		labelMap[lowerName] = internalLabels[i]
		originalNames[lowerName] = originalName

		if internalLabels[i].IsBlockStart {
			blockStartCount++
		}
	}

	if blockStartCount > 1 {
		return nil, errors.New("only one block start label is allowed")
	}

	separators := ":~-="
	if opts != nil && opts.Separators != "" {
		separators = opts.Separators
	}

	patterns := buildPatterns(internalLabels, separators)
	separatorRegex := buildSeparatorRegex(separators)

	return &Parser{
		labels:        internalLabels,
		patterns:      patterns,
		labelMap:      labelMap,
		originalNames: originalNames,
		separators:    separators,
		separatorRe:   separatorRegex,
	}, nil
}

// buildPatterns constructs regex patterns for each label.
func buildPatterns(labels []Label, separators string) []labelPattern {
	var patterns []labelPattern
	escapedSeparators := regexp.QuoteMeta(separators)
	escapedSeparators = strings.ReplaceAll(escapedSeparators, `\-`, `-`)
	if strings.Contains(escapedSeparators, "-") {
		escapedSeparators = strings.ReplaceAll(escapedSeparators, "-", "")
		escapedSeparators += "-"
	}

	for _, label := range labels {
		labelRegex := strings.Join(strings.Fields(label.Name), `\s+`)
		pattern := regexp.MustCompile(`(?i)^\s*` + labelRegex + `\s*[` + escapedSeparators + `]+\s*`)
		patterns = append(patterns, labelPattern{Name: label.Name, Pattern: pattern})
	}
	return patterns
}

// buildSeparatorRegex creates a regex for separator matching.
func buildSeparatorRegex(separators string) *regexp.Regexp {
	escapedSeparators := regexp.QuoteMeta(separators)
	escapedSeparators = strings.ReplaceAll(escapedSeparators, `\-`, `-`)
	if strings.Contains(escapedSeparators, "-") {
		escapedSeparators = strings.ReplaceAll(escapedSeparators, "-", "")
		escapedSeparators += "-"
	}
	return regexp.MustCompile(`^\s*[` + escapedSeparators + `]+`)
}
