import React, { useState, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Plus, Grid, List } from 'lucide-react';
import Header from '../components/Header';
import NoteCard from '../components/NoteCard';
import NoteEditor from '../components/NoteEditor';
import SearchBar from '../components/SearchBar';
import { useNotes } from '../hooks/useNotes';
import { Note, CreateNoteRequest, UpdateNoteRequest } from '../types';

const Dashboard: React.FC = () => {
  const { notes, loading, createNote, updateNote, deleteNote } = useNotes();
  const [isEditorOpen, setIsEditorOpen] = useState(false);
  const [editingNote, setEditingNote] = useState<Note | undefined>();
  const [searchTerm, setSearchTerm] = useState('');
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');

  const filteredNotes = useMemo(() => {
     if (!notes) return [];
    if (!searchTerm) return notes;
    return notes.filter(note =>
      note.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      note.content.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [notes, searchTerm]);

  const handleCreateNote = () => {
    setEditingNote(undefined);
    setIsEditorOpen(true);
  };

  const handleEditNote = (note: Note) => {
    setEditingNote(note);
    setIsEditorOpen(true);
  };

  const handleSaveNote = async (noteData: CreateNoteRequest | UpdateNoteRequest, id?: string) => {
    if (id) {
      await updateNote(id, noteData as UpdateNoteRequest);
    } else {
      await createNote(noteData as CreateNoteRequest);
    }
  };

  const handleDeleteNote = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this note?')) {
      await deleteNote(id);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <motion.div
          animate={{ rotate: 360 }}
          transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
          className="w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full"
        />
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <Header />
      
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Top bar */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
          <div>
            <h1 className="text-3xl font-bold text-white mb-2">My Notes</h1>
            <p className="text-gray-400">
              {filteredNotes.length} {filteredNotes.length === 1 ? 'note' : 'notes'}
              {searchTerm && ` matching "${searchTerm}"`}
            </p>
          </div>
          
          <div className="flex items-center space-x-4">
            <SearchBar searchTerm={searchTerm} onSearchChange={setSearchTerm} />
            
            <div className="flex items-center space-x-2 glass rounded-lg p-1">
              <motion.button
                onClick={() => setViewMode('grid')}
                className={`p-2 rounded-md transition-colors ${
                  viewMode === 'grid' ? 'bg-purple-600 text-white' : 'text-gray-400 hover:text-white'
                }`}
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <Grid className="w-4 h-4" />
              </motion.button>
              <motion.button
                onClick={() => setViewMode('list')}
                className={`p-2 rounded-md transition-colors ${
                  viewMode === 'list' ? 'bg-purple-600 text-white' : 'text-gray-400 hover:text-white'
                }`}
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
              >
                <List className="w-4 h-4" />
              </motion.button>
            </div>
            
            <motion.button
              onClick={handleCreateNote}
              className="btn-primary flex items-center space-x-2"
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <Plus className="w-4 h-4" />
              <span className="hidden sm:inline">New Note</span>
            </motion.button>
          </div>
        </div>

        {/* Notes grid */}
        <AnimatePresence mode="wait">
          {filteredNotes.length === 0 ? (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="text-center py-16"
            >
              <div className="glass rounded-2xl p-12 max-w-md mx-auto">
                <div className="w-16 h-16 bg-purple-600/20 rounded-full flex items-center justify-center mx-auto mb-4">
                  <Plus className="w-8 h-8 text-purple-400" />
                </div>
                <h3 className="text-xl font-semibold text-white mb-2">
                  {searchTerm ? 'No notes found' : 'No notes yet'}
                </h3>
                <p className="text-gray-400 mb-6">
                  {searchTerm 
                    ? `No notes match "${searchTerm}". Try a different search term.`
                    : 'Create your first note to get started with NoteFlow.'}
                </p>
                {!searchTerm && (
                  <motion.button
                    onClick={handleCreateNote}
                    className="btn-primary"
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                  >
                    Create Your First Note
                  </motion.button>
                )}
              </div>
            </motion.div>
          ) : (
            <motion.div
              layout
              className={
                viewMode === 'grid'
                  ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6'
                  : 'space-y-4'
              }
            >
              {filteredNotes.map((note) => (
                <NoteCard
                  key={note.id}
                  note={note}
                  onEdit={handleEditNote}
                  onDelete={handleDeleteNote}
                />
              ))}
            </motion.div>
          )}
        </AnimatePresence>
      </main>

      {/* Note Editor Modal */}
      <NoteEditor
        note={editingNote}
        isOpen={isEditorOpen}
        onClose={() => setIsEditorOpen(false)}
        onSave={handleSaveNote}
      />
    </div>
  );
};

export default Dashboard;
