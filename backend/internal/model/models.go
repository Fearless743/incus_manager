package model

import "time"

type User struct {
	ID              uint      `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	Role            string    `json:"role"` // admin, user
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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
	NetworkLimit   string    `json:"network_limit"` // "unlimited" or specific limit
	UploadLimit    string    `json:"upload_limit"`
	DownloadLimit  string    `json:"download_limit"`
	Status         string    `json:"status"`
	SharedWith     []uint    `json:"shared_with"`
	ExpiryDate     time.Time `json:"expiry_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ShareRequest struct {
	InstanceID uint  `json:"instance_id"`
	UserID     uint  `json:"user_id"`
	ExpiresAt  time.Time `json:"expires_at"`
}
