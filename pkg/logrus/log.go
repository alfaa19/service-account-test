// internal/logger/logger.go
package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// Level represents log levels
type Level string

const (
	// Log levels
	LevelDebug    Level = "DEBUG"
	LevelInfo     Level = "INFO"
	LevelWarning  Level = "WARNING"
	LevelError    Level = "ERROR"
	LevelCritical Level = "CRITICAL"
)

// Config holds logger configuration
type Config struct {
	LogLevel         Level
	LogToConsole     bool
	LogToFile        bool
	LogFilePath      string
	LogFileMaxSize   int
	LogFileMaxBackup int
	LogFileMaxAge    int
}

// CustomLogger wraps logrus logger
type CustomLogger struct {
	*logrus.Logger
}

// NewLogger creates a new logger instance that logs to both file and console
func NewLogger(cfg Config) (*CustomLogger, error) {
	// Create base logger
	logger := logrus.New()

	// Set log level
	level := convertLevel(cfg.LogLevel)
	logger.SetLevel(level)

	// Set formatter
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Create output writers
	var writers []io.Writer

	// Add console writer if enabled
	if cfg.LogToConsole {
		writers = append(writers, os.Stdout)
	}

	// Add file writer if enabled
	if cfg.LogToFile {
		if cfg.LogFilePath == "" {
			// Default to logs directory if not specified
			cfg.LogFilePath = filepath.Join("logs", "application.log")
		}

		// Ensure directory exists
		err := os.MkdirAll(filepath.Dir(cfg.LogFilePath), 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Open log file
		file, err := os.OpenFile(cfg.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		writers = append(writers, file)
	}

	// Create multi-writer if we have multiple outputs
	if len(writers) > 1 {
		multiWriter := io.MultiWriter(writers...)
		logger.SetOutput(multiWriter)
	} else if len(writers) == 1 {
		logger.SetOutput(writers[0])
	}

	return &CustomLogger{
		Logger: logger,
	}, nil
}

// WithRequestID adds a request ID field to the logger entry
func (l *CustomLogger) WithRequestID(requestID string) *logrus.Entry {
	return l.WithField("request_id", requestID)
}

// WithContext extracts request ID from context and adds it to log entry
func (l *CustomLogger) WithContext(ctx context.Context) *logrus.Entry {
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		return l.WithRequestID(requestID)
	}
	return l.WithFields(logrus.Fields{})
}

// LogOperation logs an operation with consistent format and context
func (l *CustomLogger) LogOperation(ctx context.Context, operation string, result string, data map[string]interface{}) {
	fields := logrus.Fields{
		"operation": operation,
		"result":    result,
	}

	// Add all data to fields
	for k, v := range data {
		fields[k] = v
	}

	entry := l.WithContext(ctx).WithFields(fields)

	switch result {
	case "success":
		entry.Info("Operation completed successfully")
	case "warning":
		entry.Warn("Operation completed with warnings")
	case "error":
		entry.Error("Operation failed")
	case "critical":
		entry.Fatal("Critical system failure")
	default:
		entry.Info("Operation status")
	}
}

// LogDebug logs a debug message
func (l *CustomLogger) LogDebug(ctx context.Context, msg string, data map[string]interface{}) {
	fields := logrus.Fields{}
	for k, v := range data {
		fields[k] = v
	}
	l.WithContext(ctx).WithFields(fields).Debug(msg)
}

// LogInfo logs an info message
func (l *CustomLogger) LogInfo(ctx context.Context, msg string, data map[string]interface{}) {
	fields := logrus.Fields{}
	for k, v := range data {
		fields[k] = v
	}
	l.WithContext(ctx).WithFields(fields).Info(msg)
}

// LogWarning logs a warning message
func (l *CustomLogger) LogWarning(ctx context.Context, msg string, data map[string]interface{}) {
	fields := logrus.Fields{}
	for k, v := range data {
		fields[k] = v
	}
	l.WithContext(ctx).WithFields(fields).Warn(msg)
}

// LogError logs an error message
func (l *CustomLogger) LogError(ctx context.Context, msg string, data map[string]interface{}) {
	fields := logrus.Fields{}
	for k, v := range data {
		fields[k] = v
	}
	l.WithContext(ctx).WithFields(fields).Error(msg)
}

// LogCritical logs a critical message
func (l *CustomLogger) LogCritical(ctx context.Context, msg string, data map[string]interface{}) {
	fields := logrus.Fields{}
	for k, v := range data {
		fields[k] = v
	}
	l.WithContext(ctx).WithFields(fields).Fatal(msg)
}

// // ContextWithRequestID adds a request ID to the context
// func ContextWithRequestID(ctx context.Context, requestId string) context.Context {
// 	return context.WithValue(ctx, "request_id", requestId)
// }

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	return ""
}

// Helper function to convert our level to logrus level
func convertLevel(level Level) logrus.Level {
	switch level {
	case LevelDebug:
		return logrus.DebugLevel
	case LevelInfo:
		return logrus.InfoLevel
	case LevelWarning:
		return logrus.WarnLevel
	case LevelError:
		return logrus.ErrorLevel
	case LevelCritical:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
