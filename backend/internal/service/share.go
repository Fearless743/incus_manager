package service

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"incus-manager/internal/model"
)

type SharedService struct {
	DB           *gorm.DB
	IncusService *IncusService
}

func NewSharedService(db *gorm.DB, incus *IncusService) *SharedService {
	return &SharedService{DB: db, IncusService: incus}
}

func (s *SharedService) ShareInstance(instanceID, sharedWithUserID uint, expiresAt time.Time) error {
	var instance model.Instance
	if err := s.DB.First(&instance, instanceID).Error; err != nil {
		return errors.New("instance not found")
	}

	// Add user to shared_with list
	instance.SharedWith = append(instance.SharedWith, sharedWithUserID)
	instance.ExpiryDate = expiresAt

	if err := s.DB.Save(&instance).Error; err != nil {
		return errors.New("failed to share instance")
	}

	return nil
}

func (s *SharedService) RevokeShare(instanceID, sharedWithUserID uint) error {
	var instance model.Instance
	if err := s.DB.First(&instance, instanceID).Error; err != nil {
		return errors.New("instance not found")
	}

	// Remove user from shared_with list
	newSharedWith := []uint{}
	for _, uid := range instance.SharedWith {
		if uid != sharedWithUserID {
			newSharedWith = append(newSharedWith, uid)
		}
	}
	instance.SharedWith = newSharedWith

	if err := s.DB.Save(&instance).Error; err != nil {
		return errors.New("failed to revoke share")
	}

	return nil
}
