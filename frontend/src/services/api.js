import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api';

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

export const authAPI = {
  login: (username, password) => api.post('/login', { username, password }),
};

export const userAPI = {
  create: (username, email, password) => api.post('/users', { username, email, password }),
};

export const hostAPI = {
  add: (name, address, certificate) => api.post('/hosts', { name, address, certificate }),
  getAll: () => api.get('/hosts'),
};

export const instanceAPI = {
  create: (config) => api.post('/instances', config),
  getAll: () => api.get('/instances'),
};

export const shareAPI = {
  share: (instanceId, userId, expiresAt) => api.post('/share', { instanceId, userId, expiresAt }),
  revoke: (instanceId, userId) => api.delete('/share', { data: { instanceId, userId } }),
};

export default api;
