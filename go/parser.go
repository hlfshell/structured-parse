package structuredparse

import (
	"encoding/json"
	"regexp"
	"strings"
)

var (
	codeBlockRe  = regexp.MustCompile("(?s)```(?:\\w+)?\\s*(.*?)\\s*```")
	inlineCodeRe = regexp.MustCompile("`([^`]+)`")
)

// Parser parses labeled sections from text input.
type Parser struct {
	labels        []Label           // Internal copy of labels (with lowercase names)
	patterns      []labelPattern    // Regex patterns for label matching
	labelMap      map[string]Label  // Map of lowercase label name -> Label (for lookup)
	originalNames map[string]string // Map of lowercase label name -> original name (for result keys)
	separators    string            // Allowed separator characters
	separatorRe   *regexp.Regexp    // Precompiled regex for separator matching
}

// Parse parses the text into a map of label names (preserving original casing) to their values.
// Each label can have a single value or a slice of values.
//   - Detects labels using regex patterns (case-insensitive, multiple separators)
//   - Collects multi-line values for labels
//   - Parses JSON fields if specified
//   - Validates required fields and dependencies
//   - Returns a map of results and a slice of error strings
func (p *Parser) Parse(text string) (map[string]interface{}, []string) {
	return p.parseLines(cleanText(text))
}

// parseLines parses already-cleaned text that has been split into lines.
// This is used internally to avoid double-cleaning in ParseBlocks.
func (p *Parser) parseLines(text string) (map[string]interface{}, []string) {
	lines := splitAndTrimLines(text)

	data := make(map[string][]string)
	for _, label := range p.labels {
		data[label.Name] = []string{}
	}
	var (
		currentLabel string
		currentEntry strings.Builder
	)

	for _, line := range lines {
		labelName, value := p.parseLine(line)
		if labelName != "" {
			// If we were collecting a previous entry, finalize it
			if currentLabel != "" {
				finalizeEntry(data, currentLabel, currentEntry.String())
				currentEntry.Reset()
			}
			currentLabel = strings.ToLower(labelName)
			currentEntry.WriteString(value)
		} else if currentLabel != "" {
			isLabelLine := p.isLabelLine(line)
			if !isLabelLine {
				if currentEntry.Len() > 0 {
					currentEntry.WriteString("\n")
				}
				currentEntry.WriteString(line)
			}
		}
	}
	if currentLabel != "" {
		finalizeEntry(data, currentLabel, currentEntry.String())
	}

	results, errList := p.processResults(data)
	return results, errList
}

// cleanText removes markdown code blocks and inline code from the input text.
func cleanText(text string) string {
	text = codeBlockRe.ReplaceAllStringFunc(text, func(match string) string {
		sub := codeBlockRe.FindStringSubmatch(match)
		if len(sub) > 1 {
			return sub[1]
		}
		return ""
	})
	text = inlineCodeRe.ReplaceAllString(text, "$1")
	return strings.TrimSpace(text)
}

// isLabelLine checks if a line starts with a known label.
func (p *Parser) isLabelLine(line string) bool {
	trimmed := strings.ToLower(strings.TrimSpace(line))
	for _, lbl := range p.labels {
		lowerName := strings.ToLower(lbl.Name)
		if strings.HasPrefix(trimmed, lowerName) {
			remain := trimmed[len(lowerName):]
			if p.separatorRe.MatchString(remain) {
				return true
			}
		}
	}
	return false
}

// splitAndTrimLines splits text into lines and trims right whitespace.
func splitAndTrimLines(text string) []string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t\r")
	}
	return lines
}

// parseLine tries to match a label at the start of the line.
func (p *Parser) parseLine(line string) (string, string) {
	for _, pat := range p.patterns {
		if loc := pat.Pattern.FindStringIndex(line); loc != nil {
			value := strings.TrimSpace(line[loc[1]:])
			return pat.Name, value
		}
	}
	for labelName := range p.labelMap {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(trimmed), labelName) {
			remain := trimmed[len(labelName):]
			if p.separatorRe.MatchString(remain) {
				content := p.separatorRe.ReplaceAllString(remain, "")
				return labelName, strings.TrimSpace(content)
			}
			return "", trimmed
		}
	}
	return "", ""
}

// finalizeEntry appends a non-empty entry to the data map for a label.
func finalizeEntry(data map[string][]string, labelName, entry string) {
	content := strings.TrimSpace(entry)
	if content != "" {
		data[labelName] = append(data[labelName], content)
	}
}

// processResults parses JSON fields, flattens single-value lists, and collects errors.
// Result map keys use original label names (preserving user's casing).
func (p *Parser) processResults(rawData map[string][]string) (map[string]interface{}, []string) {
	results := make(map[string]interface{})
	errList := []string{}
	for lowerName, entries := range rawData {
		originalName := p.originalNames[lowerName]
		if originalName == "" {
			originalName = lowerName
		}

		labelDef := p.labelMap[lowerName]
		parsedEntries := []interface{}{}
		for _, entry := range entries {
			if labelDef.IsJSON {
				if strings.TrimSpace(entry) == "" {
					parsedEntries = append(parsedEntries, map[string]interface{}{})
					continue
				}
				var obj interface{}
				if err := json.Unmarshal([]byte(entry), &obj); err != nil {
					parsedEntries = append(parsedEntries, entry)
					errList = append(errList, "JSON error in '"+originalName+"': "+err.Error())
				} else {
					parsedEntries = append(parsedEntries, obj)
				}
			} else {
				parsedEntries = append(parsedEntries, entry)
			}
		}
		if len(parsedEntries) == 1 {
			if str, ok := parsedEntries[0].(string); ok && str == "" {
				results[originalName] = ""
			} else {
				results[originalName] = parsedEntries[0]
			}
		} else if len(parsedEntries) == 0 {
			results[originalName] = ""
		} else {
			results[originalName] = parsedEntries
		}
	}
	errList = append(errList, p.validateDependencies(rawData)...)
	return results, errList
}
