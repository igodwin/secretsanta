package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/igodwin/secretsanta/internal/api"
	"github.com/igodwin/secretsanta/pkg/config"
)

// Build-time variables (set via -ldflags)
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

// Logger wraps the standard logger with level support
type Logger struct {
	*log.Logger
}

func NewLogger(out io.Writer) *Logger {
	return &Logger{
		Logger: log.New(out, "", 0),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.logWithLevel("INFO", format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.logWithLevel("ERROR", format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.logWithLevel("FATAL", format, v...)
	os.Exit(1)
}

func (l *Logger) logWithLevel(level, format string, v ...interface{}) {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	message := fmt.Sprintf(format, v...)
	l.Logger.Printf("%s [%s] %s", timestamp, level, message)
}

// Print without log level (for banners, etc.)
func (l *Logger) Print(msg string) {
	fmt.Println(msg)
}

func printBanner(logger *Logger) {
	banner := strings.Builder{}
	banner.WriteString("\n====================================\n")
	banner.WriteString("Secret Santa Web Server\n")
	banner.WriteString(fmt.Sprintf("Version:    %s\n", Version))
	banner.WriteString(fmt.Sprintf("Git Commit: %s\n", GitCommit))
	banner.WriteString(fmt.Sprintf("Build Time: %s\n", BuildTime))
	banner.WriteString("====================================\n")
	logger.Print(banner.String())
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP server address")
	flag.Parse()

	// Create custom logger
	logger := NewLogger(os.Stdout)

	// Print startup banner
	printBanner(logger)

	// Configure global log package to use our logger format
	log.SetFlags(0)
	log.SetOutput(&logWriter{logger: logger})

	// Load configuration at startup
	_ = config.GetConfig()

	logger.Info("Visit http://localhost%s to get started", *addr)

	server := api.NewServer(*addr)

	if err := server.Start(); err != nil {
		logger.Fatal("Server failed to start: %v", err)
	}
}

// logWriter adapts our Logger to io.Writer for the global log package
type logWriter struct {
	logger *Logger
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	// Remove trailing newline from log package
	msg := strings.TrimSuffix(string(p), "\n")
	w.logger.Info("%s", msg)
	return len(p), nil
}
