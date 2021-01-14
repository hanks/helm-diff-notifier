package pullrequest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFromCircleCIEnv(t *testing.T) {
	testCases := []struct {
		title string
		input struct {
			username    string
			reponame    string
			pullRequest string
		}
		expected GHPullRequest
	}{
		{
			title: "Should return GHPullRequest correctly with valid env vars",
			input: struct {
				username    string
				reponame    string
				pullRequest string
			}{
				username:    "foo",
				reponame:    "bar",
				pullRequest: "https://github.com/foo/bar/pull/8378",
			},
			expected: GHPullRequest{"foo", "bar", 8378},
		},
		{
			title: "Should return error with empty env vars",
			input: struct {
				username    string
				reponame    string
				pullRequest string
			}{
				username:    "",
				reponame:    "bar",
				pullRequest: "https://github.com/foo/bar/pull/8378",
			},
			expected: GHPullRequest{},
		},
		{
			title: "Should return error with invalid pull request url",
			input: struct {
				username    string
				reponame    string
				pullRequest string
			}{
				username:    "foo",
				reponame:    "bar",
				pullRequest: "https://github.com/foo/bar/pull/ab78",
			},
			expected: GHPullRequest{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			os.Setenv("CIRCLE_PROJECT_USERNAME", tc.input.username)
			os.Setenv("CIRCLE_PROJECT_REPONAME", tc.input.reponame)
			os.Setenv("CIRCLE_PULL_REQUEST", tc.input.pullRequest)

			actual, _ := NewFromCircleCIEnv()
			assert.Equal(t, tc.expected, actual)

			os.Unsetenv("CIRCLE_PROJECT_USERNAME")
			os.Unsetenv("CIRCLE_PROJECT_REPONAME")
			os.Unsetenv("CIRCLE_PULL_REQUEST")
		})

	}
}
