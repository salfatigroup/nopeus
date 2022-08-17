package logger

import (
	"context"
	"os"

	"github.com/gofrs/uuid"
	"github.com/salfatigroup/gologsnag"
	"github.com/sirupsen/logrus"
)

var (
	log     *logrus.Entry
	logsnag *gologsnag.LogSnag
)

// define a global session id to identify the flow
var sessionID = uuid.Must(uuid.NewV4()).String()

// initialize the logger with the reuiqred
// log level and format
func init() {
	initLogrusLogger()
	initLogsnagLogger()
}

// initialize the logsnag logger
func initLogsnagLogger() {
	logsnag = gologsnag.NewLogSnag(
		// public nopeus logsnag key
		"2f0420e7710703268ea2ab32f493c887",
		"salfati-group-cloud",
	)
}

// initialize the logrus logger
func initLogrusLogger() {
	// create a new logger
	l := logrus.New()

	// get the log level from the environment variable
	// and set the log level
	l.SetLevel(getLogLevel())

	// set the format of the logger
	// base on the GO_ENV environment variable
	l.SetFormatter(getLogFormat())

	// define the default logger fields
	rl := l.WithFields(logrus.Fields{
		"app":        "nopeus",
		"session-id": sessionID,
	})

	// set the logger to the global Logger
	log = rl
}

// get the logger level from the environment variable calld LOG_LEVEL
func getLogLevel() logrus.Level {
	// get the log level from the environment variable
	// and set the log level
	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.ErrorLevel
	}
	return level
}

// get the logger format based on the GO_ENV environment variable
func getLogFormat() logrus.Formatter {
	// get the format from the environment variable
	// and set the format
	format := os.Getenv("GO_ENV")
	if format == "production" {
		return &logrus.JSONFormatter{}
	}
	return &logrus.TextFormatter{}
}

// export additional withfield function
func WithField(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}

// export the different log levels logging functions
// and the format logging functions
func Trace(args ...interface{}) {
	log.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	log.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

type PublishError struct {
	Error error `json:"error"`
}

// export the logsnag publish function
func Publish(options *gologsnag.PublishOptions) *PublishError {
	if os.Getenv("GO_ENV") != "development" {
		ctx := context.Background()

		// if tags do not exists, init empty tags
		if options.Tags == nil {
			options.Tags = &gologsnag.Tags{}
		}

		// extend the tags with the session id
		options.Tags.Add("session-id", sessionID)

		// force the channel to be "nopeus-public"
		options.Channel = "nopeus-public"

		// publish the logsnag message
		err := logsnag.Publish(ctx, options)
		return &PublishError{
			Error: err,
		}
	}

	return nil
}

// export the logsnag insight function
func Insight(options *gologsnag.InsightOptions) *PublishError {
	if os.Getenv("GO_ENV") != "development" {
		ctx := context.Background()

		// send an insight to logsnag
		err := logsnag.Insight(ctx, options)
		return &PublishError{
			Error: err,
		}
	}

	return nil
}
