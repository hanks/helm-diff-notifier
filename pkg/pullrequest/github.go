package pullrequest

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hanks/helm-diff-notifier/pkg/log"
)

// GHPullRequest represents the data used for pull request in the Github API
type GHPullRequest struct {
	Owner  string
	Repo   string
	Number int
}

// NewFromCircleCIEnv returns a new GHPullRequest object from environment variables in CircleCI pipeline.
func NewFromCircleCIEnv() (GHPullRequest, error) {
	owner := os.Getenv("CIRCLE_PROJECT_USERNAME")
	repo := os.Getenv("CIRCLE_PROJECT_REPONAME")
	pullRequestURL := os.Getenv("CIRCLE_PULL_REQUEST")

	if owner == "" || repo == "" || pullRequestURL == "" {
		return GHPullRequest{}, fmt.Errorf("CircleCI ENV vars should not be empty, "+
			"owner: %s, repo: %s, pr_url: %s", owner, repo, pullRequestURL)
	}

	parts := strings.Split(pullRequestURL, "/")
	numberStr := parts[len(parts)-1]
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return GHPullRequest{}, err
	}
	log.Logger.Infof("owner: %s, repo: %s, number: %d\n", owner, repo, number)

	return GHPullRequest{
		Owner:  owner,
		Repo:   repo,
		Number: number,
	}, nil
}
