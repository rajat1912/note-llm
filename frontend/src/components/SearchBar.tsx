import React from 'react';
import { motion } from 'framer-motion';
import { Search, X } from 'lucide-react';

interface SearchBarProps {
  searchTerm: string;
  onSearchChange: (term: string) => void;
}

const SearchBar: React.FC<SearchBarProps> = ({ searchTerm, onSearchChange }) => {
  return (
    <motion.div
      initial={{ opacity: 0, y: -10 }}
      animate={{ opacity: 1, y: 0 }}
      className="relative"
    >
      <div className="relative">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
        <input
          type="text"
          placeholder="Search notes..."
          value={searchTerm}
          onChange={(e) => onSearchChange(e.target.value)}
          className="input-field pl-10 pr-10 w-full"
        />
        {searchTerm && (
          <motion.button
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            onClick={() => onSearchChange('')}
            className="absolute right-3 top-1/2 transform -translate-y-1/2 p-1 hover:bg-white/10 rounded-full transition-colors"
          >
            <X className="w-4 h-4 text-gray-400" />
          </motion.button>
        )}
      </div>
    </motion.div>
  );
};

export default SearchBar;