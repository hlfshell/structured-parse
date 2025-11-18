package structuredparse

import (
	"strings"
)

// validateDependencies checks required and required_with constraints.
func (p *Parser) validateDependencies(data map[string][]string) []string {
	errList := []string{}
	for _, label := range p.labels {
		key := label.Name
		entries, present := data[key]
		missing := !present || len(entries) == 0 || (len(entries) == 1 && entries[0] == "")
		
		originalName := p.originalNames[key]
		if originalName == "" {
			originalName = key
		}

		if label.Required && missing {
			errList = append(errList, "'"+originalName+"' is required")
		}
		if len(label.RequiredWith) > 0 {
			for _, dep := range label.RequiredWith {
				depKey := strings.ToLower(dep)
				depEntries, depPresent := data[depKey]
				depMissing := !depPresent || len(depEntries) == 0 || (len(depEntries) == 1 && depEntries[0] == "")
				if !missing {
					if depMissing {
						depOriginalName := p.originalNames[depKey]
						if depOriginalName == "" {
							depOriginalName = dep
						}
						errList = append(errList, "'"+originalName+"' requires '"+depOriginalName+"'")
					}
				}
			}
		}
	}
	return errList
}

