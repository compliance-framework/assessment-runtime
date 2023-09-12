package main

import (
	"context"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/config"
	"github.com/compliance-framework/assessment-runtime/internal/registry"
	"github.com/compliance-framework/assessment-runtime/internal/scheduling"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup

	confManager, err := config.NewConfigurationManager()
	if err != nil {
		log.Fatalf("Failed to create configuration manager: %s", err)
	}

	err = events.Connect(confManager.Config().EventBusURL)
	if err != nil {
		log.Fatalf("Failed to connect to event bus: %s", err)
	}

	confManager.Listen()

	// Download plugin packages
	pluginDownloader := registry.NewPackageDownloader(confManager.Config().PluginRegistryURL)
	err = pluginDownloader.DownloadPackages(confManager.Packages())
	if err != nil {
		log.Errorf("Error downloading some of the plugins: %s", err)
	}

	scheduler := scheduling.NewScheduler(confManager.JobConfigs())

	wg.Add(1)
	go func() {
		defer wg.Done()
		scheduler.Start(ctx)
	}()

	<-ctx.Done()

	scheduler.Stop()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Components have all shut down.")
	case <-time.After(5 * time.Second):
		fmt.Println("Timed out waiting for components to shut down; exiting anyway.")
	}

	os.Exit(0)
}
