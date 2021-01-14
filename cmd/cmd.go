package cmd

import (
	"github.com/hanks/helm-diff-notifier/config"
	"github.com/hanks/helm-diff-notifier/pkg/log"
	"github.com/hanks/helm-diff-notifier/pkg/notifier"
	"github.com/hanks/helm-diff-notifier/pkg/parse"
	"github.com/hanks/helm-diff-notifier/pkg/pullrequest"
)

// Notify is the default command to do the diff notification task.
func Notify(token string, output string) error {
	pr, err := pullrequest.NewFromCircleCIEnv()
	if err != nil {
		return err
	}

	parser := parse.NewDefaultParser()
	diffs := parser.GetDiffs(output)
	log.Logger.Debugf("diffs: %v\n", diffs)
	formatted := parser.Format(diffs)
	cleanMsg := parser.Desensitize(formatted, config.DesensitizationRules)
	log.Logger.Debugf("cleanMsg: %s\n", cleanMsg)

	notifier := notifier.NewGHNotifier(token)
	err = notifier.Comment(pr.Owner, pr.Repo, pr.Number, cleanMsg)

	return err
}
