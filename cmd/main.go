package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/unownone/osark-daemon/internal/service/logger"
	"github.com/unownone/osark-daemon/internal/service/osarkserver"
	"github.com/unownone/osark-daemon/internal/service/osquery"
)

var (
	OsarkServerURL string = "http://127.0.0.1:3000"
	LogDir         string = "logs"
)

// multiWriter is a simple io.Writer that writes to multiple io.Writers
type multiWriter struct {
	writers []io.Writer
}

func (mw *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range mw.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
	}
	return len(p), nil
}

// setupLogging initializes the logging system to write to both console and file
func setupLogging() (string, error) {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(LogDir, 0755); err != nil {
		return "", err
	}

	// Create log file with timestamp in filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(LogDir, "osark_"+timestamp+".log")
	logFile, err := os.Create(logFilePath)
	if err != nil {
		return "", err
	}

	// Create a writer that writes to both stdout and the log file
	mw := &multiWriter{
		writers: []io.Writer{os.Stdout, logFile},
	}

	// Initialize slog with text handler writing to both stdout and file
	slogHandler := slog.NewTextHandler(mw, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})
	slogLogger := slog.New(slogHandler)
	slog.SetDefault(slogLogger)

	return logFilePath, nil
}

// setupSignalHandling sets up a handler for interrupt signals
func setupSignalHandling(cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		sig := <-signalChan
		slog.Info("Received shutdown signal", "signal", sig)
		cancel()
	}()
}

// initializeServices initializes and sets up all required services
func initializeServices() (osquery.Manager, osarkserver.Manager, logger.Service, error) {
	if OsarkServerURL == "" {
		return nil, nil, nil, errorf("OSARK_SERVER_URL is not set")
	}
	
	manager, err := osquery.NewManager()
	if err != nil {
		return nil, nil, nil, errorf("failed to create manager: %v", err)
	}
	
	sysInfo, err := manager.GetSystemInfo()
	if err != nil {
		return nil, nil, nil, errorf("failed to get system info: %v", err)
	}
	
	serverManager, err := osarkserver.NewPushManager(OsarkServerURL, sysInfo)
	if err != nil {
		return nil, nil, nil, errorf("failed to create push manager: %v", err)
	}
	
	loggerService := logger.NewLoggerService(manager, serverManager, 100)
	return manager, serverManager, loggerService, nil
}

// performGracefulShutdown gracefully shuts down the service with a timeout
func performGracefulShutdown(loggerService logger.Service, timeout time.Duration) {
	slog.Info("Initiating graceful shutdown")
	
	// Create a timeout context for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()
	
	// Attempt graceful shutdown with timeout
	shutdownComplete := make(chan struct{})
	go func() {
		err := loggerService.Stop()
		if err != nil {
			slog.Error("Error during service shutdown", "error", err)
		}
		close(shutdownComplete)
	}()
	
	// Wait for shutdown to complete or timeout
	select {
	case <-shutdownComplete:
		slog.Info("Graceful shutdown completed")
	case <-shutdownCtx.Done():
		slog.Warn("Graceful shutdown timed out, forcing exit")
	}
	
	slog.Info("Application shutdown complete")
}

// errorf creates a new error with the given format and arguments
func errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func main() {
	// Setup logging
	logFilePath, err := setupLogging()
	if err != nil {
		panic("Failed to setup logging: " + err.Error())
	}
	slog.Info("Logging initialized", "logFile", logFilePath)

	// Create a context that will be canceled on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	setupSignalHandling(cancel)

	// Initialize services
	_, _, loggerService, err := initializeServices()
	if err != nil {
		slog.Error("Service initialization failed", "error", err)
		os.Exit(1)
	}

	// Start the logger service
	loggerService.Start()
	slog.Info("Logger service started")
	
	// Wait for cancel signal from context
	<-ctx.Done()
	
	// Perform graceful shutdown
	performGracefulShutdown(loggerService, 10*time.Second)
}
