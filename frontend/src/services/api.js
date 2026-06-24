import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests if available
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle 401 errors - redirect to login
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  login: (username, password) => api.post('/login', { username, password }),
};

export const userAPI = {
  create: (username, email, password) => api.post('/users', { username, email, password }),
};

export const hostAPI = {
  add: (name, address, certificate) => api.post('/hosts', { name, address, certificate }),
  test: (address, certificate) => api.post('/hosts/test', { address, certificate }),
  update: (id, data) => api.put(`/hosts/${id}`, data),
  delete: (id) => api.delete(`/hosts/${id}`),
  getAll: () => api.get('/hosts'),
};

export const instanceAPI = {
  create: (config) => api.post('/instances', config),
  getAll: () => api.get('/instances'),
  delete: (id) => api.delete(`/instances/${id}`),
  start: (id) => api.post(`/instances/start/${id}`),
  stop: (id) => api.post(`/instances/stop/${id}`),
  getImages: () => api.get('/instances/images'),
};

export const shareAPI = {
  share: (instanceId, userId, expiresAt) => api.post('/share', { instance_id: instanceId, user_id: userId, expires_at: expiresAt }),
  revoke: (instanceId, userId) => api.delete(`/share/${instanceId}/${userId}`),
};

export const statsAPI = {
  get: () => api.get('/stats'),
};

export default api;
