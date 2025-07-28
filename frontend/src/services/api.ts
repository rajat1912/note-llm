
import { Note, CreateNoteRequest, UpdateNoteRequest } from '@/types';

const API_BASE_URL = 'http://localhost:8080';

class ApiService {
  private getHeaders(token: string) {
    return {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    };
  }

  async getNotes(token: string): Promise<Note[]> {
    const response = await fetch(`${API_BASE_URL}/notes`, {
      headers: this.getHeaders(token),
    });

    if (!response.ok) {
      throw new Error('Failed to fetch notes');
    }

    return response.json();
  }

  async getNote(id: string, token: string): Promise<Note> {
    const response = await fetch(`${API_BASE_URL}/notes/${id}`, {
      headers: this.getHeaders(token),
    });

    if (!response.ok) {
      throw new Error('Failed to fetch note');
    }

    return response.json();
  }

  async createNote(note: CreateNoteRequest, token: string): Promise<Note> {
    const response = await fetch(`${API_BASE_URL}/notes`, {
      method: 'POST',
      headers: this.getHeaders(token),
      body: JSON.stringify(note),
    });

    if (!response.ok) {
      throw new Error('Failed to create note');
    }

    return response.json();
  }

  async updateNote(id: string, note: UpdateNoteRequest, token: string): Promise<Note> {
    const response = await fetch(`${API_BASE_URL}/notes/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(token),
      body: JSON.stringify(note),
    });

    if (!response.ok) {
      throw new Error('Failed to update note');
    }

    return response.json();
  }

  async deleteNote(id: string, token: string): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/notes/${id}`, {
      method: 'DELETE',
      headers: this.getHeaders(token),
    });

    if (!response.ok) {
      throw new Error('Failed to delete note');
    }
  }

  async askQuestion(question: string, token: string): Promise<string> {
    const response = await fetch(`${API_BASE_URL}/ask`, {
      method: 'POST',
      headers: this.getHeaders(token),
      body: JSON.stringify({ question }),
    });

    if (!response.ok) {
      const err = await response.json();
      throw new Error(err.error || 'Failed to get answer');
    }

    const data = await response.json();
    return data.answer;
  }
}
export const apiService = new ApiService();
