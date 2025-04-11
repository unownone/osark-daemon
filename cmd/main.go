package main

import (
	"fmt"
	"log"

	"github.com/unownone/osark-daemon/internal/service/logger"
	"github.com/unownone/osark-daemon/internal/service/osarkserver"
	"github.com/unownone/osark-daemon/internal/service/osquery"
)

var (
	OsarkServerURL string = "http://127.0.0.1:3000"
)

func main() {
	if OsarkServerURL == "" {
		log.Fatalf("OSARK_SERVER_URL is not set")
	}
	manager, err := osquery.NewManager()
	if err != nil {
		log.Fatalf("failed to create manager: %v", err)
	}
	sysInfo, err := manager.GetSystemInfo()
	if err != nil {
		log.Fatalf("failed to get system info: %v", err)
	}
	serverManager, err := osarkserver.NewPushManager(OsarkServerURL, sysInfo)
	if err != nil {
		log.Fatalf("failed to create push manager: %v", err)
	}
	loggerService := logger.NewLoggerService(manager, serverManager, 100)
	loggerService.Start()
	fmt.Println("Logger service started")
	loggerService.Wait()
	fmt.Println("Logger service stopped")
}
