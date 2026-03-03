import { lazy } from 'react';

export const routes = [
  {
    path: '/login',
    name: '登录',
    component: lazy(() => import('../pages/Auth/LoginPage')),
  },
  {
    path: '/register',
    name: '注册',
    component: lazy(() => import('../pages/Auth/RegisterPage')),
  },
  {
    path: '/profile',
    name: '个人中心',
    component: lazy(() => import('../pages/Profile/ProfilePage')),
  },
];

export default routes;
