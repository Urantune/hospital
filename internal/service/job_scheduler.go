package service

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type JobScheduler struct {
	ticker           *time.Ticker
	stopChan         chan struct{}
	wg               sync.WaitGroup
	reminderInterval time.Duration
	retryInterval    time.Duration
	isRunning        bool
	mu               sync.Mutex
}

func NewJobScheduler(reminderInterval, retryInterval time.Duration) *JobScheduler {
	if reminderInterval == 0 {
		reminderInterval = 1 * time.Minute
	}
	if retryInterval == 0 {
		retryInterval = 5 * time.Minute
	}

	return &JobScheduler{
		reminderInterval: reminderInterval,
		retryInterval:    retryInterval,
	}
}

func (js *JobScheduler) Start() error {
	js.mu.Lock()
	defer js.mu.Unlock()

	if js.isRunning {
		return fmt.Errorf("job scheduler is already running")
	}

	js.stopChan = make(chan struct{})
	js.ticker = time.NewTicker(js.reminderInterval)
	js.isRunning = true

	js.wg.Add(1)
	go js.run()

	js.wg.Add(1)
	go js.runRetryJob()

	log.Printf("[JobScheduler] started: reminder interval=%s, retry interval=%s", js.reminderInterval, js.retryInterval)
	return nil
}

func (js *JobScheduler) Stop() error {
	js.mu.Lock()
	if !js.isRunning {
		js.mu.Unlock()
		return fmt.Errorf("job scheduler is not running")
	}

	js.isRunning = false
	if js.ticker != nil {
		js.ticker.Stop()
	}
	close(js.stopChan)
	js.mu.Unlock()

	done := make(chan struct{})
	go func() {
		js.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("[JobScheduler] stopped gracefully")
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("job scheduler stop timeout")
	}
}

func (js *JobScheduler) run() {
	defer js.wg.Done()

	for {
		select {
		case <-js.stopChan:
			log.Println("[JobScheduler] reminder job stopped")
			return
		case <-js.ticker.C:
			js.processRemindersJob()
		}
	}
}

func (js *JobScheduler) runRetryJob() {
	defer js.wg.Done()

	retryTicker := time.NewTicker(js.retryInterval)
	defer retryTicker.Stop()

	for {
		select {
		case <-js.stopChan:
			log.Println("[JobScheduler] retry job stopped")
			return
		case <-retryTicker.C:
			js.processRetryJob()
		}
	}
}

func (js *JobScheduler) processRemindersJob() {
	startTime := time.Now()
	sentCount, failedCount, err := ProcessPendingReminders()
	if err != nil {
		log.Printf("[JobScheduler] pending reminder processing failed: %v", err)
		return
	}

	log.Printf("[JobScheduler] reminder processing complete: sent=%d failed=%d duration=%v", sentCount, failedCount, time.Since(startTime))
}

func (js *JobScheduler) processRetryJob() {
	startTime := time.Now()
	retriedCount, failedCount, err := RetryFailedReminders()
	if err != nil {
		log.Printf("[JobScheduler] failed reminder retry failed: %v", err)
		return
	}

	log.Printf("[JobScheduler] retry processing complete: retried=%d failed=%d duration=%v", retriedCount, failedCount, time.Since(startTime))
}

func (js *JobScheduler) IsRunning() bool {
	js.mu.Lock()
	defer js.mu.Unlock()
	return js.isRunning
}

var globalScheduler *JobScheduler

func InitializeJobScheduler(reminderInterval, retryInterval time.Duration) (*JobScheduler, error) {
	globalScheduler = NewJobScheduler(reminderInterval, retryInterval)
	if err := globalScheduler.Start(); err != nil {
		return nil, err
	}
	return globalScheduler, nil
}

func GetGlobalScheduler() *JobScheduler {
	return globalScheduler
}

func ShutdownJobScheduler() error {
	if globalScheduler != nil && globalScheduler.IsRunning() {
		return globalScheduler.Stop()
	}
	return nil
}
