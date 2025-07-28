import axios from 'axios';
import { Note, CreateNoteRequest, UpdateNoteRequest } from '../types';

// Configure your backend API base URL
const API_BASE_URL = 'http://localhost:8080'; // Update this to your backend URL

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  googleAuth: () => {
    // Redirect to your backend's Google OAuth endpoint
    window.location.href = `${API_BASE_URL}/auth/google`;
  },
};

export const notesAPI = {
  getAllNotes: async (): Promise<Note[]> => {
    const response = await api.get('/notes');
    return response.data;
  },

  createNote: async (note: CreateNoteRequest): Promise<Note> => {
    const response = await api.post('/notes', note);
    return response.data;
  },

  updateNote: async (id: string, note: UpdateNoteRequest): Promise<Note> => {
    const response = await api.put(`/notes/${id}`, note);
    return response.data;
  },

  deleteNote: async (id: string): Promise<void> => {
    await api.delete(`/notes/${id}`);
  },

  getNote: async (id: string): Promise<Note> => {
    const response = await api.get(`/notes/${id}`);
    return response.data;
  },
  askQuestion: async (question: string): Promise<string> => {
    const response = await api.post('/notes/ask', { question });
    return response.data.answer;
  },
};

export default api;