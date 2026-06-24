package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func (h *Handler) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /api/login", h.login)
	mux.HandleFunc("POST /api/users", h.createUser)

	// Protected routes
	auth := middleware.Authenticate(h.authService)
	mux.HandleFunc("POST /api/hosts", auth(h.addHost))
	mux.HandleFunc("GET /api/hosts", auth(h.getHosts))
	mux.HandleFunc("GET /api/instances", auth(h.getInstances))
	mux.HandleFunc("POST /api/instances", auth(h.createInstance))
	mux.HandleFunc("DELETE /api/instances/", auth(h.deleteInstance))
	mux.HandleFunc("POST /api/instances/start/", auth(h.startInstance))
	mux.HandleFunc("POST /api/instances/stop/", auth(h.stopInstance))
	mux.HandleFunc("POST /api/share", auth(h.shareInstance))
	mux.HandleFunc("DELETE /api/share/", auth(h.revokeShare))
	mux.HandleFunc("GET /api/instances/images", auth(h.getImages))
	mux.HandleFunc("GET /api/stats", auth(h.getStats))

	return mux
}

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
	instanceID, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		writeError(w, "Invalid instance ID", http.StatusBadRequest)
		return
	}

	err = h.instanceService.DeleteInstance(uint(instanceID), userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Instance deleted"})
}

func (h *Handler) startInstance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	instanceID, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		writeError(w, "Invalid instance ID", http.StatusBadRequest)
		return
	}

	err = h.instanceService.StartInstance(uint(instanceID), userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Instance started"})
}

func (h *Handler) stopInstance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	instanceID, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
	if err != nil {
		writeError(w, "Invalid instance ID", http.StatusBadRequest)
		return
	}

	err = h.instanceService.StopInstance(uint(instanceID), userID)
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
	instanceID, err := strconv.ParseUint(r.PathValue("instanceId"), 10, 32)
	if err != nil {
		writeError(w, "Invalid instance ID", http.StatusBadRequest)
		return
	}

	sharedWithUserID, err := strconv.ParseUint(r.PathValue("userId"), 10, 32)
	if err != nil {
		writeError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.sharedService.RevokeShare(uint(instanceID), uint(sharedWithUserID))
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
