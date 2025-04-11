package logger

import (
	"fmt"
	"sync"
	"time"

	"github.com/unownone/osark-daemon/internal/service/osarkserver"
	"github.com/unownone/osark-daemon/internal/service/osquery"
	"github.com/unownone/osark-daemon/internal/utils"
	"github.com/unownone/osark-daemon/models"
)

// Debt: we currently track all app

// Service is the interface for the logger service
type Service interface {
	Start() error
	Stop() error
	Wait()
}

// A Highlevel service that manages the system logger
// It is responsible for logging events to the system logger
// and pushing them to the server
type loggerService struct {
	oqManager        osquery.Manager
	serverManager    osarkserver.Manager
	eventChan        chan *models.LogEvent
	waitGroup        *sync.WaitGroup
	delay            time.Duration
	batchSize        int
	stopChan         chan *struct{}
	trackedBundleIDs []string
}

// NewLoggerService creates a new logger service
func NewLoggerService(oqManager osquery.Manager, serverManager osarkserver.Manager, batchSize int) Service {
	return &loggerService{
		oqManager:     oqManager,
		serverManager: serverManager,
		delay:         1 * time.Second,
		waitGroup:     &sync.WaitGroup{},
		eventChan:     make(chan *models.LogEvent),
		batchSize:     batchSize,
	}
}

// Start starts the logger service
func (s *loggerService) Start() error {
	go s.pusher()       // push events to the server
	go s.recordWorker() // record events
	// Send the init event

	return nil
}

// Wait waits for the logger service to finish
func (s *loggerService) Wait() {
	s.waitGroup.Wait()
}

// Stop stops the logger service
func (s *loggerService) Stop() error {
	s.stopChan <- &struct{}{} // Send a signal to the recordWorker to stop
	close(s.stopChan)
	close(s.eventChan)
	s.waitGroup.Wait()
	fmt.Println("Logger service stopped")
	return nil
}

// pusher pushes events to the server
func (s *loggerService) pusher() error {
	s.waitGroup.Add(1)
	ticker := time.NewTicker(s.delay) // we wait for the delay to push the events that got collected
	batch := utils.NewBatchStore[*models.LogEvent](s.batchSize)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if data, err := batch.GetAndReset(); err != nil {
				s.serverManager.PushError(err)
			} else if len(data) > 0 {
				s.serverManager.Push(data)
			}
		case event, ok := <-s.eventChan:
			if !ok {
				// Channel is closed, flush remaining events and exit
				if data, err := batch.GetAndReset(); err != nil {
					s.serverManager.PushError(err)
				} else {
					s.serverManager.Push(data)
				}
				s.waitGroup.Done()
				return nil
			}
			if data, err := batch.Push(event); err != nil {
				s.serverManager.PushError(err)
			} else if data != nil {
				s.serverManager.Push(*data)
			}
		}
	}
}

// recorder records events
func (s *loggerService) recorder() error {
	var err error
	defer func() {
		if err != nil {
			s.serverManager.PushError(err) // push error to the server
		}
	}()
	processes, err := s.oqManager.GetCurrentRunningProcesses(s.trackedBundleIDs)
	if err != nil {
		return err
	}
	s.eventChan <- &models.LogEvent{
		Intent:    models.IntentRunningProcesses,
		Processes: processes,
	}
	return nil
}

// recordWorker records events
func (s *loggerService) recordWorker() error {
	ticker := time.NewTicker(s.delay)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.recorder()
		case <-s.stopChan:
			return nil
		}
	}
}

// sendInitEvent sends the init event
func (s *loggerService) sendInitEvent() error {
	apps, err := s.oqManager.GetApps()
	if err != nil {
		return err
	}
	sysInfo, err := s.oqManager.GetSystemInfo()
	if err != nil {
		return err
	}
	s.trackedBundleIDs = make([]string, 0, len(apps))
	// TODO: we should track targetted apps
	for _, app := range apps[:10] {
		s.trackedBundleIDs = append(s.trackedBundleIDs, app.BundleID)
	}
	s.eventChan <- &models.LogEvent{
		Intent:     models.IntentInit,
		AppInfo:    apps,
		SystemInfo: sysInfo,
	}
	return nil
}
