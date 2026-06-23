package service

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"incus-manager/internal/model"
)

type HostService struct {
	DB           *gorm.DB
	IncusService *IncusService
}

func NewHostService(db *gorm.DB, incus *IncusService) *HostService {
	return &HostService{DB: db, IncusService: incus}
}

func (s *HostService) AddHost(name, address, certificate string, userID uint) (*model.Host, error) {
	host := &model.Host{
		Name:        name,
		Address:     address,
		Certificate: certificate,
		UserID:      userID,
		Status:      "active",
	}

	// Generate project name based on user and host
	host.Project = generateProjectName(userID, host.Name)

	if err := s.DB.Create(host).Error; err != nil {
		return nil, errors.New("failed to add host")
	}

	return host, nil
}

func (s *HostService) GetHostsByUser(userID uint) ([]model.Host, error) {
	var hosts []model.Host
	if err := s.DB.Where("user_id = ?", userID).Find(&hosts).Error; err != nil {
		return nil, errors.New("failed to get hosts")
	}
	return hosts, nil
}

func (s *HostService) GetHostByID(id uint) (*model.Host, error) {
	var host model.Host
	if err := s.DB.First(&host, id).Error; err != nil {
		return nil, errors.New("host not found")
	}
	return &host, nil
}

func (s *HostService) DeleteHost(id, userID uint) error {
	var host model.Host
	if err := s.DB.First(&host, id).Error; err != nil {
		return errors.New("host not found")
	}

	if host.UserID != userID {
		return errors.New("access denied")
	}

	return s.DB.Delete(&host).Error
}

func generateProjectName(userID uint, hostName string) string {
	return fmt.Sprintf("host-%s-%d", hostName, userID)
}
