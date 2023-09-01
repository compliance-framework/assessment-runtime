package main

import (
	"context"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/config"
	plugins2 "github.com/compliance-framework/assessment-runtime/internal/plugins"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup

	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	confManager, err := config.NewConfigurationManager()
	if err != nil {
		log.Fatalf("Failed to create configuration manager: %s", err)
	}

	// Download plugin packages
	packages, err := confManager.Packages()
	if err != nil {
		log.Fatalf("Failed to get packages: %s", err)
	}

	pluginDownloader := plugins2.NewPackageDownloader(confManager.Config().PluginRegistryURL)
	err = pluginDownloader.DownloadPackages(packages)
	if err != nil {
		log.Errorf("Error downloading some of the plugins: %s", err)
		// TODO: If the download error keeps occurring, we should report it back to the control plane.
	}

	scheduler := plugins2.NewScheduler(confManager.Assessments())

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
