package service

import (
	"errors"
	"strings"

	"incus-manager/internal/model"

	"gorm.io/gorm"
)

type HostService struct {
	DB           *gorm.DB
	IncusFactory *IncusServiceFactory
}

func NewHostService(db *gorm.DB, factory *IncusServiceFactory) *HostService {
	return &HostService{DB: db, IncusFactory: factory}
}

type HostWithStatus struct {
	model.Host
	Status string `json:"status"`
}

func (s *HostService) AddHost(name, address, certificate string, userID uint) (*model.Host, error) {
	address = normalizeAddress(address)
	host := &model.Host{
		Name:        name,
		Address:     address,
		Certificate: certificate,
		UserID:      userID,
	}

	host.Project = generateProjectName(name, userID)

	if err := s.DB.Create(host).Error; err != nil {
		return nil, errors.New("failed to add host")
	}

	return host, nil
}

func (s *HostService) TestHost(address, certificate string) (bool, string, error) {
	address = normalizeAddress(address)
	client := NewIncusClient(address, certificate, "")
	err := client.Ping()
	if err != nil {
		return false, err.Error(), nil
	}
	return true, "连接成功", nil
}

func (s *HostService) UpdateHost(hostID, userID uint, name, address, certificate string) (*model.Host, error) {
	address = normalizeAddress(address)
	var host model.Host
	if err := s.DB.First(&host, hostID).Error; err != nil {
		return nil, errors.New("主机不存在")
	}

	if host.UserID != userID {
		return nil, errors.New("无权操作")
	}

	host.Name = name
	host.Address = address
	host.Certificate = certificate
	if address != "" {
		host.Project = generateProjectName(name, userID)
	}

	if err := s.DB.Save(&host).Error; err != nil {
		return nil, errors.New("更新主机失败")
	}

	return &host, nil
}

func (s *HostService) DeleteHost(hostID, userID uint) error {
	var host model.Host
	if err := s.DB.First(&host, hostID).Error; err != nil {
		return errors.New("主机不存在")
	}

	if host.UserID != userID {
		return errors.New("无权操作")
	}

	return s.DB.Delete(&host).Error
}

func (s *HostService) GetHostsByUser(userID uint) ([]HostWithStatus, error) {
	var hosts []model.Host
	if err := s.DB.Where("user_id = ?", userID).Find(&hosts).Error; err != nil {
		return nil, errors.New("failed to get hosts")
	}

	result := make([]HostWithStatus, 0, len(hosts))
	for _, h := range hosts {
		withStatus := HostWithStatus{Host: h}
		client := s.IncusFactory.GetClient(h.ID, h.Address, h.Certificate)
		if err := client.Ping(); err == nil {
			withStatus.Status = "active"
		} else {
			withStatus.Status = "error"
		}
		result = append(result, withStatus)
	}
	return result, nil
}

func (s *HostService) GetHostByID(id uint) (*model.Host, error) {
	var host model.Host
	if err := s.DB.First(&host, id).Error; err != nil {
		return nil, errors.New("host not found")
	}
	return &host, nil
}

func normalizeAddress(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return addr
	}
	if !strings.Contains(addr, "://") {
		addr = "https://" + addr
	}
	return addr
}
