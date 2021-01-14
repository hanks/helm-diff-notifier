package notifier

import (
	"bytes"
	"context"
	"text/template"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

// GHNotifier is used to call the GitHub API to comment on PR review.
type GHNotifier struct {
	Client *github.IssuesService
}

// hasDiffTMPL is the template for diff result.
const hasDiffTMPL = "⚠️ Detected helm diff:\n\n```diff\n{{ .Body }}\n```\n"

// noDiffMSG is the template for no diff result.
const noDiffMSG = "✅ No helm diff detected.\n"

// NewGHNotifier returns an object containing the authorized  client.
func NewGHNotifier(token string) Notifier {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return GHNotifier{Client: client.Issues}
}

// Comment writes the diff result to the PR review page.
func (notifier GHNotifier) Comment(owner string, repo string, number int, msg string) error {
	var output string
	var err error

	if msg != "" {
		output, err = notifier.renderCommentMSG(hasDiffTMPL, msg)
	} else {
		output = noDiffMSG
	}

	if err != nil {
		return err
	}

	issueComment := &github.IssueComment{Body: &output}
	ctx := context.Background()
	_, _, err = notifier.Client.CreateComment(ctx, owner, repo, number, issueComment)

	return err
}

// renderCommentMSG returns the rendered message by specified template.
func (notifier GHNotifier) renderCommentMSG(tmpl string, msg string) (string, error) {
	t := template.Must(template.New("comment").Parse(tmpl))
	var tpl bytes.Buffer
	data := struct {
		Body string
	}{
		Body: msg,
	}
	err := t.Execute(&tpl, data)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}
