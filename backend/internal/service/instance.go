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

func (s *InstanceService) CreateInstance(name, image string, ports []int, cpu int, memory, disk int64, 
	networkLimit, uploadLimit, downloadLimit string, expiryDate time.Time, hostID, userID uint) (*model.Instance, error) {
	
	host, err := s.getHostAndCheckAccess(hostID, userID)
	if err != nil {
		return nil, err
	}

	instance := &model.Instance{
		Name:           name,
		HostID:         hostID,
		UserID:         userID,
		Image:          image,
		Ports:          ports,
		CPU:            cpu,
		Memory:         memory,
		Disk:           disk,
		NetworkLimit:   networkLimit,
		UploadLimit:    uploadLimit,
		DownloadLimit:  downloadLimit,
		Status:         "created",
		ExpiryDate:     expiryDate,
	}

	// TODO: Create instance in Incus using project
	if err := s.DB.Create(instance).Error; err != nil {
		return nil, errors.New("failed to create instance")
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
