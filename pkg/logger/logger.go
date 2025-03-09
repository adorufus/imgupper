package logger

import (
	"os"
	"strings"

	"github.com/adorufus/imgupper/config"
	"github.com/rs/zerolog"
)

// Logger is the interface for logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

// ZerologLogger implements Logger using zerolog
type ZerologLogger struct {
	logger zerolog.Logger
}

// New creates a new logger
func New(cfg config.LoggerConfig) (Logger, error) {
	// Set up logger output
	var output *os.File = os.Stdout
	if cfg.File != "" {
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		output = file
	}

	// Set up log level
	level, err := parseLogLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(level)

	// Create logger
	logger := zerolog.New(output).With().Timestamp().Caller().Logger()

	return &ZerologLogger{logger: logger}, nil
}

// Debug logs debug messages
func (l *ZerologLogger) Debug(msg string, args ...interface{}) {
	l.logEvent(l.logger.Debug(), msg, args...)
}

// Info logs info messages
func (l *ZerologLogger) Info(msg string, args ...interface{}) {
	l.logEvent(l.logger.Info(), msg, args...)
}

// Warn logs warning messages
func (l *ZerologLogger) Warn(msg string, args ...interface{}) {
	l.logEvent(l.logger.Warn(), msg, args...)
}

// Error logs error messages
func (l *ZerologLogger) Error(msg string, args ...interface{}) {
	l.logEvent(l.logger.Error(), msg, args...)
}

// Fatal logs fatal messages and exits
func (l *ZerologLogger) Fatal(msg string, args ...interface{}) {
	l.logEvent(l.logger.Fatal(), msg, args...)
}

// logEvent adds fields to the log event
func (l *ZerologLogger) logEvent(event *zerolog.Event, msg string, args ...interface{}) {
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key, ok := args[i].(string)
			if !ok {
				event.Interface("invalid_key", args[i])
				continue
			}
			event.Interface(key, args[i+1])
		}
	}
	event.Msg(msg)
}

func parseLogLevel(level string) (zerolog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel, nil
	case "info":
		return zerolog.InfoLevel, nil
	case "warn":
		return zerolog.WarnLevel, nil
	case "error":
		return zerolog.ErrorLevel, nil
	case "fatal":
		return zerolog.FatalLevel, nil
	case "panic":
		return zerolog.PanicLevel, nil
	default:
		return zerolog.InfoLevel, nil
	}
}
