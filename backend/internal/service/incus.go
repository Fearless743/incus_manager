package service

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type IncusService struct {
	URL       string
	Client    *http.Client
}

func NewIncusService(url, cert string) *IncusService {
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

func (s *IncusService) GetInstances(project string) ([]map[string]interface{}, error) {
	resp, err := s.Client.Get(fmt.Sprintf("%s/1.0/instances?project=%s", s.URL, project))
	if err != nil {
		return nil, fmt.Errorf("failed to get instances: %w", err)
	}
	defer resp.Body.Close()

	// Parse response (simplified)
	return nil, fmt.Errorf("not implemented")
}

func (s *IncusService) CreateInstance(config InstanceConfig) error {
	// Implement Incus API call to create instance
	return fmt.Errorf("not implemented")
}

type InstanceConfig struct {
	Name        string
	Image       string
	Project     string
	Ports       []int
	CPU         int
	Memory      int64
	Disk        int64
	NetworkLimit string
	UploadLimit  string
	DownloadLimit string
}
