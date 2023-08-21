package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/compliance-framework/assessment-runtime/config"
	"github.com/compliance-framework/assessment-runtime/plugins"
)

const configFilePath = "./config.yaml"

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	confManager := config.NewConfigurationManager()

	cfg, err := confManager.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load cfg: %s", err)
	}

	log.Infof("Cfg loaded successfully: %v", cfg)

	pluginDownloader := plugins.NewPluginDownloader(cfg)
	err = pluginDownloader.DownloadPlugins()
	if err != nil {
		log.Fatalf("Error downloading some of the plugins: %s", err)
	}

	pluginManager := plugins.NewPluginManager(cfg)
	err = pluginManager.InitPlugins()
	if err != nil {
		log.Fatalf("Error initializing plugins: %s", err)
	}

	scheduler := plugins.NewScheduler()
	for _, plugin := range cfg.Plugins {
		scheduler.AddJob(plugin.Schedule, func() {
			pluginManager.StartPlugin(plugin.Name)
		})
	}
	scheduler.Start()

	select {}
}
