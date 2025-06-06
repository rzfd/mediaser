package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config holds logger configuration
type Config struct {
	Level      string `json:"level" yaml:"level"`
	Output     string `json:"output" yaml:"output"`
	TimeFormat string `json:"time_format" yaml:"time_format"`
	LogFile    string `json:"log_file" yaml:"log_file"`
	ServiceName string `json:"service_name" yaml:"service_name"`
}

// Logger represents a structured logger
type Logger struct {
	*zerolog.Logger
}

// New creates a new logger instance
func New(config Config) *Logger {
	// Set log level
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure time format
	if config.TimeFormat == "" {
		config.TimeFormat = time.RFC3339
	}
	zerolog.TimeFieldFormat = config.TimeFormat

	// Configure output
	var output io.Writer
	switch config.Output {
	case "file":
		if config.LogFile != "" {
			// Create log directory if not exists
			if err := os.MkdirAll(filepath.Dir(config.LogFile), 0755); err != nil {
				output = os.Stdout
			} else {
				file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					output = os.Stdout
				} else {
					output = file
				}
			}
		} else {
			output = os.Stdout
		}
	case "both":
		if config.LogFile != "" {
			// Create log directory if not exists
			if err := os.MkdirAll(filepath.Dir(config.LogFile), 0755); err != nil {
				output = os.Stdout
			} else {
				file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					output = os.Stdout
				} else {
					output = io.MultiWriter(os.Stdout, file)
				}
			}
		} else {
			output = os.Stdout
		}
	default: // stdout
		output = os.Stdout
	}

	// Create console writer for better formatting in development
	if config.Output == "stdout" || config.Output == "" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
			NoColor:    false,
		}
	}

	// Create logger
	logger := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Str("service", config.ServiceName).
		Logger()

	return &Logger{
		Logger: &logger,
	}
}

// DefaultConfig returns default logger configuration
func DefaultConfig(serviceName string) Config {
	return Config{
		Level:       "info",
		Output:      "stdout",
		TimeFormat:  time.RFC3339,
		ServiceName: serviceName,
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	event := l.Logger.Info()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			event = event.Interface(fields[i].(string), fields[i+1])
		}
	}
	event.Msg(msg)
}

// Error logs an error message
func (l *Logger) Error(err error, msg string, fields ...interface{}) {
	event := l.Logger.Error().Err(err)
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			event = event.Interface(fields[i].(string), fields[i+1])
		}
	}
	event.Msg(msg)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	event := l.Logger.Debug()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			event = event.Interface(fields[i].(string), fields[i+1])
		}
	}
	event.Msg(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	event := l.Logger.Warn()
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			event = event.Interface(fields[i].(string), fields[i+1])
		}
	}
	event.Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(err error, msg string, fields ...interface{}) {
	event := l.Logger.Fatal().Err(err)
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			event = event.Interface(fields[i].(string), fields[i+1])
		}
	}
	event.Msg(msg)
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := l.Logger.With().Interface(key, value).Logger()
	return &Logger{Logger: &newLogger}
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	event := l.Logger.With()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	newLogger := event.Logger()
	return &Logger{Logger: &newLogger}
}

// Global logger instance
var GlobalLogger *Logger

// Init initializes the global logger
func Init(config Config) {
	GlobalLogger = New(config)
	log.Logger = *GlobalLogger.Logger
}

// GetLogger returns the global logger
func GetLogger() *Logger {
	if GlobalLogger == nil {
		GlobalLogger = New(DefaultConfig("default"))
	}
	return GlobalLogger
} 