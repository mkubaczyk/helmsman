package app

import (
	"net/url"
	"os"

	"github.com/apsdehal/go-logger"
)

type Logger struct {
	*logger.Logger
	SlackWebhook   string
	MSTeamsWebhook string
}

func (l *Logger) Debug(message string) {
	if flags.debug {
		l.Logger.Debug(message)
	}
}

func (l *Logger) Verbose(message string) {
	if flags.verbose {
		l.Logger.Info(message)
	}
}

func (l *Logger) Error(message string) {
	l.notifyAboutFailureUsingWebhooks(message)
	l.Logger.Error(message)
}

func (l *Logger) Fatal(message string) {
	l.notifyAboutFailureUsingWebhooks(message)
	l.Logger.Fatal(message)
}

func (l *Logger) notifyAboutFailureUsingWebhooks(message string) {
	if _, err := url.ParseRequestURI(l.SlackWebhook); err == nil {
		notifySlack(message, l.SlackWebhook, true, flags.apply)
	}
	if _, err := url.ParseRequestURI(l.MSTeamsWebhook); err == nil {
		notifyMSTeams(message, l.MSTeamsWebhook, true, flags.apply)
	}
}

func initLogs(verbose bool, noColors bool) {
	logger.SetDefaultFormat("%{time:2006-01-02 15:04:05} %{level}: %{message}")
	logLevel := logger.InfoLevel
	if verbose {
		logLevel = logger.DebugLevel
	}
	colors := 1
	if noColors {
		colors = 0
	}
	log.Logger, _ = logger.New("logger", colors, os.Stdout, logLevel)
}
