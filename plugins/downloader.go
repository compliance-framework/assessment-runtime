package plugins

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/compliance-framework/assessment-runtime/config"
)

const (
	PluginPath = "./plugins"
)

type PluginDownloader struct {
	cfg    config.Config
	client *http.Client
}

func NewPluginDownloader(cfg config.Config) *PluginDownloader {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	pluginsPath := filepath.Join(filepath.Dir(ex), "./plugins")

	if _, err := os.Stat(pluginsPath); os.IsNotExist(err) {
		err = os.MkdirAll(pluginsPath, 0755)
		if err != nil {
			log.Errorf("Failed to create directory: %v", err)
			return nil
		}
	}

	return &PluginDownloader{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (m *PluginDownloader) DownloadPlugins() error {
	var wg sync.WaitGroup
	var errorCh = make(chan error)

	for _, plugin := range m.cfg.Plugins {
		wg.Add(1)
		go func(p config.PluginConfig) {
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

func (m *PluginDownloader) downloadPlugin(p config.PluginConfig) error {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	pluginsPath := filepath.Join(filepath.Dir(ex), "./plugins")

	pluginPath := filepath.Join(pluginsPath, fmt.Sprintf("%s/%s", p.Name, p.Version))
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		err = os.MkdirAll(pluginPath, 0755)
		if err != nil {
			log.Errorf("Failed to create directory: %v", err)
			return nil
		}
	}

	resp, err := m.client.Get(fmt.Sprintf("%s/%s/%s/%s", m.cfg.PluginRegistryURL, p.Name, p.Version, p.Name))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(pluginPath + "/" + p.Name)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = os.Chmod(pluginPath+"/"+p.Name, 0755)
	return err
}
