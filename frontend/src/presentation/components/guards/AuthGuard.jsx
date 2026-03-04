import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../../../shared/hooks/useRedux';

/**
 * 认证守卫组件
 * 用于保护需要登录才能访问的路由
 */
const AuthGuard = ({ children }) => {
  const { isAuthenticated } = useAuth();
  const location = useLocation();

  if (!isAuthenticated) {
    // 未登录，重定向到登录页，并保存当前路径
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return children;
};

/**
 * 游客守卫组件
 * 用于保护只有未登录用户才能访问的路由（如登录页、注册页）
 * 已登录用户访问时会自动重定向到个人中心
 */
const GuestGuard = ({ children }) => {
  const { isAuthenticated } = useAuth();
  const location = useLocation();

  if (isAuthenticated) {
    // 已登录，重定向到个人中心或来源页面
    const from = location.state?.from?.pathname || '/profile';
    return <Navigate to={from} replace />;
  }

  return children;
};

export { AuthGuard, GuestGuard };
