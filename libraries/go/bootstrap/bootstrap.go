package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/AyyItsAljo/smart-room/libraries/go/api"
	"github.com/AyyItsAljo/smart-room/libraries/go/config"
	"github.com/AyyItsAljo/smart-room/libraries/go/firehose"
	"github.com/AyyItsAljo/smart-room/libraries/go/slog"

	"github.com/go-redis/redis"
)

// Process is a long-running task that provides service functionality
type Process interface {
	// GetName returns a friendly name for the process for use in logs
	GetName() string

	// Start kicks off the task and only returns when the task has finished
	Start() error

	// Stop will try to gracefully end the task and should be safe to run regardless of whether the process is currently running
	Stop(context.Context) error
}

// Boot performs standard service startup tasks
func Init(serviceName string) error {
	// Create default API client
	apiGateway := os.Getenv("API_GATEWAY")
	if apiGateway == "" {
		return fmt.Errorf("API_GATEWAY env var not set")
	}
	apiClient, err := api.New(apiGateway, "data")
	if err != nil {
		return err
	}
	api.DefaultClient = apiClient

	// Load config
	var configRsp map[string]interface{}
	_, err = api.Get(fmt.Sprintf("service.config/read/%s", serviceName), &configRsp)
	if err != nil {
		return err
	}
	config.DefaultProvider = config.New(configRsp)

	// Connect to Redis
	if config.Has("redis.host") {
		host := config.Get("redis.host").String()
		port := config.Get("redis.port").Int()
		addr := fmt.Sprintf("%s:%d", host, port)
		slog.Info("Connecting to Redis at address %s", addr)
		firehose.DefaultPublisher = firehose.New(redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
		}))
	}

	return nil
}

// Run takes a number of processes and concurrently runs them all. It will stop if all processes
// terminate or if a signal (SIGINT or SIGTERM) is received.
func Run(processes ...Process) {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}

	// Start all of the processes in goroutines
	for _, process := range processes {
		process := process

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := process.Start(); err != nil {
				slog.Error("Process %s stopped with error: %v", process.GetName(), err)
			} else {
				slog.Debug("Process %s stopped", process.GetName())
			}
		}()
	}

	// Close the done channel when all processes return
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for all processes to return or for a signal
	select {
	case <-done:
		slog.Warn("All processes stopped prematurely")
		os.Exit(1)
	case s := <-sig:
		slog.Info("Received %v signal", s)
	}

	// A short timeout because Docker will kill us after 10 seconds anyway
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Simultaneously stop all processes
	for _, process := range processes {
		process := process

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := process.Stop(ctx); err != nil {
				slog.Error("Failed to stop %s gracefully: %v", process.GetName(), err)
			}
		}()
	}

	// Wait for everything to finish or a timeout to be hit
	select {
	case <-time.After(time.Second * 9):
		slog.Error("Failed to stop processes in time")
		os.Exit(1)
	case <-done:
		slog.Info("All processes stopped")
		os.Exit(0)
	}
}
