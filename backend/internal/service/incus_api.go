package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"incus-manager/internal/model"
)

type IncusClient struct {
	URL      string
	Client   *http.Client
	CertData []byte
	KeyData  []byte
}

func NewIncusClient(url, certPEM, keyPEM string) *IncusClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &IncusClient{
		URL:      url,
		Client:   &http.Client{Transport: tr, Timeout: 30 * time.Second},
		CertData: []byte(certPEM),
		KeyData:  []byte(keyPEM),
	}
}

func (c *IncusClient) doRequest(method, path string, body interface{}) (map[string]interface{}, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := fmt.Sprintf("%s%s", c.URL, path)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("API error: %s", string(bodyBytes))
	}

	return result, nil
}

func (c *IncusClient) GetInstances(project string) ([]model.Instance, error) {
	result, err := c.doRequest("GET", fmt.Sprintf("/1.0/instances?project=%s&state=all", project), nil)
	if err != nil {
		return nil, err
	}

	instances := make([]model.Instance, 0)
	if entries, ok := result["entries"].([]interface{}); ok {
		for _, entry := range entries {
			if entryMap, ok := entry.(map[string]interface{}); ok {
				instance := model.Instance{}
				if name, ok := entryMap["name"].(string); ok {
					instance.Name = name
				}
				if status, ok := entryMap["status"].(string); ok {
					instance.Status = status
				}
				instances = append(instances, instance)
			}
		}
	}

	return instances, nil
}

func (c *IncusClient) CreateInstance(config model.InstanceConfig) error {
	instanceData := map[string]interface{}{
		"name":      config.Name,
		"type":      "container",
		"ephemeral": false,
		"source": map[string]interface{}{
			"server":   config.Image,
			"type":     "image",
			"protocol": "simplestreams",
			"mode":     "pull",
		},
		"project": config.Project,
		"config": map[string]interface{}{
			"limits.cpu":            fmt.Sprintf("%d", config.CPU),
			"limits.memory":         fmt.Sprintf("%dMB", config.Memory),
			"limits.root.size":      fmt.Sprintf("%dGB", config.Disk),
			"security.privileged":   "false",
		},
	}

	// Add port forwarding if ports are specified
	if len(config.Ports) > 0 {
		configMap := instanceData["config"].(map[string]interface{})
		for i, port := range config.Ports {
			configMap[fmt.Sprintf("proxy.%d", i)] = fmt.Sprintf("tcp::%d::-:%d", port, port)
		}
	}

	_, err := c.doRequest("POST", "/1.0/instances", instanceData)
	return err
}

func (c *IncusClient) DeleteInstance(name, project string) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/1.0/instances/%s?project=%s", name, project), nil)
	return err
}

func (c *IncusClient) StartInstance(name, project string) error {
	_, err := c.doRequest("PUT", fmt.Sprintf("/1.0/instances/%s/state?project=%s", name, project), map[string]interface{}{
		"action":  "start",
		"timeout": -1,
		"force":   false,
	})
	return err
}

func (c *IncusClient) StopInstance(name, project string) error {
	_, err := c.doRequest("PUT", fmt.Sprintf("/1.0/instances/%s/state?project=%s", name, project), map[string]interface{}{
		"action":  "stop",
		"timeout": -1,
		"force":   false,
	})
	return err
}

func (c *IncusClient) GetImages(project string) ([]string, error) {
	result, err := c.doRequest("GET", fmt.Sprintf("/1.0/images?project=%s", project), nil)
	if err != nil {
		return []string{}, err
	}

	images := make([]string, 0)
	if entries, ok := result["entries"].([]interface{}); ok {
		for _, entry := range entries {
			if entryMap, ok := entry.(map[string]interface{}); ok {
				if properties, ok := entryMap["properties"].(map[string]interface{}); ok {
					if distro, ok := properties["distro_name"].(string); ok {
						if release, ok := properties["release"].(string); ok {
							images = append(images, fmt.Sprintf("%s/%s", distro, release))
						}
					}
				}
			}
		}
	}

	return images, nil
}

// IncusServiceFactory manages per-host Incus clients
type IncusServiceFactory struct {
	clients map[uint]*IncusClient
	mu      sync.RWMutex
}

func NewIncusServiceFactory() *IncusServiceFactory {
	return &IncusServiceFactory{
		clients: make(map[uint]*IncusClient),
	}
}

func (f *IncusServiceFactory) GetClient(hostID uint, address, certificate string) *IncusClient {
	f.mu.RLock()
	if client, ok := f.clients[hostID]; ok {
		f.mu.RUnlock()
		return client
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check after acquiring write lock
	if client, ok := f.clients[hostID]; ok {
		return client
	}

	client := NewIncusClient(address, certificate, "")
	f.clients[hostID] = client
	return client
}
