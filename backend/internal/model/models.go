package model

import "time"

type InstanceConfig struct {
	Name          string   `json:"name"`
	Image         string   `json:"image"`
	Project       string   `json:"project"`
	Ports         []int    `json:"ports"`
	CPU           int      `json:"cpu"`
	Memory        int64    `json:"memory"`
	Disk          int64    `json:"disk"`
	NetworkLimit  string   `json:"network_limit"`
	UploadLimit   string   `json:"upload_limit"`
	DownloadLimit string   `json:"download_limit"`
	ExpiryDate    time.Time `json:"expiry_date"`
	HostID        uint     `json:"host_id"`
	UserID        uint     `json:"user_id"`
	MappingIP     string   `json:"mapping_ip"`
}

type User struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Host struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	UserID      uint      `json:"user_id"`
	Address     string    `json:"address"`
	Certificate string    `json:"certificate"`
	Project     string    `json:"project"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IncusURL    string    `json:"-"`
}

type Instance struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	HostID         uint      `json:"host_id"`
	UserID         uint      `json:"user_id"`
	Image          string    `json:"image"`
	Ports          []int     `json:"ports"`
	CPU            int       `json:"cpu"`
	Memory         int64     `json:"memory"`
	Disk           int64     `json:"disk"`
	NetworkLimit   string    `json:"network_limit"`
	UploadLimit    string    `json:"upload_limit"`
	DownloadLimit  string    `json:"download_limit"`
	Status         string    `json:"status"`
	SharedWith     []uint    `json:"shared_with"`
	ExpiryDate     time.Time `json:"expiry_date"`
	MappingIP      string    `json:"mapping_ip"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
