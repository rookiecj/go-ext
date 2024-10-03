// Description: A simple logger package that wraps slog with colored output.
package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// LogLevel is a custom type for log levels
type LogLevel int

const (
	// VerboseLevel is the verbose log level
	// This level of message is used for in detail logging like function entry and exit
	VerboseLevel LogLevel = iota
	// DebugLevel is the debug log level(default log level)
	// less or equal to this level of message is used in development environment
	DebugLevel
	// InfoLevel is the info log level
	// greater or equal to this level of message shows up in production environment
	InfoLevel
	// WarnLevel is the warn log level
	// This can be used to warn the user about something, but not to report an error
	WarnLevel
	// ErrorLevel is the error log level
	// This can be used to report errors
	ErrorLevel
	// PanicLevel is the panic log level
	// This raise a panic after logging the message
	PanicLevel
)

func (l LogLevel) String() string {
	switch l {
	case VerboseLevel:
		return "VRB"
	case DebugLevel:
		return "DBG"
	case InfoLevel:
		return "INF"
	case WarnLevel:
		return "WRN"
	case ErrorLevel:
		return "ERR"
	case PanicLevel:
		return "PNC"
	default:
		return "DBG"
	}
}

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger(os.Stderr)
}

// ANSI color codes
const (
	reset = "\033[0m"
	red   = "\033[31m" // error
	pink  = "\033[35m" // warn
	green = "\033[32m" // info
	//yellow = "\033[33m" 	// info
	white = "\033[37m" // debug
)

const DefaultTimeFormatPattern = "2006-01-02 15:04:05"

// LogFormatter is an interface for log formatters
type LogFormatter interface {
	Format(logTime time.Time, level LogLevel, msg string) string
}

type DefaultFormatter struct{}

// NewDefaultFormatter creates a new default log formatter
func NewDefaultFormatter() LogFormatter {
	return &DefaultFormatter{}
}

func (f *DefaultFormatter) Format(logTime time.Time, level LogLevel, msg string) string {
	return fmt.Sprintf("%s %s %s", logTime.Format(DefaultTimeFormatPattern), level.String(), msg)
}

type ColoredLevelFormatter struct{}

// NewColoredLevelFormatter creates a new loevel colored log formatter
func NewColoredLevelFormatter() LogFormatter {
	return &ColoredLevelFormatter{}
}

func (f *ColoredLevelFormatter) Format(logTime time.Time, level LogLevel, msg string) string {
	var color string
	switch level {
	case VerboseLevel:
		color = white
	case DebugLevel:
		color = white
	case InfoLevel:
		color = green
	case WarnLevel:
		color = pink
	case ErrorLevel:
		color = red
	case PanicLevel:
		color = red
	default:
		color = white
	}

	//timeStr := logTime.Format(time.RFC3339)
	timeStr := logTime.Format(DefaultTimeFormatPattern)
	//coloredLevelStr := fmt.Sprintf("%s%s%s", color, level.String(), reset)
	//logMsg := fmt.Sprintf("%s %s %s", timeStr, coloredLevelStr, msg)
	logMsg := fmt.Sprintf("%s %s%s%s %s", timeStr, color, level.String(), reset, msg)
	return logMsg
}

type Logger struct {
	ctx       context.Context
	mu        sync.Mutex
	wr        io.Writer
	logLevel  LogLevel
	formatter LogFormatter
}

// NewLogger creates a new logger with the given writer
func NewLogger(wr io.Writer) *Logger {
	logger := &Logger{
		ctx: context.Background(),
		mu:  sync.Mutex{},
		wr:  wr,
		logLevel: func() LogLevel {
			logLevel := DebugLevel
			if level := os.Getenv("LOG_LEVEL"); level != "" {
				switch level {
				case "verbose":
					logLevel = VerboseLevel
				case "debug":
					logLevel = DebugLevel
				case "info":
					logLevel = InfoLevel
				case "warn":
					logLevel = WarnLevel
				case "error":
					logLevel = ErrorLevel
				case "panic":
					logLevel = PanicLevel
				default:
					logLevel = DebugLevel
				}
			}
			return logLevel
		}(),
	}
	logger.formatter = &DefaultFormatter{}
	return logger
}

// NewLoggerWithFormatter creates a new logger with the given writer and formatter
func NewLoggerWithFormatter(wr io.Writer, formatter LogFormatter) *Logger {
	logger := NewLogger(wr)
	logger.SetFormatter(formatter)
	return logger
}

// GetLogger returns the default logger
func GetLogger() *Logger {
	return defaultLogger
}

func (c *Logger) SetLogLevel(level LogLevel) {
	c.logLevel = level
}

func (c *Logger) SetFormatter(formatter LogFormatter) {
	c.formatter = formatter
}

func (c *Logger) SetWriter(wr io.Writer) {
	Debugf("Setting writer to %v\n", wr)
	c.mu.Lock()
	c.wr = wr
	c.mu.Unlock()
}

func (c *Logger) Write(p []byte) (n int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.wr != nil {
		return c.wr.Write(p)
	}
	return 0, nil
}

// Log logs a message with the given log level
func Log(ctx context.Context, level LogLevel, msg string, args ...any) {
	if defaultLogger.logLevel > level {
		return
	}
	logMsg := defaultLogger.formatter.Format(time.Now(), level, fmt.Sprintf(msg, args...))
	_, err := defaultLogger.Write([]byte(logMsg))
	if err != nil {
		fmt.Printf("Error writing log message: %v\n", err)
	}
}

// SetLogLevel sets the log level of the default logger
func SetLogLevel(level LogLevel) {
	defaultLogger.SetLogLevel(level)
}

// SetFormatter sets the formatter of the default logger
func SetFormatter(formatter LogFormatter) {
	defaultLogger.SetFormatter(formatter)
}

// SetWriter sets the writer of the default logger
func SetWriter(wr io.Writer) {
	defaultLogger.SetWriter(wr)
}

// Logf logs a message with the given log level
func Logf(level LogLevel, msg string, args ...any) {
	Log(defaultLogger.ctx, level, msg, args...)
}

// Verbose logs a message with the verbose log level
func Verbose(ctx context.Context, fmt string, args ...interface{}) {
	Log(ctx, VerboseLevel, fmt, args...)
}

// Verbose logs a message with the verbose log level
func Verbosef(fmt string, args ...interface{}) {
	Logf(VerboseLevel, fmt, args...)
}

// Debug logs a message with the debug log level
func Debug(ctx context.Context, fmt string, args ...interface{}) {
	Log(ctx, DebugLevel, fmt, args...)
}

// Debug logs a message with the debug log level
func Debugf(fmt string, args ...interface{}) {
	Logf(DebugLevel, fmt, args...)
}

// Info logs a message with the info log level
func Info(ctx context.Context, fmt string, args ...interface{}) {
	Log(ctx, InfoLevel, fmt, args...)
}

// Info logs a message with the info log level
func Infof(fmt string, args ...interface{}) {
	Logf(InfoLevel, fmt, args...)
}

// Warn logs a message with the warn log level
func Warn(ctx context.Context, fmt string, args ...interface{}) {
	Log(ctx, WarnLevel, fmt, args...)
}

// Warn logs a message with the warn log level
func Warnf(fmt string, args ...interface{}) {
	Logf(WarnLevel, fmt, args...)
}

// Error logs a message with the error log level
func Error(ctx context.Context, fmt string, args ...interface{}) {
	Log(ctx, ErrorLevel, fmt, args...)
}

// Error logs a message with the error log level
func Errorf(fmt string, args ...interface{}) {
	Logf(ErrorLevel, fmt, args...)
}

// Panic logs a message with the panic log level and raises a panic
func Panic(ctx context.Context, fmtStr string, args ...interface{}) {
	Log(ctx, PanicLevel, fmtStr, args...)
	panic(fmt.Sprintf(fmtStr, args...))
}

// Panic logs a message with the panic log level and raises a panic
func Panicf(fmtStr string, args ...interface{}) {
	Logf(PanicLevel, fmtStr, args...)
	panic(fmt.Sprintf(fmtStr, args...))
}
