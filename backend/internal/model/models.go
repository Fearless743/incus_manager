package model

import "time"

type InstanceConfig struct {
	Name          string    `json:"name"`
	Image         string    `json:"image"`
	Project       string    `json:"project"`
	Ports         []int     `json:"ports"`
	CPU           int       `json:"cpu"`
	Memory        int64     `json:"memory"`
	Disk          int64     `json:"disk"`
	NetworkLimit  string    `json:"network_limit"`
	UploadLimit   string    `json:"upload_limit"`
	DownloadLimit string    `json:"download_limit"`
	ExpiryDate    time.Time `json:"expiry_date"`
	HostID        uint      `json:"host_id"`
	UserID        uint      `json:"user_id"`
	MappingIP     string    `json:"mapping_ip"`
}

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;size:100;not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	Role         string    `json:"role" gorm:"size:50;default:user"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Host struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;size:100;not null"`
	UserID      uint      `json:"user_id" gorm:"index;not null"`
	Address     string    `json:"address" gorm:"size:255;not null"`
	Certificate string    `json:"certificate" gorm:"type:text"`
	Project     string    `json:"project" gorm:"size:100;not null"`
	CreatedAt   time.Time `json:"created_at"`
}

type Instance struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"uniqueIndex;size:100;not null"`
	HostID        uint      `json:"host_id" gorm:"index;not null"`
	UserID        uint      `json:"user_id" gorm:"index;not null"`
	Image         string    `json:"image" gorm:"size:100;not null"`
	Ports         string    `json:"ports" gorm:"type:text"`
	CPU           int       `json:"cpu"`
	Memory        int64     `json:"memory"`
	Disk          int64     `json:"disk"`
	NetworkLimit  string    `json:"network_limit" gorm:"size:50"`
	UploadLimit   string    `json:"upload_limit" gorm:"size:50"`
	DownloadLimit string    `json:"download_limit" gorm:"size:50"`
	Status        string    `json:"status" gorm:"size:50;default:created"`
	SharedWith    string    `json:"shared_with" gorm:"type:text"`
	ExpiryDate    time.Time `json:"expiry_date"`
	MappingIP     string    `json:"mapping_ip" gorm:"size:45"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
