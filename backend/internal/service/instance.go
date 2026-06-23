package service

import (
	"errors"

	"gorm.io/gorm"
	"incus-manager/internal/model"
)

type InstanceService struct {
	DB           *gorm.DB
	IncusService *IncusService
	IPManager    *IPManager
}

func NewInstanceService(db *gorm.DB, incus *IncusService, ipManager *IPManager) *InstanceService {
	return &InstanceService{DB: db, IncusService: incus, IPManager: ipManager}
}

func (s *InstanceService) CreateInstance(config model.InstanceConfig, userID uint) (*model.Instance, error) {
	host, err := s.getHostAndCheckAccess(config.HostID, userID)
	if err != nil {
		return nil, err
	}

	// Allocate IP
	mappingIP, err := s.IPManager.AllocateIP(host.ID)
	if err != nil {
		return nil, errors.New("failed to allocate IP")
	}

	instance := &model.Instance{
		Name:           config.Name,
		HostID:         host.ID,
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

	// Create in Incus
	incusConfig := model.InstanceConfig{
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
		s.IPManager.ReleaseIP(host.ID, mappingIP)
		return nil, errors.New("failed to create instance in Incus")
	}

	instance.Status = "created"
	if err := s.DB.Create(instance).Error; err != nil {
		return nil, errors.New("failed to save instance")
	}

	return instance, nil
}

func (s *InstanceService) GetInstancesByUser(userID uint) ([]model.Instance, error) {
	var instances []model.Instance
	if err := s.DB.Where("user_id = ? OR shared_with @> ARRAY[?]", userID, userID).Find(&instances).Error; err != nil {
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

	if err := s.IncusService.DeleteInstance(instance.Name, host.Project); err != nil {
		return err
	}

	s.IPManager.ReleaseIP(host.ID, instance.MappingIP)

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

func (s *InstanceService) StartInstance(instanceID, userID uint) error {
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

	return s.IncusService.StartInstance(instance.Name, host.Project)
}

func (s *InstanceService) StopInstance(instanceID, userID uint) error {
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

	return s.IncusService.StopInstance(instance.Name, host.Project)
}
