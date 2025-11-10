package logger

import (
	"log"
	"os"
	"time"
)

var (
	infoLog  *log.Logger
	errorLog *log.Logger
	debugLog *log.Logger
)

func init() {
	infoLog = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLog = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime)
}

// Info logs informational messages
func Info(format string, v ...interface{}) {
	infoLog.Printf(format, v...)
}

// Error logs error messages with file and line number
func Error(format string, v ...interface{}) {
	errorLog.Printf(format, v...)
}

// Debug logs debug messages
func Debug(format string, v ...interface{}) {
	debugLog.Printf(format, v...)
}

// Fatal logs error and exits
func Fatal(format string, v ...interface{}) {
	errorLog.Fatalf(format, v...)
}

// LogDuration logs function duration
func LogDuration(start time.Time, name string) {
	duration := time.Since(start)
	infoLog.Printf("%s took %v", name, duration)
}
