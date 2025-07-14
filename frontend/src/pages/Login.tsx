import React, { useEffect } from 'react';
import { motion } from 'framer-motion';
import { PenTool, Sparkles, Shield, Zap, Chrome } from 'lucide-react';
import { authAPI } from '../lib/api';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';

const Login: React.FC = () => {
  const { user } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (user) {
      navigate('/dashboard');
    }
  }, [user, navigate]);

  const handleGoogleLogin = () => {
    authAPI.googleAuth();
  };

  const features = [
    {
      icon: <Sparkles className="w-6 h-6" />,
      title: "Beautiful Interface",
      description: "Clean, modern design with smooth animations"
    },
    {
      icon: <Shield className="w-6 h-6" />,
      title: "Secure & Private",
      description: "Your notes are encrypted and protected"
    },
    {
      icon: <Zap className="w-6 h-6" />,
      title: "Lightning Fast",
      description: "Instant search and real-time synchronization"
    }
  ];

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-6xl grid lg:grid-cols-2 gap-12 items-center">
        {/* Left side - Branding */}
        <motion.div
          initial={{ opacity: 0, x: -50 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6 }}
          className="text-center lg:text-left"
        >
          <motion.div
            className="flex items-center justify-center lg:justify-start space-x-3 mb-8"
            whileHover={{ scale: 1.05 }}
          >
            <div className="p-3 bg-gradient-to-r from-purple-600 to-purple-700 rounded-xl">
              <PenTool className="w-8 h-8 text-white" />
            </div>
            <h1 className="text-4xl font-bold gradient-text">NoteFlow</h1>
          </motion.div>

          <h2 className="text-3xl lg:text-4xl font-bold text-white mb-6">
            Your thoughts,{' '}
            <span className="gradient-text">beautifully organized</span>
          </h2>

          <p className="text-gray-300 text-lg mb-8 leading-relaxed">
            Experience the future of note-taking with our elegant, fast, and secure platform. 
            Capture ideas, organize thoughts, and boost your productivity.
          </p>

          <div className="grid gap-4 mb-8">
            {features.map((feature, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.2 + index * 0.1 }}
                className="flex items-center space-x-4 p-4 glass rounded-lg"
              >
                <div className="text-purple-400">{feature.icon}</div>
                <div>
                  <h3 className="font-semibold text-white">{feature.title}</h3>
                  <p className="text-gray-400 text-sm">{feature.description}</p>
                </div>
              </motion.div>
            ))}
          </div>
        </motion.div>

        {/* Right side - Login */}
        <motion.div
          initial={{ opacity: 0, x: 50 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6, delay: 0.2 }}
          className="flex justify-center"
        >
          <div className="glass glow-border rounded-2xl p-8 w-full max-w-md">
            <div className="text-center mb-8">
              <h3 className="text-2xl font-bold text-white mb-2">Welcome Back</h3>
              <p className="text-gray-400">Sign in to access your notes</p>
            </div>

            <motion.button
              onClick={handleGoogleLogin}
              className="w-full flex items-center justify-center space-x-3 bg-white hover:bg-gray-100 text-gray-900 font-medium py-4 px-6 rounded-lg transition-all duration-200 shadow-lg hover:shadow-xl"
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <Chrome className="w-5 h-5" />
              <span>Continue with Google</span>
            </motion.button>

            <div className="mt-6 text-center">
              <p className="text-xs text-gray-500">
                By signing in, you agree to our Terms of Service and Privacy Policy
              </p>
            </div>
          </div>
        </motion.div>
      </div>
    </div>
  );
};

export default Login;