package service

import (
	"fmt"
	"strconv"
	"time"

	"incus-manager/internal/model"
)

type InstanceConfig struct {
	Name          string
	Image         string
	Project       string
	Ports         []int
	CPU           int
	Memory        int64
	Disk          int64
	NetworkLimit  string
	UploadLimit   string
	DownloadLimit string
	ExpiryDate    time.Time
}

func (s *InstanceService) generateMappingIP(host *model.Host) string {
	// Generate IP based on host ID and current time
	lastOctet := int(time.Now().UnixNano()) % 254 + 1
	return fmt.Sprintf("10.0.%d.%d", host.ID, lastOctet)
}

func (s *InstanceService) parseMemoryLimit(limit string) int64 {
	if limit == "unlimited" || limit == "" {
		return 0
	}
	
	value := int64(0)
	unit := "MB"
	
	if len(limit) > 2 {
		unit = limit[len(limit)-2:]
		value, _ = strconv.ParseInt(limit[:len(limit)-2], 10, 64)
	} else {
		value, _ = strconv.ParseInt(limit, 10, 64)
		unit = "MB"
	}
	
	switch unit {
	case "GB":
		return value * 1024
	case "TB":
		return value * 1024 * 1024
	default:
		return value
	}
}
