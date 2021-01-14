package parse

import (
	"regexp"
	"strings"
)

// DefaultParser represents the default parser.
type DefaultParser struct {
}

// NewDefaultParser returns a new parser object.
func NewDefaultParser() Parser {
	return DefaultParser{}
}

// Desensitize is doing secret mask to the original helm diff result by using the rules.
func (p DefaultParser) Desensitize(msg string, rules []DesensitizationPattern) string {
	desensitizedMSG := msg

	for _, rule := range rules {
		pattern := regexp.MustCompile(rule.Target)
		submatchPattern := regexp.MustCompile(rule.Submatch)
		desensitizedMSG = pattern.ReplaceAllStringFunc(desensitizedMSG, func(m string) string {
			return submatchPattern.ReplaceAllString(m, rule.Replacement)
		})
	}

	return desensitizedMSG
}

// GetDiffs returns a list of diff sections from the original result, that could be helpful for the latter processing.
func (p DefaultParser) GetDiffs(msg string) []string {
	diffs := p.findDiffs(msg)
	keyDiffs := p.simplifyDiffs(diffs)

	return keyDiffs
}

// simplifyDiffs is picking the key diff changes in the original result
func (p DefaultParser) simplifyDiffs(diffs []string) []string {
	var result []string

	for _, diff := range diffs {
		lines := strings.Split(diff, "\n")
		targetLines := []string{lines[0]}
		for i, line := range lines {
			if i == 1 && strings.HasPrefix(line, " ") {
				targetLines = append(targetLines, line)
				continue
			}

			if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") {
				targetLines = append(targetLines, line)
			}
		}
		result = append(result, strings.Join(targetLines[:], "\n"))
	}

	return result
}

// findDiffs returns a list of diff result section by diff headers.
func (p DefaultParser) findDiffs(msg string) []string {
	titlePattern := `(?m)^\w+.*?:$`
	re := regexp.MustCompile(titlePattern)
	titles := re.FindAllString(msg, -1)
	diffs := removeEmpty(re.Split(msg, -1))

	for i, diff := range diffs {
		diffs[i] = titles[i] + diff
	}

	return diffs
}

// removeEmpty is a helper function to remove empty items in a string list
func removeEmpty(s []string) []string {
	var result = []string{}
	for _, str := range s {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}

// Format returns a formatted string from a list of diff section.
func (p DefaultParser) Format(diffs []string) string {
	return strings.Join(diffs[:], "\n\n")
}
