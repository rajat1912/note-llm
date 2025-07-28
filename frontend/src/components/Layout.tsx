import React, { useState } from 'react';
import { motion } from 'framer-motion';
import AskAIOverlay from './AskAIOverlay';
import { notesAPI } from '../lib/api';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [showAI, setShowAI] = useState(false);

  return (
    <div className="min-h-screen bg-gradient-to-br from-dark-950 via-dark-900 to-dark-950">
      {/* Animated background elements */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <motion.div
          className="absolute -top-40 -right-40 w-80 h-80 bg-purple-500/10 rounded-full blur-3xl"
          animate={{
            scale: [1, 1.2, 1],
            opacity: [0.3, 0.5, 0.3],
          }}
          transition={{
            duration: 8,
            repeat: Infinity,
            ease: 'easeInOut',
          }}
        />
        <motion.div
          className="absolute -bottom-40 -left-40 w-80 h-80 bg-blue-500/10 rounded-full blur-3xl"
          animate={{
            scale: [1.2, 1, 1.2],
            opacity: [0.2, 0.4, 0.2],
          }}
          transition={{
            duration: 10,
            repeat: Infinity,
            ease: 'easeInOut',
          }}
        />
      </div>

      {/* Main Content */}
      <div className="relative z-10">
        {children}

        {/* Ask AI Floating Button */}
        <button
          onClick={() => setShowAI(true)}
          className="fixed bottom-6 right-6 z-40 px-4 py-2 rounded-full bg-gradient-to-r from-blue-600 to-purple-600 text-white shadow-lg hover:scale-105 transition"
        >
          Ask AI
        </button>
      </div>

      {/* AI Overlay */}
      {showAI && (
        <AskAIOverlay
          onClose={() => setShowAI(false)}
        />
      )}
    </div>
  );
};

export default Layout;
