import React, { Suspense, useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate, useLocation } from 'react-router-dom';
import './index.css';
import { routes } from './presentation/routes';
import { AuthGuard, GuestGuard } from './presentation/components/guards/AuthGuard';
import Header from './presentation/components/common/Header/Header';

const Loading = () => (
  <div className="flex justify-center items-center h-screen text-lg text-gray-500">
    加载中...
  </div>
);

// 不需要 Header 的页面路径
const noHeaderPaths = ['/login', '/register'];

function AppContent() {
  const location = useLocation();
  const showHeader = !noHeaderPaths.includes(location.pathname);

  return (
    <div className="min-h-screen bg-gray-50">
      {showHeader && <Header />}
      <Routes>
        {/* 游客路由（登录、注册） */}
        {routes.map((route) => {
          // 个人中心需要认证守卫
          if (route.path === '/profile') {
            return (
              <Route
                key={route.path}
                path={route.path}
                element={
                  <AuthGuard>
                    <Suspense fallback={<Loading />}>
                      <route.component />
                    </Suspense>
                  </AuthGuard>
                }
              />
            );
          }
          
          // 登录/注册页面使用游客守卫（已登录用户不能访问）
          if (route.path === '/login' || route.path === '/register') {
            return (
              <Route
                key={route.path}
                path={route.path}
                element={
                  <GuestGuard>
                    <Suspense fallback={<Loading />}>
                      <route.component />
                    </Suspense>
                  </GuestGuard>
                }
              />
            );
          }
          
          // 其他路由直接渲染
          return (
            <Route
              key={route.path}
              path={route.path}
              element={
                <Suspense fallback={<Loading />}>
                  <route.component />
                </Suspense>
              }
            />
          );
        })}
        <Route path="*" element={<Navigate to="/login" />} />
      </Routes>
    </div>
  );
}

function App() {
  return (
    <BrowserRouter>
      <AppContent />
    </BrowserRouter>
  );
}

export default App;
