package main

import (
	"github.com/compliance-framework/assessment-runtime/plugins"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"

	"github.com/compliance-framework/assessment-runtime/config"
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

	// Load assessment configs
	assessmentConfigsPath := filepath.Join(execDir, "assessments")
	err = confManager.LoadAssessmentConfigs(assessmentConfigsPath)
	if err != nil {
		log.Fatalf("Failed to load assessment configs: %s", err)
	}

	// Download plugin packages
	packages, err := confManager.Packages()
	if err != nil {
		log.Fatalf("Failed to get packages: %s", err)
	}

	pluginDownloader := plugins.NewPackageDownloader(cfg.PluginRegistryURL)
	err = pluginDownloader.DownloadPackages(packages)
	if err != nil {
		log.Errorf("Error downloading some of the plugins: %s", err)
		// TODO: If the download error keeps occurring, we should report it back to the control plane.
	}

	//pluginManager := plugins.NewAssessment(cfg)
	//
	//err = pluginManager.Init()
	//if err != nil {
	//	log.Errorf("Error starting plugins: %s", err)
	//}

	//scheduler := plugins.NewScheduler()
	//for _, plugin := range cfg.Plugins {
	//	scheduler.AddJob(plugin.Schedule, func() {
	//		// TODO: This should come from the control plane. We're just simulating it for now.
	//		err := pluginManager.executePlugin(plugin.Name, plugins.ActionInput{
	//			SSPId:        "123",
	//			ControlId:    "123",
	//			ComponentId:  "123",
	//			AssessmentId: "123",
	//		})
	//		if err != nil {
	//			log.Errorf("Error starting plugin %s: %s", plugin.Name, err)
	//			return
	//		}
	//	})
	//}
	//scheduler.Init()
	//
	select {}
}
