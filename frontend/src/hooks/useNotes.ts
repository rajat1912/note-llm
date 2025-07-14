import { useState, useEffect } from 'react';
import { Note, CreateNoteRequest, UpdateNoteRequest } from '../types';
import { notesAPI } from '../lib/api';
import { toast } from 'react-hot-toast';

export const useNotes = () => {
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchNotes = async () => {
    try {
      setLoading(true);
      const fetchedNotes = await notesAPI.getAllNotes();
      setNotes(fetchedNotes);
      setError(null);
    } catch (err) {
      setError('Failed to fetch notes');
      console.error('Error fetching notes:', err);
    } finally {
      setLoading(false);
    }
  };

  const createNote = async (noteData: CreateNoteRequest) => {
    try {
      const newNote = await notesAPI.createNote(noteData);
      setNotes(prev => [newNote, ...prev]);
      toast.success('Note created successfully!');
      return newNote;
    } catch (err) {
      toast.error('Failed to create note');
      throw err;
    }
  };

  const updateNote = async (id: string, noteData: UpdateNoteRequest) => {
    try {
      const updatedNote = await notesAPI.updateNote(id, noteData);
      setNotes(prev => prev.map(note => 
        note.id === id ? updatedNote : note
      ));
      toast.success('Note updated successfully!');
      return updatedNote;
    } catch (err) {
      toast.error('Failed to update note');
      throw err;
    }
  };

  const deleteNote = async (id: string) => {
    try {
      await notesAPI.deleteNote(id);
      setNotes(prev => prev.filter(note => note.id !== id));
      toast.success('Note deleted successfully!');
    } catch (err) {
      toast.error('Failed to delete note');
      throw err;
    }
  };

  useEffect(() => {
    fetchNotes();
  }, []);

  return {
    notes,
    loading,
    error,
    createNote,
    updateNote,
    deleteNote,
    refetch: fetchNotes,
  };
};