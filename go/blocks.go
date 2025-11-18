package structuredparse

import (
	"strings"
)

// ParseBlocks parses the text into blocks, splitting at the block start label.
func (p *Parser) ParseBlocks(text string) ([]map[string]interface{}, []string) {
	blockLabel := ""
	for _, label := range p.labels {
		if label.IsBlockStart {
			blockLabel = label.Name
			break
		}
	}
	if blockLabel == "" {
		return nil, []string{"no block start label defined - must have at least one"}
	}

	cleaned := cleanText(text)
	lines := splitAndTrimLines(cleaned)

	var (
		blocks       [][]string
		currentBlock []string
		inBlock      bool
	)

	for _, line := range lines {
		labelName, _ := p.parseLine(line)
		if strings.ToLower(labelName) == blockLabel {
			if inBlock && len(currentBlock) > 0 {
				blocks = append(blocks, currentBlock)
				currentBlock = []string{}
			}
			inBlock = true
		}
		if inBlock {
			currentBlock = append(currentBlock, line)
		}
	}
	if inBlock && len(currentBlock) > 0 {
		blocks = append(blocks, currentBlock)
	}

	var (
		results []map[string]interface{}
		errList []string
	)
	for _, blockLines := range blocks {
		blockText := strings.Join(blockLines, "\n")
		result, blockErr := p.parseLines(blockText)
		if len(blockErr) > 0 {
			errList = append(errList, blockErr...)
		}
		results = append(results, result)
	}
	return results, errList
}
