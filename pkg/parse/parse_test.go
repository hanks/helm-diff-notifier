package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestRemoveEmpty(t *testing.T) {
	testCases := []struct {
		title    string
		input    []string
		expected []string
	}{
		{
			title:    "Should remove one empty string with valid input",
			input:    []string{"a", "b", "", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			title:    "Should remove multi empty strings with valid input",
			input:    []string{"a", "b", "", "c", ""},
			expected: []string{"a", "b", "c"},
		},
		{
			title:    "Should work with empty input",
			input:    []string{},
			expected: []string{},
		},
		{
			title:    "Should work with empty input",
			input:    []string{""},
			expected: []string{},
		},
		{
			title:    "Should work with empty input",
			input:    []string{"", ""},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual := removeEmpty(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

type ParseTestSuite struct {
	suite.Suite
	parser DefaultParser
}

func (suite *ParseTestSuite) SetupTest() {
	suite.parser = DefaultParser{}
}

func (suite *ParseTestSuite) TestFormat() {
	testCases := []struct {
		title    string
		input    []string
		expected string
	}{
		{
			title:    "Should format correctly with valid inputs",
			input:    []string{"a", "b", "c"},
			expected: "a\n\nb\n\nc",
		},
		{
			title:    "Should format correctly with empty inputs",
			input:    []string{},
			expected: "",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.title, func(t *testing.T) {
			actual := suite.parser.Format(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func (suite *ParseTestSuite) TestFindDiffs() {
	testCases := []struct {
		title    string
		input    string
		expected []string
	}{
		{
			title: "Should get diffs correctly with build-in pattern",
			input: `default, foo, ConfigMap (v1) has changed:
  # Source: foo/templates/configmap.yaml
  apiVersion: v1
  kind: ConfigMap
default, foo, Secret (v1) has been removed:
- # Source: foo/templates/secret.yaml
- apiVersion: v1
- kind: Secret
default, foo, StatefulSet (apps) has been added:
+ # Source: foo/templates/statefulset.yaml
+ apiVersion: apps/v1
+ kind: StatefulSet
`,
			expected: []string{
				`default, foo, ConfigMap (v1) has changed:
  # Source: foo/templates/configmap.yaml
  apiVersion: v1
  kind: ConfigMap
`,
				`default, foo, Secret (v1) has been removed:
- # Source: foo/templates/secret.yaml
- apiVersion: v1
- kind: Secret
`,
				`default, foo, StatefulSet (apps) has been added:
+ # Source: foo/templates/statefulset.yaml
+ apiVersion: apps/v1
+ kind: StatefulSet
`,
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.title, func(t *testing.T) {
			actual := suite.parser.findDiffs(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func (suite *ParseTestSuite) TestGetDiffs() {
	testCases := []struct {
		title    string
		input    string
		expected []string
	}{
		{
			title: "Should get diffs correctly",
			input: `default, foo, ConfigMap (v1) has changed:
  # Source: foo/templates/configmap.yaml
  apiVersion: v1
+ kind: ConfigMap
- metadata:
    name: graylog
default, foo, Secret (v1) has been removed:
- # Source: foo/templates/secret.yaml
- apiVersion: v1
- kind: Secret
default, foo, StatefulSet (apps) has been added:
+ # Source: foo/templates/statefulset.yaml
+ apiVersion: apps/v1
+ kind: StatefulSet
`,
			expected: []string{
				`default, foo, ConfigMap (v1) has changed:
  # Source: foo/templates/configmap.yaml
+ kind: ConfigMap
- metadata:`,
				`default, foo, Secret (v1) has been removed:
- # Source: foo/templates/secret.yaml
- apiVersion: v1
- kind: Secret`,
				`default, foo, StatefulSet (apps) has been added:
+ # Source: foo/templates/statefulset.yaml
+ apiVersion: apps/v1
+ kind: StatefulSet`,
			},
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.title, func(t *testing.T) {
			actual := suite.parser.GetDiffs(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func (suite *ParseTestSuite) TestDesensitize() {
	testCases := []struct {
		title string
		input struct {
			msg   string
			rules []DesensitizationPattern
		}
		expected string
	}{
		{
			title: "Should get diffs correctly",
			input: struct {
				msg   string
				rules []DesensitizationPattern
			}{
				msg: `default, foo, ConfigMap (v1) has changed:
  # Source: foo/templates/configmap.yaml
  apiVersion: v1
+ kind: ConfigMap
- metadata:
	name: graylog
-     mongodb_uri = mongodb://user:123456@mongodb.sre.svc.cluster.local/bar
+     mongodb_uri = ${MONGODB_URI}
`,
				rules: []DesensitizationPattern{
					{
						Name:        "MongoDB URI",
						Target:      `mongodb://.*`,
						Submatch:    `(:)\w+(@)`,
						Replacement: `${1}****${2}`,
					},
				},
			},
			expected: `default, foo, ConfigMap (v1) has changed:
  # Source: foo/templates/configmap.yaml
  apiVersion: v1
+ kind: ConfigMap
- metadata:
	name: graylog
-     mongodb_uri = mongodb://user:****@mongodb.sre.svc.cluster.local/bar
+     mongodb_uri = ${MONGODB_URI}
`,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.title, func(t *testing.T) {
			actual := suite.parser.Desensitize(tc.input.msg, tc.input.rules)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func (suite *ParseTestSuite) TestsimplifyDiffs() {
	testCases := []struct {
		title    string
		input    []string
		expected []string
	}{
		{
			title: "Should simplify diffs correctly",
			input: []string{
				`default, foo, ConfigMap (v1) has changed:
  # Source: foo/templates/configmap.yaml
  apiVersion: v1
+ kind: ConfigMap
- metadata:
  name: foo
`,
				`default, foo, Secret (v1) has been removed:
- # Source: foo/templates/secret.yaml
- apiVersion: v1
- kind: Secret
`,
				`default, foo, StatefulSet (apps) has been added:
+ # Source: foo/templates/statefulset.yaml
+ apiVersion: apps/v1
+ kind: StatefulSet
`,
			},
			expected: []string{
				`default, foo, ConfigMap (v1) has changed:
  # Source: foo/templates/configmap.yaml
+ kind: ConfigMap
- metadata:`,
				`default, foo, Secret (v1) has been removed:
- # Source: foo/templates/secret.yaml
- apiVersion: v1
- kind: Secret`,
				`default, foo, StatefulSet (apps) has been added:
+ # Source: foo/templates/statefulset.yaml
+ apiVersion: apps/v1
+ kind: StatefulSet`,
			},
		},
	}
	for _, tc := range testCases {
		suite.T().Run(tc.title, func(t *testing.T) {
			actual := suite.parser.simplifyDiffs(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestParseTestSuite(t *testing.T) {
	suite.Run(t, new(ParseTestSuite))
}
