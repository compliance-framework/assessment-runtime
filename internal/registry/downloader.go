package registry

import (
	"fmt"
	"github.com/compliance-framework/assessment-runtime/internal/model"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Downloader struct {
	registryURL string
	client      *http.Client
}

func NewPackageDownloader(registryURL string) *Downloader {
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

	return &Downloader{
		registryURL: registryURL,
		client:      &http.Client{},
	}
}

func (m *Downloader) DownloadPackages(packages []model.Package) error {
	var wg sync.WaitGroup
	var errorCh = make(chan error)

	for _, pkg := range packages {
		wg.Add(1)
		go func(p model.Package) {
			defer wg.Done()
			log.WithFields(log.Fields{
				"package": p.Name,
				"version": p.Version,
			}).Info("Downloading package")
			if err := m.downloadPackage(p); err != nil {
				errorCh <- err
			} else {
				log.WithFields(log.Fields{
					"package": p.Name,
					"version": p.Version,
				}).Info("Downloaded package")
			}
		}(pkg)
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

func (m *Downloader) downloadPackage(p model.Package) error {
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

	resp, err := m.client.Get(fmt.Sprintf("%s/%s/%s/%s", m.registryURL, p.Name, p.Version, p.Name))
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
