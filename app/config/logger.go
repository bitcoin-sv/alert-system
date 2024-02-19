package config

import (
	"fmt"
	"log"
	"os"
)

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
	LogLevel() string
	Panic(args ...interface{})
	Panicf(msg string, args ...interface{})
	Warn(args ...interface{})
	Warnf(msg string, args ...interface{})
	Printf(format string, v ...interface{}) // Custom method for go-api-router
	CloseWriter() error
	// GetLogLevel() gocore.logLevel
}

// ExtendedLogger is the extended logger to satisfy the LoggerInterface
type ExtendedLogger struct {
	*log.Logger
	logLevel string
	writer   *os.File
}

// CloseWriter close the log writer
func (es *ExtendedLogger) CloseWriter() error {
	return es.writer.Close()
}

// Printf will print the log message to the console
func (es *ExtendedLogger) Printf(format string, v ...interface{}) {
	es.Logger.Printf(format, v...)
}

// Debugf will print debug messages to the console
func (es *ExtendedLogger) Debugf(format string, v ...interface{}) {
	if es.logLevel != "debug" {
		return
	}
	es.Logger.Printf(fmt.Sprintf("\033[1;34m| DEBUG | %s\033[0m", format), v...)
}

// Debug will print debug messages to the console
func (es *ExtendedLogger) Debug(v ...interface{}) {
	if es.logLevel != "debug" {
		return
	}
	es.Logger.Printf("%v", v...)
}

// Error will print debug messages to the console
func (es *ExtendedLogger) Error(v ...interface{}) {
	es.Logger.Printf("%v", v...)
}

// Errorf will print debug messages to the console
func (es *ExtendedLogger) Errorf(format string, v ...interface{}) {
	es.Logger.Printf(fmt.Sprintf("\033[1;31m| ERROR |: %s\033[0m", format), v...)
}

// ErrorWithStack will print debug messages to the console
func (es *ExtendedLogger) ErrorWithStack(format string, v ...interface{}) {
	es.Logger.Printf(format, v...)
}

// Info will print info messages to the console
func (es *ExtendedLogger) Info(v ...interface{}) {
	es.Logger.Printf("%v", v...)
}

// Infof will print info messages to the console
func (es *ExtendedLogger) Infof(format string, v ...interface{}) {
	es.Logger.Printf(fmt.Sprintf("\033[1;32m| INFO  | %s\033[0m", format), v...)
}

// LogLevel returns the logging level
func (es *ExtendedLogger) LogLevel() string {
	return es.logLevel
}

// Warn will print warning messages to the console
func (es *ExtendedLogger) Warn(v ...interface{}) {
	es.Logger.Printf("%v", v...)
}

// Warnf will print warning messages to the console
func (es *ExtendedLogger) Warnf(format string, v ...interface{}) {
	es.Logger.Printf(format, v...)
}
