package database

import (
	"gorm.io/gorm"
	"incus-manager/internal/model"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Host{},
		&model.Instance{},
	)
}

func SeedInitialAdmin(db *gorm.DB) error {
	var count int64
	db.Model(&model.User{}).Count(&count)
	
	if count == 0 {
		admin := model.User{
			Username: "admin",
			Email:    "admin@incus-manager.local",
			Role:     "admin",
		}
		// TODO: Hash password before saving
		return db.Create(&admin).Error
	}
	
	return nil
}
