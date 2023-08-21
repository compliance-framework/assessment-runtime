package plugins

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/compliance-framework/assessment-runtime/config"
)

const (
	PluginPath = "./plugins"
)

type PluginManager struct {
	cfg    config.Config
	client *http.Client
}

func NewPluginManager(cfg config.Config) *PluginManager {
	if _, err := os.Stat(PluginPath); os.IsNotExist(err) {
		err = os.MkdirAll(PluginPath, 0755)
		if err != nil {
			log.Errorf("Failed to create directory: %v", err)
			return nil
		}
	}

	return &PluginManager{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (m *PluginManager) DownloadPlugins() error {
	var wg sync.WaitGroup
	var errorCh = make(chan error)

	for _, plugin := range m.cfg.Plugins {
		wg.Add(1)
		go func(p config.Plugin) {
			defer wg.Done()
			if err := m.downloadPlugin(p); err != nil {
				errorCh <- err
			} else {
				log.Infof("Downloaded plugin: %s", p.Name)
			}
		}(plugin)
	}

	go func() {
		wg.Wait()
		close(errorCh)
	}()

	var errs []error
	for err := range errorCh {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors during download: %v", len(errs), errs)
	}

	return nil
}

func (m *PluginManager) downloadPlugin(p config.Plugin) error {
	pluginPath := filepath.Join(PluginPath, fmt.Sprintf("%s-%s", p.Name, p.Version))

	if _, err := os.Stat(pluginPath); !os.IsNotExist(err) {
		return nil
	}

	resp, err := m.client.Get(fmt.Sprintf("%s/%s/%s/%s", m.cfg.PluginRegistryURL, p.Name, p.Version, p.Name))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(pluginPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
