import React from 'react';
import { motion } from 'framer-motion';
import { LogOut, User, PenTool } from 'lucide-react';
import { useAuth } from '../contexts/AuthContext';

const Header: React.FC = () => {
  const { user, logout } = useAuth();

  return (
    <motion.header
      initial={{ y: -20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      className="glass border-b border-white/10 sticky top-0 z-50"
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <motion.div
            className="flex items-center space-x-3"
            whileHover={{ scale: 1.05 }}
          >
            <div className="p-2 bg-gradient-to-r from-purple-600 to-purple-700 rounded-lg">
              <PenTool className="w-6 h-6 text-white" />
            </div>
            <h1 className="text-xl font-bold gradient-text">NoteFlow</h1>
          </motion.div>

          {/* User menu */}
          <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-3">
              {/* {user?.picture ? (
                <img
                  src={user.picture}
                  alt={user.name}
                  className="w-8 h-8 rounded-full border-2 border-purple-500/30"
                />
              ) : (
                <div className="w-8 h-8 bg-purple-600 rounded-full flex items-center justify-center">
                  <User className="w-4 h-4 text-white" />
                </div>
              )} */}
              <span className="text-sm text-gray-300 hidden sm:block">
                {user?.email || user?.email}
              </span>
            </div>
            
            <motion.button
              onClick={logout}
              className="p-2 glass glass-hover rounded-lg"
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              <LogOut className="w-4 h-4 text-gray-300" />
            </motion.button>
          </div>
        </div>
      </div>
    </motion.header>
  );
};

export default Header;