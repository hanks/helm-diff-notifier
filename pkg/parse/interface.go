package parse

// DesensitizationPattern represents the rules data used in desensitization process.
type DesensitizationPattern struct {
	Name        string
	Target      string
	Submatch    string
	Replacement string
}

// Parser is the interface contains the standard methods for parsing helm diff result.
type Parser interface {
	Desensitize(msg string, rules []DesensitizationPattern) string
	GetDiffs(msg string) []string
	Format(diffs []string) string
}
