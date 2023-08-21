package main

import (
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/compliance-framework/assessment-runtime/config"
	"github.com/compliance-framework/assessment-runtime/plugins"
)

const configFilePath = "config.yaml"

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
		log.Error("Error downloading some of the plugins:", err)
	}
}
