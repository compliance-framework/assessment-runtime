package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"

	"github.com/compliance-framework/assessment-runtime/config"
	"github.com/compliance-framework/assessment-runtime/plugins"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	confManager := config.NewConfigurationManager()

	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execDir := filepath.Dir(execPath)
	configFilePath := filepath.Join(execDir, "config.yml")

	cfg, err := confManager.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load cfg: %s", err)
	}

	log.Infof("Configuration loaded successfully: %v", cfg)

	pluginDownloader := plugins.NewPluginDownloader(cfg)
	err = pluginDownloader.DownloadPlugins()
	if err != nil {
		log.Errorf("Error downloading some of the plugins: %s", err)
		// TODO: If the download error keeps occurring, we should report it back to the control plane.
	}

	pluginManager := plugins.NewPluginManager(cfg)

	err = pluginManager.Start()
	if err != nil {
		log.Errorf("Error starting plugins: %s", err)
	}

	scheduler := plugins.NewScheduler()
	for _, plugin := range cfg.Plugins {
		scheduler.AddJob(plugin.Schedule, func() {
			err := pluginManager.Execute(plugin.Name, plugins.ActionInput{})
			if err != nil {
				log.Errorf("Error starting plugin %s: %s", plugin.Name, err)
				return
			}
		})
	}
	scheduler.Start()

	select {}
}
