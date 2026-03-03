import React from 'react';
import { Link } from 'react-router-dom';
import Button from '../../ui/Button/Button';

const Header = () => {
  return (
    <header className="bg-white shadow-md">
      <div className="container mx-auto px-4 py-3 flex justify-between items-center">
        <Link to="/" className="text-xl font-bold text-blue-600">
          MathFun
        </Link>
        
        <nav className="hidden md:flex space-x-6">
          <Link to="/" className="text-gray-700 hover:text-blue-600 transition-colors">
            首页
          </Link>
          <Link to="/examples" className="text-gray-700 hover:text-blue-600 transition-colors">
            示例
          </Link>
          <Link to="/learning" className="text-gray-700 hover:text-blue-600 transition-colors">
            学习
          </Link>
        </nav>
        
        <div className="flex items-center space-x-3">
          <Button variant="outline" size="sm">
            登录
          </Button>
          <Button variant="primary" size="sm">
            注册
          </Button>
        </div>
      </div>
    </header>
  );
};

export default Header;