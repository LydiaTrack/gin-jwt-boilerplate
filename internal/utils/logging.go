package utils

import (
	"io"
	"log"
	"os"
)

var logger *log.Logger

// InitLogging initializes logging by creating a new logger and setting its flags and output
func InitLogging() {
	// Save log to file with date and time
	logfile, err := os.Create("app.log")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	// Save on both stdout and file
	write := io.MultiWriter(os.Stdout, logfile)
	logger = log.New(write, "app ", log.LstdFlags)

	defer logfile.Close()
	Log("Logging initialized")
}

// Log with severity INFO
func Log(message string, args ...interface{}) {
	logger.Printf("[INFO] "+message, args...)
}

// LogError with severity ERROR
func LogError(message string, args ...interface{}) {
	logger.Printf("[ERROR] "+message, args...)
}

// LogFatal with severity FATAL
func LogFatal(message string, args ...interface{}) {
	logger.Fatalf("[FATAL] "+message, args...)
}

// LogPanic with severity PANIC
func LogPanic(message string) {
	logger.Panicln("[PANIC] " + message)
}

// LogWarning with severity WARNING
func LogWarning(message string) {
	logger.Println("[WARNING] " + message)
}

// LogDebug with severity DEBUG
func LogDebug(message string) {
	logger.Println("[DEBUG] " + message)
}
