import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../../shared/hooks/useRedux';
import Button from '../../ui/Button/Button';

const Header = () => {
  const navigate = useNavigate();
  const { isAuthenticated, logout } = useAuth();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <header className="bg-white shadow-md">
      <div className="container mx-auto px-4 py-3 flex justify-between items-center">
        <Link to="/" className="text-xl font-bold text-blue-600">
          Go DDD Scaffold
        </Link>
        
        <nav className="hidden md:flex space-x-6">
          {!isAuthenticated ? (
            <>
              <Link to="/login" className="text-gray-700 hover:text-blue-600 transition-colors">
                首页
              </Link>
              <Link to="/register" className="text-gray-700 hover:text-blue-600 transition-colors">
                示例
              </Link>
            </>
          ) : (
            <Link to="/profile" className="text-gray-700 hover:text-blue-600 transition-colors">
              个人中心
            </Link>
          )}
        </nav>
        
        <div className="flex items-center space-x-3">
          {isAuthenticated ? (
            <>
              <Button variant="outline" size="sm" onClick={handleLogout}>
                退出
              </Button>
              <Button 
                variant="primary" 
                size="sm"
                onClick={() => navigate('/profile')}
              >
                个人中心
              </Button>
            </>
          ) : (
            <>
              <Button 
                variant="outline" 
                size="sm"
                onClick={() => navigate('/login')}
              >
                登录
              </Button>
              <Button 
                variant="primary" 
                size="sm"
                onClick={() => navigate('/register')}
              >
                注册
              </Button>
            </>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header;