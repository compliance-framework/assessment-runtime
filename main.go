package main

import (
	"fmt"
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
		fmt.Printf("failed to load cfg: %s", err)
		os.Exit(1)
	}

	fmt.Printf("cfg loaded successfully: %v", cfg)

	pluginManager := plugins.NewPluginManager(cfg)
	err = pluginManager.DownloadPlugins()
	if err != nil {
		fmt.Println("Error downloading plugins:", err)
	}
}
