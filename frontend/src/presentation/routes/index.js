import { lazy } from 'react';

export const routes = [
  {
    path: '/',
    name: '首页',
    component: lazy(() => import('../pages/Home/HomePage')),
    exact: true,
  },
  {
    path: '/knowledge-map',
    name: '知识地图',
    component: lazy(() => import('../pages/KnowledgeMap/KnowledgeMapPage')),
  },
  {
    path: '/learning/:domainId/:trunkId',
    name: '学习页面',
    component: lazy(() => import('../pages/Learning/LearningPage')),
  },
  {
    path: '/achievements',
    name: '成就中心',
    component: lazy(() => import('../pages/Achievements/AchievementsPage')),
  },
  {
    path: '/profile',
    name: '个人中心',
    component: lazy(() => import('../pages/Profile/ProfilePage')),
  },
  {
    path: '/parent',
    name: '家长端',
    component: lazy(() => import('../pages/Parent/ParentDashboard')),
  },
  {
    path: '/3d/:sceneType',
    name: '3D场景',
    component: lazy(() => import('../pages/ThreeD/ThreeDScenePage')),
  },
];

export default routes;
