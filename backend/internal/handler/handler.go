package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"incus-manager/internal/middleware"
	"incus-manager/internal/model"
	"incus-manager/internal/service"
)

type Handler struct {
	authService     *service.AuthService
	userService     *service.UserService
	hostService     *service.HostService
	instanceService *service.InstanceService
	sharedService   *service.SharedService
	ipManager       *service.IPManager
}

func NewHandler(auth *service.AuthService, user *service.UserService, host *service.HostService,
	instance *service.InstanceService, shared *service.SharedService, ipMgr *service.IPManager) *Handler {
	return &Handler{
		authService:     auth,
		userService:     user,
		hostService:     host,
		instanceService: instance,
		sharedService:   shared,
		ipManager:       ipMgr,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) { h.login(w, r) }
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) { h.createUser(w, r) }
func (h *Handler) AddHost(w http.ResponseWriter, r *http.Request) { h.addHost(w, r) }
func (h *Handler) TestHost(w http.ResponseWriter, r *http.Request) { h.testHost(w, r) }
func (h *Handler) UpdateHost(w http.ResponseWriter, r *http.Request) { h.updateHost(w, r) }
func (h *Handler) DeleteHost(w http.ResponseWriter, r *http.Request) { h.deleteHost(w, r) }
func (h *Handler) GetHosts(w http.ResponseWriter, r *http.Request) { h.getHosts(w, r) }
func (h *Handler) GetInstances(w http.ResponseWriter, r *http.Request) { h.getInstances(w, r) }
func (h *Handler) CreateInstance(w http.ResponseWriter, r *http.Request) { h.createInstance(w, r) }
func (h *Handler) DeleteInstance(w http.ResponseWriter, r *http.Request) { h.deleteInstance(w, r) }
func (h *Handler) StartInstance(w http.ResponseWriter, r *http.Request) { h.startInstance(w, r) }
func (h *Handler) StopInstance(w http.ResponseWriter, r *http.Request) { h.stopInstance(w, r) }
func (h *Handler) ShareInstance(w http.ResponseWriter, r *http.Request) { h.shareInstance(w, r) }
func (h *Handler) RevokeShare(w http.ResponseWriter, r *http.Request) { h.revokeShare(w, r) }
func (h *Handler) GetImages(w http.ResponseWriter, r *http.Request) { h.getImages(w, r) }
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) { h.getStats(w, r) }

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		writeError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	writeJSON(w, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, user, http.StatusCreated)
}

func (h *Handler) addHost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req struct {
		Name        string `json:"name"`
		Address     string `json:"address"`
		Certificate string `json:"certificate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	host, err := h.hostService.AddHost(req.Name, req.Address, req.Certificate, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, host, http.StatusCreated)
}

func (h *Handler) testHost(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Address     string `json:"address"`
		Certificate string `json:"certificate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ok, msg, err := h.hostService.TestHost(req.Address, req.Certificate)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		writeError(w, msg, http.StatusBadGateway)
		return
	}

	writeJSON(w, map[string]interface{}{
		"success": true,
		"message": msg,
	})
}

func (h *Handler) updateHost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	hostID := parseUint(strings.TrimPrefix(r.URL.Path, "/api/hosts/"))
	if hostID == 0 {
		writeError(w, "无效的主机ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Address     string `json:"address"`
		Certificate string `json:"certificate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	host, err := h.hostService.UpdateHost(hostID, userID, req.Name, req.Address, req.Certificate)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, host)
}

func (h *Handler) deleteHost(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	hostID := parseUint(strings.TrimPrefix(r.URL.Path, "/api/hosts/"))
	if hostID == 0 {
		writeError(w, "无效的主机ID", http.StatusBadRequest)
		return
	}

	err := h.hostService.DeleteHost(hostID, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "主机已删除"})
}

func (h *Handler) getHosts(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	hosts, err := h.hostService.GetHostsByUser(userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, hosts)
}

func (h *Handler) getInstances(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	instances, err := h.instanceService.GetInstancesByUser(userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, instances)
}

func (h *Handler) createInstance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req model.InstanceConfig
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ExpiryDate.IsZero() {
		req.ExpiryDate = time.Now().Add(30 * 24 * time.Hour)
	}

	instance, err := h.instanceService.CreateInstance(req, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, instance, http.StatusCreated)
}

func (h *Handler) deleteInstance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	instanceID := getInstanceIDFromPath(r.URL.Path)
	if instanceID == 0 {
		writeError(w, "Invalid instance ID", http.StatusBadRequest)
		return
	}

	err := h.instanceService.DeleteInstance(instanceID, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Instance deleted"})
}

func (h *Handler) startInstance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	instanceID := getInstanceIDFromPath(r.URL.Path)
	if instanceID == 0 {
		writeError(w, "Invalid instance ID", http.StatusBadRequest)
		return
	}

	err := h.instanceService.StartInstance(instanceID, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Instance started"})
}

func (h *Handler) stopInstance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	instanceID := getInstanceIDFromPath(r.URL.Path)
	if instanceID == 0 {
		writeError(w, "Invalid instance ID", http.StatusBadRequest)
		return
	}

	err := h.instanceService.StopInstance(instanceID, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Instance stopped"})
}

func (h *Handler) shareInstance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		InstanceID uint   `json:"instance_id"`
		UserID     uint   `json:"user_id"`
		ExpiresAt  string `json:"expires_at"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err != nil {
		writeError(w, "Invalid expires_at format", http.StatusBadRequest)
		return
	}

	err = h.sharedService.ShareInstance(req.InstanceID, req.UserID, expiresAt)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Instance shared"})
}

func (h *Handler) revokeShare(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	instanceID := uint(0)
	sharedWithUserID := uint(0)
	for i, seg := range parts {
		if i == 0 && seg == "share" {
			continue
		}
		id := parseUint(seg)
		if id == 0 {
			writeError(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		if instanceID == 0 {
			instanceID = id
		} else {
			sharedWithUserID = id
		}
	}

	if instanceID == 0 || sharedWithUserID == 0 {
		writeError(w, "Missing instance ID or user ID", http.StatusBadRequest)
		return
	}

	err := h.sharedService.RevokeShare(instanceID, sharedWithUserID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Share revoked"})
}

func (h *Handler) getImages(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, []string{"ubuntu/22.04", "ubuntu/20.04", "centos/8", "centos/9"})
}

func (h *Handler) getStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	instances, _ := h.instanceService.GetInstancesByUser(userID)
	hosts, _ := h.hostService.GetHostsByUser(userID)

	writeJSON(w, map[string]interface{}{
		"total_hosts":        len(hosts),
		"total_instances":    len(instances),
		"running_instances":  0,
		"shared_instances":   0,
	})
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode ...int) {
	status := http.StatusOK
	if len(statusCode) > 0 {
		status = statusCode[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func splitPath(s string) []string {
	result := []string{}
	current := ""
	for _, c := range s {
		if c == '/' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func getInstanceIDFromPath(path string) uint {
	parts := splitPath(path)
	for i, p := range parts {
		if p == "instances" && i+1 < len(parts) {
			return parseUint(parts[i+1])
		}
	}
	return 0
}

func parseUint(s string) uint {
	var n uint64
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return 0
	}
	return uint(n)
}
