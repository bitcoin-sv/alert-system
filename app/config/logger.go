package config

import "github.com/ordishs/gocore"

// LoggerInterface is the interface for the logger
// This is used to allow the logger to be mocked and tested
// These methods are the same as the gocore.Logger methods
type LoggerInterface interface {
	Debug(args ...interface{})
	Debugf(msg string, args ...interface{})
	Error(args ...interface{})
	ErrorWithStack(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(msg string, args ...interface{})
	Info(args ...interface{})
	Infof(msg string, args ...interface{})
	LogLevel() int
	Panic(args ...interface{})
	Panicf(msg string, args ...interface{})
	Warn(args ...interface{})
	Warnf(msg string, args ...interface{})
	Printf(format string, v ...interface{}) // Custom method for go-api-router
	// GetLogLevel() gocore.logLevel
}

// ExtendedLogger is the extended logger to satisfy the LoggerInterface
type ExtendedLogger struct {
	*gocore.Logger
}

// Printf will print the log message to the console
func (es *ExtendedLogger) Printf(format string, v ...interface{}) {
	es.Infof(format, v...)
}
