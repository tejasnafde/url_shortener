package logger

import (
	"os"
	"runtime"
	"strings"
	"strconv"
	"github.com/sirupsen/logrus"
	_ "golang.org/x/text"
)

// TODO: Configure Logrus logger with desired settings
//
// What to implement:
// 1. Set log level (Debug, Info, Warn, Error) -- Done??
// 2. Enable caller reporting (file name and line number)
// 3. Use text formatter (not JSON) for readable logs -- Done
// 4. Add custom fields for context (request_id, user_id, etc.) -- TODO
// 5. Create different loggers for different components (DB, HTTP, Service) -- PARKED for now
//
// Good practices to follow?
// - Structured logging with fields: log.WithField("component", "database").Info("message")
// - Contextual logging: log.WithFields(logrus.Fields{"user_id": 123, "action": "create_url"})
// - Different log levels for different environments (dev vs prod)
//
// Example usage:
// logger.Info("Application started")
// logger.WithField("url", "https://example.com").Info("URL shortened")
// logger.WithError(err).Error("Database connection failed")

var Logger = logrus.New()

// TODO: Initialize logger with proper configuration
func Init(level string) {
	lvl, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		lvl = logrus.InfoLevel
	}
	Logger.SetLevel(lvl)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: "02-01-2006 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			return "", f.File + ":" + strconv.Itoa(f.Line)
		},
	})
	Logger.SetReportCaller(true)
	Logger.SetOutput(os.Stdout)
	Logger.Info("Test log print")
}

// WithComponent adds a component field for structured logs
func WithComponent(name string) *logrus.Entry {
	return Logger.WithField("identifier", name)
}
