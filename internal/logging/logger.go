package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/logger"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger      *log.Logger
	WailsLogger *WailsFileLogger
)

// WailsFileLogger implements the Wails logger.Logger interface with file rotation
type WailsFileLogger struct {
	logger   *log.Logger
	logFile  *lumberjack.Logger
	logLevel logger.LogLevel
}

// InitLogger initializes the logger with file rotation and Wails compatibility
// The log files will be saved in es_gaze_logs folder within the provided base directory
func InitLogger(baseDir string) (*WailsFileLogger, error) {
	// Create es_gaze_logs directory
	logDir := filepath.Join(baseDir, "es_gaze_logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// Configure lumberjack for log rotation
	// MaxSize is in megabytes, so 150MB
	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "elasticgaze.log"),
		MaxSize:    150,  // megabytes
		MaxBackups: 5,    // keep 5 old log files
		MaxAge:     30,   // keep logs for 30 days
		Compress:   true, // compress old log files
	}

	// Create multi-writer to write to both file and console
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Initialize the logger
	stdLogger := log.New(multiWriter, "", log.LstdFlags|log.Lshortfile)

	// Create Wails-compatible logger
	wailsLogger := &WailsFileLogger{
		logger:   stdLogger,
		logFile:  logFile,
		logLevel: logger.INFO, // Default log level
	}

	// Set global variables for backward compatibility
	Logger = stdLogger
	WailsLogger = wailsLogger

	wailsLogger.Info("Logger initialized successfully")
	return wailsLogger, nil
}

// Wails Logger interface implementation
// These methods implement the logger.Logger interface required by Wails

// Print logs a raw message
func (w *WailsFileLogger) Print(message string) {
	w.logger.Print("[PRINT] " + message)
}

// Trace logs a trace level message
func (w *WailsFileLogger) Trace(message string) {
	if w.logLevel <= logger.TRACE {
		w.logger.Print("[TRACE] " + message)
	}
}

// Debug logs a debug level message
func (w *WailsFileLogger) Debug(message string) {
	if w.logLevel <= logger.DEBUG {
		w.logger.Print("[DEBUG] " + message)
	}
}

// Info logs an info level message
func (w *WailsFileLogger) Info(message string) {
	if w.logLevel <= logger.INFO {
		w.logger.Print("[INFO] " + message)
	}
}

// Warning logs a warning level message
func (w *WailsFileLogger) Warning(message string) {
	if w.logLevel <= logger.WARNING {
		w.logger.Print("[WARNING] " + message)
	}
}

// Error logs an error level message
func (w *WailsFileLogger) Error(message string) {
	if w.logLevel <= logger.ERROR {
		w.logger.Print("[ERROR] " + message)
	}
}

// Fatal logs a fatal level message
func (w *WailsFileLogger) Fatal(message string) {
	w.logger.Print("[FATAL] " + message)
}

// SetLogLevel sets the minimum log level
func (w *WailsFileLogger) SetLogLevel(level logger.LogLevel) {
	w.logLevel = level
}

// Backward compatibility functions for existing code

// Info logs an info message
func Info(v ...interface{}) {
	if Logger != nil {
		Logger.Println(append([]interface{}{"[INFO]"}, v...)...)
	}
}

// Error logs an error message
func Error(v ...interface{}) {
	if Logger != nil {
		Logger.Println(append([]interface{}{"[ERROR]"}, v...)...)
	}
}

// Debug logs a debug message
func Debug(v ...interface{}) {
	if Logger != nil {
		Logger.Println(append([]interface{}{"[DEBUG]"}, v...)...)
	}
}

// Warn logs a warning message
func Warn(v ...interface{}) {
	if Logger != nil {
		Logger.Println(append([]interface{}{"[WARN]"}, v...)...)
	}
}

// Infof logs an info message with formatting
func Infof(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[INFO] "+format, v...)
	}
}

// Errorf logs an error message with formatting
func Errorf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[ERROR] "+format, v...)
	}
}

// Debugf logs a debug message with formatting
func Debugf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[DEBUG] "+format, v...)
	}
}

// Warnf logs a warning message with formatting
func Warnf(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[WARN] "+format, v...)
	}
}
