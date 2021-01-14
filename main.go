package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/hanks/helm-diff-notifier/cmd"
	"github.com/hanks/helm-diff-notifier/pkg/log"
	"github.com/hanks/helm-diff-notifier/version"
)

func execute() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	output := strings.Join(lines, "\n")
	token := os.Getenv("GITHUB_API_TOKEN")
	err := cmd.Notify(token, output)
	if err != nil {
		log.Logger.Fatal(err)
	}
	log.Logger.Info("main - Send notification is succeed")
}

func main() {
	app := cli.NewApp()
	app.Name = "helm-diff-notify"
	app.Description = "Notify helm diff result to github pullrequest"
	app.Usage = "helm diff upgrade -n foo release_bar | helm-diff-notify"
	app.Version = version.Version

	app.Action = func(c *cli.Context) error {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			cli.ShowAppHelpAndExit(c, 0)
		}

		execute()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Logger.Fatal(err)
	}
}
