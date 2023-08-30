package main

import (
	"context"
	"github.com/compliance-framework/assessment-runtime/config"
	"github.com/compliance-framework/assessment-runtime/plugins"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		cancel()
	}()

	log.SetFormatter(&log.JSONFormatter{})
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

	pluginDownloader := plugins.NewPackageDownloader(confManager.Config().PluginRegistryURL)
	err = pluginDownloader.DownloadPackages(packages)
	if err != nil {
		log.Errorf("Error downloading some of the plugins: %s", err)
		// TODO: If the download error keeps occurring, we should report it back to the control plane.
	}

	scheduler := plugins.NewScheduler(confManager.Assessments())

	wg.Add(1)
	go func() {
		defer wg.Done()
		scheduler.Start(ctx)
	}()

	<-ctx.Done()

	// Wait for all components to finish their cleanup
	wg.Wait()
}
