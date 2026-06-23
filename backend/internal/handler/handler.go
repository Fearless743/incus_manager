package handler

import (
	"encoding/json"
	"net/http"

	"incus-manager/internal/service"
)

type Handler struct {
	authService     *service.AuthService
	userService     *service.UserService
	hostService     *service.HostService
	instanceService *service.InstanceService
	sharedService   *service.SharedService
}

func NewHandler(auth, user, host, instance, shared *service.Service) *Handler {
	return &Handler{
		authService:     auth,
		userService:     user,
		hostService:     host,
		instanceService: instance,
		sharedService:   shared,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/login":
		h.login(w, r)
	case "/api/users":
		h.createUser(w, r)
	case "/api/hosts":
		h.addHost(w, r)
	case "/api/instances":
		h.createInstance(w, r)
	case "/api/share":
		h.shareInstance(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	// TODO: Implement login logic
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement user creation
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) addHost(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement host addition
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) createInstance(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement instance creation
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) shareInstance(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement sharing logic
	w.WriteHeader(http.StatusNotImplemented)
}
