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
	URL      string
	Client   *http.Client
	CertFile string
	KeyFile  string
}

func NewIncusService(url, certFile, keyFile string) *IncusService {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &IncusService{
		URL:      url,
		Client:   &http.Client{Transport: tr, Timeout: 30 * time.Second},
		CertFile: certFile,
		KeyFile:  keyFile,
	}
}

func (s *IncusService) doRequest(method, path string, body interface{}) (map[string]interface{}, error) {
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
	req.Header.Set("X-Cert-Hash", "manual")

	resp, err := s.Client.Do(req)
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

func (s *IncusService) GetInstances(project string) ([]model.Instance, error) {
	result, err := s.doRequest("GET", fmt.Sprintf("/1.0/instances?project=%s&state=all", project), nil)
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

func (s *IncusService) CreateInstance(config model.InstanceConfig) error {
	instanceData := map[string]interface{}{
		"name":    config.Name,
		"type":    "container",
		"ephemeral": false,
		"source": map[string]interface{}{
			"server":  config.Image,
			"type":    "image",
			"protocol": "simplestreams",
			"mode":    "pull",
		},
		"project": config.Project,
		"config": map[string]interface{}{
			"limits.cpu":          fmt.Sprintf("%d", config.CPU),
			"limits.memory":       fmt.Sprintf("%dMB", config.Memory),
			"limits.disk":         fmt.Sprintf("%dGB", config.Disk),
			"security.privileged": "false",
		},
	}

	_, err := s.doRequest("POST", "/1.0/instances", instanceData)
	return err
}

func (s *IncusService) DeleteInstance(name, project string) error {
	_, err := s.doRequest("DELETE", fmt.Sprintf("/1.0/instances/%s?project=%s", name, project), nil)
	return err
}

func (s *IncusService) StartInstance(name, project string) error {
	_, err := s.doRequest("PUT", fmt.Sprintf("/1.0/instances/%s/state?project=%s", name, project), map[string]interface{}{
		"action":  "start",
		"timeout": -1,
		"force":   false,
	})
	return err
}

func (s *IncusService) StopInstance(name, project string) error {
	_, err := s.doRequest("PUT", fmt.Sprintf("/1.0/instances/%s/state?project=%s", name, project), map[string]interface{}{
		"action":  "stop",
		"timeout": -1,
		"force":   false,
	})
	return err
}

func (s *IncusService) GetImages(project string) ([]string, error) {
	result, err := s.doRequest("GET", fmt.Sprintf("/1.0/images?project=%s", project), nil)
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

func (s *IncusService) GetHosts() ([]map[string]interface{}, error) {
	result, err := s.doRequest("GET", "/1.0/hosts", nil)
	if err != nil {
		return nil, err
	}

	hosts := make([]map[string]interface{}, 0)
	if entries, ok := result["entries"].([]interface{}); ok {
		for _, entry := range entries {
			if entryMap, ok := entry.(map[string]interface{}); ok {
				hosts = append(hosts, entryMap)
			}
		}
	}

	return hosts, nil
}

func (s *IncusService) GetNetworks(project string) ([]string, error) {
	result, err := s.doRequest("GET", fmt.Sprintf("/1.0/networks?project=%s", project), nil)
	if err != nil {
		return []string{}, err
	}

	networks := make([]string, 0)
	if entries, ok := result["entries"].([]interface{}); ok {
		for _, entry := range entries {
			if entryMap, ok := entry.(map[string]interface{}); ok {
				if name, ok := entryMap["name"].(string); ok {
					networks = append(networks, name)
				}
			}
		}
	}

	return networks, nil
}

func (s *IncusService) GetProject(projectName string) (map[string]interface{}, error) {
	result, err := s.doRequest("GET", fmt.Sprintf("/1.0/projects/%s", projectName), nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *IncusService) CreateProject(projectName, description string) error {
	projectData := map[string]interface{}{
		"name":        projectName,
		"config":      map[string]interface{}{},
		"description": description,
	}

	_, err := s.doRequest("POST", "/1.0/projects", projectData)
	return err
}
