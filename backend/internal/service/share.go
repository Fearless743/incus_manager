package service

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
	"incus-manager/internal/model"
)

type SharedService struct {
	DB           *gorm.DB
	IncusFactory *IncusServiceFactory
}

func NewSharedService(db *gorm.DB, factory *IncusServiceFactory) *SharedService {
	return &SharedService{DB: db, IncusFactory: factory}
}

func (s *SharedService) ShareInstance(instanceID, sharedWithUserID uint, expiresAt time.Time) error {
	var instance model.Instance
	if err := s.DB.First(&instance, instanceID).Error; err != nil {
		return errors.New("instance not found")
	}

	var sharedWith []uint
	json.Unmarshal([]byte(instance.SharedWith), &sharedWith)

	for _, uid := range sharedWith {
		if uid == sharedWithUserID {
			return errors.New("instance already shared with this user")
		}
	}

	sharedWith = append(sharedWith, sharedWithUserID)
	sharedWithJSON, _ := json.Marshal(sharedWith)
	instance.SharedWith = string(sharedWithJSON)

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

	var sharedWith []uint
	json.Unmarshal([]byte(instance.SharedWith), &sharedWith)

	newSharedWith := []uint{}
	for _, uid := range sharedWith {
		if uid != sharedWithUserID {
			newSharedWith = append(newSharedWith, uid)
		}
	}

	sharedWithJSON, _ := json.Marshal(newSharedWith)
	instance.SharedWith = string(sharedWithJSON)

	if err := s.DB.Save(&instance).Error; err != nil {
		return errors.New("failed to revoke share")
	}

	return nil
}

func (s *SharedService) GetSharedInstances(userID uint) ([]model.Instance, error) {
	var instances []model.Instance
	if err := s.DB.Find(&instances).Error; err != nil {
		return nil, errors.New("failed to get shared instances")
	}
	return instances, nil
}
