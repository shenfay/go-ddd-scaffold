import React, { Suspense } from 'react';
import { BrowserRouter, Routes, Route, Navigate, Link } from 'react-router-dom';
import { routes } from './presentation/routes';

const Navigation = () => (
  <nav style={{
    position: 'fixed',
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'white',
    borderTop: '1px solid #e5e7eb',
    padding: '12px 20px',
    display: 'flex',
    justifyContent: 'space-around',
    zIndex: 1000,
    boxShadow: '0 -2px 10px rgba(0,0,0,0.05)'
  }}>
    <NavLink to="/" icon="🏠" label="首页" />
    <NavLink to="/knowledge-map" icon="🗺️" label="知识地图" />
    <NavLink to="/achievements" icon="🏆" label="成就" />
    <NavLink to="/profile" icon="👤" label="我的" />
  </nav>
);

const NavLink = ({ to, icon, label }) => (
  <Link to={to} style={{
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    textDecoration: 'none',
    color: '#6b7280',
    fontSize: '12px',
    gap: '4px'
  }}>
    <span style={{ fontSize: '20px' }}>{icon}</span>
    <span>{label}</span>
  </Link>
);

const Loading = () => (
  <div style={{
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '100vh',
    fontSize: '18px',
    color: '#6b7280'
  }}>
    加载中...
  </div>
);

function App() {
  return (
    <BrowserRouter>
      <div style={{ paddingBottom: '80px' }}>
        <Routes>
          {routes.map((route) => (
            <Route
              key={route.path}
              path={route.path}
              element={
                <Suspense fallback={<Loading />}>
                  <route.component />
                </Suspense>
              }
            />
          ))}
          <Route path="*" element={<Navigate to="/" />} />
        </Routes>
        <Navigation />
      </div>
    </BrowserRouter>
  );
}

export default App;
