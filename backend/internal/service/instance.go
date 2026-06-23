package service

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"incus-manager/internal/model"
)

type InstanceService struct {
	DB           *gorm.DB
	IncusService *IncusService
}

func NewInstanceService(db *gorm.DB, incus *IncusService) *InstanceService {
	return &InstanceService{DB: db, IncusService: incus}
}

func (s *InstanceService) CreateInstance(config model.InstanceConfig, hostID, userID uint) (*model.Instance, error) {
	host, err := s.getHostAndCheckAccess(hostID, userID)
	if err != nil {
		return nil, err
	}

	// Generate mapping IP based on host
	mappingIP := s.generateMappingIP(host)

	instance := &model.Instance{
		Name:           config.Name,
		HostID:         hostID,
		UserID:         userID,
		Image:          config.Image,
		Ports:          config.Ports,
		CPU:            config.CPU,
		Memory:         config.Memory,
		Disk:           config.Disk,
		NetworkLimit:   config.NetworkLimit,
		UploadLimit:    config.UploadLimit,
		DownloadLimit:  config.DownloadLimit,
		Status:         "creating",
		MappingIP:      mappingIP,
		ExpiryDate:     config.ExpiryDate,
	}

	// Create instance in Incus
	incusConfig := service.InstanceConfig{
		Name:          config.Name,
		Image:         config.Image,
		Project:       host.Project,
		Ports:         config.Ports,
		CPU:           config.CPU,
		Memory:        config.Memory,
		Disk:          config.Disk,
		NetworkLimit:  config.NetworkLimit,
		UploadLimit:   config.UploadLimit,
		DownloadLimit: config.DownloadLimit,
	}

	if err := s.IncusService.CreateInstance(incusConfig); err != nil {
		return nil, err
	}

	instance.Status = "created"
	if err := s.DB.Create(instance).Error; err != nil {
		return nil, errors.New("failed to save instance")
	}

	return instance, nil
}

func (s *InstanceService) GetInstancesByUser(userID uint) ([]model.Instance, error) {
	var instances []model.Instance
	if err := s.DB.Where("user_id = ? OR ? IN shared_with", userID, userID).Find(&instances).Error; err != nil {
		return nil, errors.New("failed to get instances")
	}
	return instances, nil
}

func (s *InstanceService) DeleteInstance(instanceID, userID uint) error {
	var instance model.Instance
	if err := s.DB.First(&instance, instanceID).Error; err != nil {
		return errors.New("instance not found")
	}

	if instance.UserID != userID {
		return errors.New("access denied")
	}

	host, err := s.getHostAndCheckAccess(instance.HostID, userID)
	if err != nil {
		return err
	}

	// Delete from Incus
	if err := s.IncusService.DeleteInstance(instance.Name, host.Project); err != nil {
		return err
	}

	// Delete from database
	return s.DB.Delete(&instance).Error
}

func (s *InstanceService) getHostAndCheckAccess(hostID, userID uint) (*model.Host, error) {
	var host model.Host
	if err := s.DB.First(&host, hostID).Error; err != nil {
		return nil, errors.New("host not found")
	}

	if host.UserID != userID {
		return nil, errors.New("access denied")
	}

	return &host, nil
}

func (s *InstanceService) generateMappingIP(host *model.Host) string {
	// Simple IP generation based on host ID and timestamp
	// In production, you'd want a more sophisticated IP management system
	lastOctet := int(time.Now().Unix()) % 254 + 1
	return fmt.Sprintf("10.0.%d.%d", host.ID, lastOctet)
}
