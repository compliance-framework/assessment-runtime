package main

import (
	"fmt"
	"os"

	"github.com/compliance-framework/assessment-runtime/config"
)

const configFilePath = "assets/config.yaml"

func main() {
	confManager := config.NewConfigurationManager()

	config, err := confManager.LoadConfig(configFilePath)
	if err != nil {
		fmt.Printf("failed to load config: %s", err)
		os.Exit(1)
	}

	fmt.Printf("config loaded successfully: %v", config)
}
