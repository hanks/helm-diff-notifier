package log

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Logger is the global log object using cross the project, providing the
// consistent logging style.
var Logger = log.New()

func init() {
	debug := strings.ToLower(os.Getenv("HELM_DIFF_DEBUG"))
	if debug == "true" || debug == "yes" {
		Logger.SetLevel(log.DebugLevel)
	} else {
		Logger.SetLevel(log.InfoLevel)
	}

	Logger.SetReportCaller(true)
	Logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
