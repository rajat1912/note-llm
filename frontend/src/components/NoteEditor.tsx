import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, Save, FileText } from 'lucide-react';
import { Note, CreateNoteRequest, UpdateNoteRequest } from '../types';

interface NoteEditorProps {
  note?: Note;
  isOpen: boolean;
  onClose: () => void;
  onSave: (data: CreateNoteRequest | UpdateNoteRequest, id?: string) => Promise<void>;
}

const NoteEditor: React.FC<NoteEditorProps> = ({ note, isOpen, onClose, onSave }) => {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (note) {
      setTitle(note.title);
      setContent(note.content);
    } else {
      setTitle('');
      setContent('');
    }
  }, [note]);

  const handleSave = async () => {
    if (!title.trim() && !content.trim()) return;
    
    setIsSaving(true);
    try {
      const noteData = {
        title: title.trim() || 'Untitled',
        content: content.trim(),
      };
      
      await onSave(noteData, note?.id);
      onClose();
    } catch (error) {
      console.error('Error saving note:', error);
    } finally {
      setIsSaving(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.ctrlKey && e.key === 's') {
      e.preventDefault();
      handleSave();
    }
    if (e.key === 'Escape') {
      onClose();
    }
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
          onClick={onClose}
        >
          <motion.div
            initial={{ scale: 0.95, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.95, opacity: 0 }}
            className="glass w-full max-w-4xl h-[80vh] rounded-2xl overflow-hidden"
            onClick={(e) => e.stopPropagation()}
            onKeyDown={handleKeyDown}
          >
            {/* Header */}
            <div className="flex items-center justify-between p-6 border-b border-white/10">
              <div className="flex items-center space-x-3">
                <FileText className="w-5 h-5 text-purple-400" />
                <h2 className="text-lg font-semibold text-white">
                  {note ? 'Edit Note' : 'New Note'}
                </h2>
              </div>
              <div className="flex items-center space-x-3">
                <motion.button
                  onClick={handleSave}
                  disabled={isSaving}
                  className="btn-primary flex items-center space-x-2 disabled:opacity-50"
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                >
                  <Save className="w-4 h-4" />
                  <span>{isSaving ? 'Saving...' : 'Save'}</span>
                </motion.button>
                <motion.button
                  onClick={onClose}
                  className="p-2 glass glass-hover rounded-lg"
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                >
                  <X className="w-5 h-5 text-gray-400" />
                </motion.button>
              </div>
            </div>

            {/* Content */}
            <div className="flex flex-col h-full p-6 space-y-4">
              <input
                type="text"
                placeholder="Note title..."
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="input-field text-xl font-semibold bg-transparent border-none focus:ring-0 p-0"
                autoFocus
              />
              
              <textarea
                placeholder="Start writing your note..."
                value={content}
                onChange={(e) => setContent(e.target.value)}
                className="input-field flex-1 resize-none bg-transparent border-none focus:ring-0 p-0 font-mono text-sm leading-relaxed"
              />
              
              <div className="text-xs text-gray-500 flex justify-between">
                <span>Press Ctrl+S to save, Esc to close</span>
                <span>{content.length} characters</span>
              </div>
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};

export default NoteEditor;