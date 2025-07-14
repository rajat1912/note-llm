import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { User } from '@/types';

const AuthCallback: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { login } = useAuth();

  useEffect(() => {
    const params = new URLSearchParams(window.location.hash.substring(1))
    const token = params.get("token")
    const email = params.get("email")

    if (token && email) {
      try {
        const user: User = {
            email: email
        }
        console.log(user)
        login(token, user);
        navigate('/dashboard');
      } catch (error) {
        console.error('Error parsing user data:', error);
        navigate('/login');
      }
    } else {
      navigate('/login');
    }
  }, [searchParams, login, navigate]);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <div className="w-8 h-8 border-2 border-purple-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
        <p className="text-gray-400">Completing authentication...</p>
      </div>
    </div>
  );
};

export default AuthCallback;