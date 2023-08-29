package main

import (
	"github.com/compliance-framework/assessment-runtime/config"
	"github.com/compliance-framework/assessment-runtime/plugins"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
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
	scheduler.Start()

	select {}
}
