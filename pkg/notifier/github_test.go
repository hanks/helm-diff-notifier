package notifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderCommentMSG(t *testing.T) {
	testCases := []struct {
		title string
		input struct {
			tmpl string
			msg  string
		}
		expected string
	}{
		{
			title: "Render HasDiffTMPL with valid message should be ok",
			input: struct {
				tmpl string
				msg  string
			}{
				tmpl: hasDiffTMPL,
				msg:  "hello world",
			},
			expected: "⚠️ Detected helm diff:\n\n```diff\nhello world\n```\n",
		},
		{
			title: "Render HasDiffTMPL with empty message should be ok",
			input: struct {
				tmpl string
				msg  string
			}{
				tmpl: hasDiffTMPL,
				msg:  "",
			},
			expected: "⚠️ Detected helm diff:\n\n```diff\n\n```\n",
		},
		{
			title: "Render template with wrong key name should be failed",
			input: struct {
				tmpl string
				msg  string
			}{
				tmpl: "⚠️ Detected helm diff:\n\n```diff\n{{ .NoBody }}\n```\n",
				msg:  "",
			},
			expected: "can't evaluate field",
		},
	}

	notifier := GHNotifier{}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual, err := notifier.renderCommentMSG(tc.input.tmpl, tc.input.msg)
			if err != nil {
				assert.Contains(t, err.Error(), tc.expected)
			} else {
				assert.Equal(t, tc.expected, actual)
			}

		})
	}
}
