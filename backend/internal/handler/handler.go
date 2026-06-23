package handler

import (
	"encoding/json"
	"net/http"
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
}

func NewHandler(auth *service.AuthService, user *service.UserService, host *service.HostService, 
	instance *service.InstanceService, shared *service.SharedService) *Handler {
	return &Handler{
		authService:     auth,
		userService:     user,
		hostService:     host,
		instanceService: instance,
		sharedService:   shared,
	}
}

func (h *Handler) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /api/login", h.login)
	mux.HandleFunc("POST /api/users", h.createUser)

	// Protected routes
	protected := middleware.Authenticate(h.authService)(mux)

	protected.HandleFunc("POST /api/hosts", h.addHost)
	protected.HandleFunc("GET /api/hosts", h.getHosts)
	protected.HandleFunc("GET /api/instances", h.getInstances)
	protected.HandleFunc("POST /api/instances", h.createInstance)
	protected.HandleFunc("DELETE /api/instances/{id}", h.deleteInstance)
	protected.HandleFunc("POST /api/share", h.shareInstance)
	protected.HandleFunc("DELETE /api/share/{instanceId}/{userId}", h.revokeShare)
	protected.HandleFunc("GET /api/instances/{id}/logs", h.getInstanceLogs)
	protected.HandleFunc("POST /api/instances/{id}/start", h.startInstance)
	protected.HandleFunc("POST /api/instances/{id}/stop", h.stopInstance)
	protected.HandleFunc("GET /api/instances/images", h.getImages)

	return protected
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

	req.ExpiryDate = time.Now().Add(30 * 24 * time.Hour) // Default 30 days

	instance, err := h.instanceService.CreateInstance(req, req.HostID, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, instance, http.StatusCreated)
}

func (h *Handler) deleteInstance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	instanceID := uint(mux.Vars(r)["id"])

	err := h.instanceService.DeleteInstance(instanceID, userID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Instance deleted"})
}

func (h *Handler) shareInstance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

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
	userID := middleware.GetUserID(r)
	vars := mux.Vars(r)
	instanceID := uint(vars["instanceId"])
	sharedWithUserID := uint(vars["userId"])

	err := h.sharedService.RevokeShare(instanceID, sharedWithUserID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"message": "Share revoked"})
}

func (h *Handler) getInstanceLogs(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement instance logs
	writeJSON(w, map[string]string{"logs": "No logs available"})
}

func (h *Handler) startInstance(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement start instance
	writeJSON(w, map[string]string{"message": "Instance started"})
}

func (h *Handler) stopInstance(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement stop instance
	writeJSON(w, map[string]string{"message": "Instance stopped"})
}

func (h *Handler) getImages(w http.ResponseWriter, r *http.Request) {
	// TODO: Get available images from Incus
	writeJSON(w, []string{})
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
