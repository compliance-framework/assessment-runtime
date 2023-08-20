package plugins

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/compliance-framework/assessment-runtime/config"
)

type PluginManager struct {
	cfg    config.Config
	client *http.Client
}

func NewPluginManager(cfg config.Config) *PluginManager {
	return &PluginManager{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (m *PluginManager) DownloadPlugins() error {
	var wg sync.WaitGroup
	var errc = make(chan error)

	for _, plugin := range m.cfg.Plugins {
		wg.Add(1)
		go func(p config.Plugin) {
			defer wg.Done()
			if err := m.downloadPlugin(p); err != nil {
				errc <- err
			}
		}(plugin)
	}
	// create a goroutine that will close the errc channel
	// after all downloader goroutines are done.
	go func() {
		wg.Wait()
		close(errc)
	}()
	// consume the errc channel
	for err := range errc {
		if err != nil {
			return err // return early on error
		}
	}
	return nil
}

func (m *PluginManager) downloadPlugin(p config.Plugin) error {
	pluginPath := filepath.Join("./plugins", fmt.Sprintf("%s-%s.plugins", p.Name, p.Version))
	if _, err := os.Stat(pluginPath); !os.IsNotExist(err) {
		// File exists, no need to download
		return nil
	}

	resp, err := m.client.Get(fmt.Sprintf("%s/%s/%s", m.cfg.PluginRegistryURL, p.Name, p.Version))
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
