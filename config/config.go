package config

import "github.com/hanks/helm-diff-notifier/pkg/parse"

// DesensitizationRules presents the patterns of desensitization rules.
var DesensitizationRules []parse.DesensitizationPattern = []parse.DesensitizationPattern{
	{
		Name:        "MongoDB URI",
		Target:      `mongodb://.*`,
		Submatch:    `(:)\w+(@)`,
		Replacement: `${1}****${2}`,
	},
}
