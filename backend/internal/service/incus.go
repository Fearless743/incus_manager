package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"incus-manager/internal/model"
)

type IncusService struct {
	URL       string
	Client    *http.Client
}

func NewIncusService(url string) *IncusService {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &IncusService{
		URL: url,
		Client: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}
}

func (s *IncusService) DoRequest(method, path string, body interface{}) (interface{}, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := fmt.Sprintf("%s%s", s.URL, path)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *IncusService) GetInstances(project string) ([]model.Instance, error) {
	path := fmt.Sprintf("/1.0/instances?project=%s&state=all", project)
	result, err := s.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resultMap := result.(map[string]interface{})
	entries, ok := resultMap["entries"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	instances := make([]model.Instance, 0, len(entries))
	for _, entry := range entries {
		entryMap := entry.(map[string]interface{})
		name := entryMap["name"].(string)
		
		// Get detailed instance info
		instancePath := fmt.Sprintf("/1.0/instances/%s?project=%s", name, project)
		instanceResult, err := s.DoRequest("GET", instancePath, nil)
		if err != nil {
			continue
		}

		instance := model.Instance{
			Name: name,
		}
		
		if state, ok := instanceResult.(map[string]interface{})["status"]; ok {
			instance.Status = state.(string)
		}
		
		instances = append(instances, instance)
	}

	return instances, nil
}

func (s *IncusService) CreateInstance(config InstanceConfig) error {
	instanceData := map[string]interface{}{
		"name":          config.Name,
		"type":          "container",
		"ephemeral":     false,
		"config": map[string]interface{}{
			"limits.cpu":          config.CPU,
			"limits.memory":       config.Memory,
			"limits.network":      config.NetworkLimit,
			"limits.disk":         config.Disk,
			"security.privileged": "false",
		},
		"source": map[string]interface{}{
			"server":  config.Image,
			"type":    "image",
			"protocol": "simplestreams",
			"mode":    "pull",
		},
		"project": config.Project,
	}

	// Add port forwarding if ports are specified
	if len(config.Ports) > 0 {
		networkConfig := make(map[string]interface{})
		for i, port := range config.Ports {
			networkConfig[fmt.Sprintf("proxy.%d", i)] = fmt.Sprintf("tcp::%d::-:%d", port, port)
		}
		instanceData["config"] = networkConfig
	}

	path := "/1.0/instances"
	_, err := s.DoRequest("POST", path, instanceData)
	return err
}

func (s *IncusService) DeleteInstance(name, project string) error {
	path := fmt.Sprintf("/1.0/instances/%s?project=%s", name, project)
	_, err := s.DoRequest("DELETE", path, nil)
	return err
}

func (s *IncusService) StartInstance(name, project string) error {
	path := fmt.Sprintf("/1.0/instances/%s/state?project=%s", name, project)
	_, err := s.DoRequest("PUT", path, map[string]interface{}{
		"action": "start",
		"timeout": -1,
		"force": false,
	})
	return err
}

func (s *IncusService) StopInstance(name, project string) error {
	path := fmt.Sprintf("/1.0/instances/%s/state?project=%s", name, project)
	_, err := s.DoRequest("PUT", path, map[string]interface{}{
		"action":  "stop",
		"timeout": -1,
		"force":   false,
	})
	return err
}

func (s *IncusService) GetImages(project string) ([]string, error) {
	path := fmt.Sprintf("/1.0/images?project=%s", project)
	result, err := s.DoRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	resultMap := result.(map[string]interface{})
	entries, ok := resultMap["entries"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	images := make([]string, 0, len(entries))
	for _, entry := range entries {
		entryMap := entry.(map[string]interface{})
		if alias, ok := entryMap["aliases"].([]interface{}); ok && len(alias) > 0 {
			aliasMap := alias[0].(map[string]interface{})
			if name, ok := aliasMap["name"].(string); ok {
				images = append(images, name)
			}
		}
	}

	return images, nil
}
