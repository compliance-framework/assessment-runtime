package plugins

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/compliance-framework/assessment-runtime/config"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type clientConfig struct {
	Host    string `json:"host"`
	Network string `json:"network"`
	Port    int    `json:"port"`
}

// This will hold the process and the grpc client
type pluginWrapper struct {
	pluginConfig config.PluginConfig
	process      *os.Process
	mutex        sync.Mutex
}

type PluginManager struct {
	wrappers map[string]*pluginWrapper
}

func NewPluginManager(cfg config.Config) *PluginManager {
	pluginManager := &PluginManager{
		wrappers: make(map[string]*pluginWrapper),
	}

	for _, plugin := range cfg.Plugins {
		pluginManager.wrappers[plugin.Name] = &pluginWrapper{
			pluginConfig: plugin,
		}
	}

	return pluginManager
}

func (p *PluginManager) InitPlugins() error {
	log.Info("Initializing plugins")
	return nil
}

func (p *PluginManager) StartPlugin(name string) error {
	wrapper, ok := p.wrappers[name]
	if !ok {
		return fmt.Errorf("plugin %s not found", name)
	}

	execPath, err := os.Executable()
	execDir := filepath.Dir(execPath)
	pluginConfig := wrapper.pluginConfig
	binaryPath := filepath.Join(execDir, "plugins", pluginConfig.Name, pluginConfig.Version, pluginConfig.Name)
	log.Info("Starting plugin: ", binaryPath)
	cmd := exec.Command(binaryPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get stdout pipe: %v", err)
	}

	// TODO: We need to handle stderr as well
	//stderr, err := cmd.StderrPipe()
	//if err != nil {
	//	log.Fatalf("Failed to get stderr pipe: %v", err)
	//}

	err = cmd.Start()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			cmd.Process.Kill()
			panic(r)
		}
	}()

	wrapper.process = cmd.Process
	log.Infof("Started plugin: %s", pluginConfig.Name)

	// TODO: Need to tie these back to the wrapper
	doneCtx, ctxCancel := context.WithCancel(context.Background())

	go func() {
		defer ctxCancel()

		err := cmd.Wait()
		if err != nil {
			log.Errorf("Error waiting for plugin: %s", err)
		} else {
			log.Infof("Plugin exited successfully")
		}

		os.Stderr.Sync()
	}()

	linesCh := make(chan string)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			linesCh <- scanner.Text()
		}
	}()

	timeout := time.After(5 * time.Second)
	select {
	case <-timeout:
		return errors.New("plugin timed out")
	case <-doneCtx.Done():
		return errors.New("plugin exited")
	case line := <-linesCh:
		jsonData := []byte(line)
		data := clientConfig{}
		if err := json.Unmarshal(jsonData, &data); err != nil {
			return fmt.Errorf("failed to get plugin network information: %s", err)
		}
		log.Infof("Plugin network information: %d", data.Port)

		// TODO: Create the grpc client and connect to the plugin
	}

	return nil
}

func (p *PluginManager) StopPlugin(name string) error {
	log.Infof("Stopping plugin: %s", name)
	return nil
}
