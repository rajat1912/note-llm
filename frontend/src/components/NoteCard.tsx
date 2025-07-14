import React from 'react';
import { motion } from 'framer-motion';
import { Edit3, Trash2, Calendar } from 'lucide-react';
import { Note } from '../types';

interface NoteCardProps {
  note: Note;
  onEdit: (note: Note) => void;
  onDelete: (id: string) => void;
}

const NoteCard: React.FC<NoteCardProps> = ({ note, onEdit, onDelete }) => {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  const truncateContent = (content: string, maxLength: number = 150) => {
    if (content.length <= maxLength) return content;
    return content.substring(0, maxLength) + '...';
  };

  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -20 }}
      whileHover={{ y: -5 }}
      className="note-card group"
    >
      <div className="flex justify-between items-start mb-3">
        <h3 className="text-lg font-semibold text-white truncate pr-2">
          {note.title || 'Untitled'}
        </h3>
        <div className="flex space-x-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
          <motion.button
            onClick={(e) => {
              e.stopPropagation();
              onEdit(note);
            }}
            className="p-1.5 glass glass-hover rounded-md"
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
          >
            <Edit3 className="w-4 h-4 text-purple-400" />
          </motion.button>
          <motion.button
            onClick={(e) => {
              e.stopPropagation();
              onDelete(note.id);
            }}
            className="p-1.5 glass glass-hover rounded-md hover:bg-red-500/20"
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.9 }}
          >
            <Trash2 className="w-4 h-4 text-red-400" />
          </motion.button>
        </div>
      </div>

      <div
        className="text-gray-300 text-sm mb-4 cursor-pointer"
        onClick={() => onEdit(note)}
      >
        {truncateContent(note.content)}
      </div>

      <div className="flex items-center text-xs text-gray-500">
        <Calendar className="w-3 h-3 mr-1" />
        <span>Modified {formatDate(note.modified_at)}</span>
      </div>
    </motion.div>
  );
};

export default NoteCard;