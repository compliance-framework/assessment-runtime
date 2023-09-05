package controlplane

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	apiUrl string
}

func NewClient(apiUrl string) *Client {
	return &Client{apiUrl: apiUrl}
}

func (c *Client) GetRuntimeConfigurationJobs() ([]*RuntimeConfigurationJob, error) {
	client := resty.New()

	resp, err := client.R().Get(c.apiUrl)
	if err != nil {
		return nil, err
	}

	var jobs []*RuntimeConfigurationJob
	err = json.Unmarshal(resp.Body(), &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}
