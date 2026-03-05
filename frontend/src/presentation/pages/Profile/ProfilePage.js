import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../../shared/hooks/useRedux';
import Button from '../../components/ui/Button/Button';
import userService from '../../../data/api/services/userService';

/**
 * 个人中心页面
 */
const ProfilePage = () => {
  const navigate = useNavigate();
  const { isAuthenticated, logout, setCurrentTenant } = useAuth();
  const [isLoading, setIsLoading] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [submitError, setSubmitError] = useState('');
  const [submitSuccess, setSubmitSuccess] = useState('');
  
  // 租户相关状态
  const [tenants, setTenants] = useState([]);
  const [currentTenantId, setCurrentTenantIdLocal] = useState('');
  const [showCreateTenant, setShowCreateTenant] = useState(false);
  const [newTenantName, setNewTenantName] = useState('');
  const [createTenantLoading, setCreateTenantLoading] = useState(false);
  const [createTenantError, setCreateTenantError] = useState('');

  // 用户数据
  const [profileData, setProfileData] = useState({
    nickname: '',
    email: '',
    avatar: '',
    phone: '',
    bio: ''
  });

  // 加载用户信息和租户列表
  useEffect(() => {
    loadUserProfile();
    loadUserTenants();
  }, []);

  const loadUserProfile = async () => {
    try {
      const response = await userService.getInfo();
      const userData = response.data?.data || response.data;
      
      if (userData) {
        setProfileData({
          nickname: userData.nickname || '测试用户',
          email: userData.email || '',
          avatar: userData.avatar !== null && userData.avatar !== undefined ? userData.avatar : '',
          phone: userData.phone !== null && userData.phone !== undefined ? userData.phone : '',
          bio: userData.bio !== null && userData.bio !== undefined ? userData.bio : ''
        });
      }
    } catch (error) {
      console.error('加载用户信息失败:', error);
    }
  };

  const loadUserTenants = async () => {
    try {
      const response = await userService.getUserTenants();
      if (response.data?.data) {
        setTenants(response.data.data);
        // 如果没有当前租户，选择第一个
        const savedTenantId = userService.getCurrentTenantId();
        if (!savedTenantId && response.data.data.length > 0) {
          handleSelectTenant(response.data.data[0].id);
        } else if (savedTenantId) {
          setCurrentTenantIdLocal(savedTenantId);
        }
      }
    } catch (error) {
      console.error('加载租户列表失败:', error);
    }
  };

  const handleSelectTenant = (tenantId) => {
    setCurrentTenantIdLocal(tenantId);
    userService.selectTenant(tenantId);
    setCurrentTenant(tenantId);
  };

  const handleCreateTenant = async () => {
    if (!newTenantName.trim()) {
      setCreateTenantError('请输入租户名称');
      return;
    }

    setCreateTenantLoading(true);
    setCreateTenantError('');

    try {
      await userService.createTenant({ 
        name: newTenantName.trim(),
        maxMembers: 10 // 默认允许 10 个成员
      });
      setCreateTenantLoading(false);
      setShowCreateTenant(false);
      setNewTenantName('');
      // 重新加载租户列表
      loadUserTenants();
    } catch (error) {
      setCreateTenantLoading(false);
      setCreateTenantError(error.message || '创建租户失败，请稍后重试');
    }
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleSave = async () => {
    setIsLoading(true);
    setSubmitError('');
    setSubmitSuccess('');

    try {
      await userService.updateProfile({
        nickname: profileData.nickname,
        phone: profileData.phone,
        bio: profileData.bio
      });
      
      setSubmitSuccess('资料保存成功！');
      setIsEditing(false);
    } catch (error) {
      setSubmitError(error.message || '保存失败，请稍后重试');
    } finally {
      setIsLoading(false);
    }
  };

  if (!isAuthenticated) {
    navigate('/login');
    return null;
  }

  return (
    <div className="min-h-screen bg-secondary py-8">
      <div className="max-w-3xl mx-auto px-page">
        {/* 头部信息卡片 */}
        <div className="card mb-card overflow-hidden">
          <div className="bg-gradient-to-r from-primary to-primary-dark h-32"></div>
          <div className="px-6 pb-6">
            <div className="flex justify-between items-end -mt-12 mb-4">
              <div className="w-24 h-24 bg-white rounded-full p-1 shadow-lg flex items-center justify-center text-4xl">
                {profileData.avatar || '👤'}
              </div>
              <Button 
                variant="outline" 
                size="sm"
                onClick={() => setIsEditing(!isEditing)}
              >
                {isEditing ? '取消' : '编辑资料'}
              </Button>
            </div>
            
            <h1 className="text-title text-text-primary">{profileData.nickname}</h1>
            <p className="text-body text-text-secondary mt-1">{profileData.email}</p>
            
            {/* 租户区域 */}
            <div className="mt-4 pt-4 border-t border-border">
              <div className="flex justify-between items-center mb-2">
                <label className="block text-sm font-medium text-text-primary">
                  我的租户
                </label>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowCreateTenant(!showCreateTenant)}
                >
                  {showCreateTenant ? '取消' : '+ 创建租户'}
                </Button>
              </div>
              
              {/* 创建租户表单 */}
              {showCreateTenant && (
                <div className="mb-3 p-3 bg-secondary rounded-lg">
                  <input
                    type="text"
                    value={newTenantName}
                    onChange={(e) => setNewTenantName(e.target.value)}
                    placeholder="输入租户名称"
                    className="input w-full max-w-xs mb-2"
                    onKeyPress={(e) => e.key === 'Enter' && handleCreateTenant()}
                  />
                  <div className="flex space-x-2">
                    <Button
                      variant="primary"
                      size="sm"
                      onClick={handleCreateTenant}
                      disabled={createTenantLoading}
                    >
                      {createTenantLoading ? '创建中...' : '创建'}
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => {
                        setShowCreateTenant(false);
                        setNewTenantName('');
                      }}
                    >
                      取消
                    </Button>
                  </div>
                  {createTenantError && (
                    <p className="text-small text-error mt-2">{createTenantError}</p>
                  )}
                </div>
              )}
              
              {/* 租户列表 - 有租户时显示选择器 */}
              {tenants.length > 0 ? (
                <select
                  value={currentTenantId}
                  onChange={(e) => handleSelectTenant(e.target.value)}
                  className="input w-full max-w-xs"
                >
                  {tenants.map(tenant => (
                    <option key={tenant.id} value={tenant.id}>
                      {tenant.name} ({tenant.role})
                    </option>
                  ))}
                </select>
              ) : (
                <p className="text-small text-text-secondary">
                  暂无租户，请创建第一个租户
                </p>
              )}
            </div>
          </div>
        </div>

        {/* 成功提示 */}
        {submitSuccess && (
          <div className="success-message mb-card">
            {submitSuccess}
          </div>
        )}

        {/* 错误提示 */}
        {submitError && (
          <div className="error-message mb-card">
            {submitError}
          </div>
        )}

        {/* 个人信息表单 */}
        {isEditing && (
          <div className="card mb-card">
            <h2 className="text-xl font-semibold text-text-primary mb-4">编辑个人资料</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-body font-medium text-text-primary mb-1">
                  昵称
                </label>
                <input
                  type="text"
                  value={profileData.nickname}
                  onChange={(e) => setProfileData({...profileData, nickname: e.target.value})}
                  className="input"
                />
              </div>
              
              <div>
                <label className="block text-body font-medium text-text-primary mb-1">
                  手机号
                </label>
                <input
                  type="tel"
                  value={profileData.phone}
                  onChange={(e) => setProfileData({...profileData, phone: e.target.value})}
                  className="input"
                  placeholder="请输入手机号"
                />
              </div>

              <div>
                <label className="block text-body font-medium text-text-primary mb-1">
                  个人简介
                </label>
                <textarea
                  value={profileData.bio}
                  onChange={(e) => setProfileData({...profileData, bio: e.target.value})}
                  rows={3}
                  className="input"
                  placeholder="介绍一下自己..."
                />
              </div>

              <div className="flex space-x-3">
                <Button 
                  variant="primary" 
                  onClick={handleSave}
                  disabled={isLoading}
                >
                  {isLoading ? '保存中...' : '保存'}
                </Button>
                <Button 
                  variant="outline" 
                  onClick={() => setIsEditing(false)}
                  disabled={isLoading}
                >
                  取消
                </Button>
              </div>
            </div>
          </div>
        )}

        {/* 账号安全 */}
        <div className="card mb-card">
          <h2 className="text-lg font-semibold text-text-primary mb-4">账号与安全</h2>
          <div className="divide-y divide-secondary">
            <div className="py-4 flex items-center justify-between">
              <div>
                <h3 className="text-body font-medium text-text-primary">邮箱地址</h3>
                <p className="text-small text-text-secondary mt-1">{profileData.email}</p>
              </div>
              <Button variant="outline" size="sm">修改</Button>
            </div>

            <div className="py-4 flex items-center justify-between">
              <div>
                <h3 className="text-body font-medium text-text-primary">密码</h3>
                <p className="text-small text-text-secondary mt-1">••••••••</p>
              </div>
              <Button variant="outline" size="sm">修改</Button>
            </div>

            <div className="py-4 flex items-center justify-between">
              <div>
                <h3 className="text-body font-medium text-text-primary">手机绑定</h3>
                <p className="text-small text-text-secondary mt-1">
                  {profileData.phone ? profileData.phone : '未绑定'}
                </p>
              </div>
              <Button variant="outline" size="sm">
                {profileData.phone ? '修改' : '绑定'}
              </Button>
            </div>
          </div>
        </div>

        {/* 设置选项 */}
        <div className="card mb-card">
          <h2 className="text-lg font-semibold text-text-primary mb-4">设置</h2>
          <div className="divide-y divide-secondary">
            <div className="py-4 flex items-center justify-between">
              <span className="text-body text-text-secondary">通知提醒</span>
              <label className="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" defaultChecked className="sr-only peer" />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>

            <div className="py-4 flex items-center justify-between">
              <span className="text-body text-text-secondary">深色模式</span>
              <label className="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" className="sr-only peer" />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
          </div>
        </div>

        {/* 退出登录按钮 */}
        <Button 
          variant="danger" 
          fullWidth 
          size="lg"
          onClick={handleLogout}
          className="shadow-md hover:shadow-lg transition-shadow"
        >
          退出登录
        </Button>
      </div>
    </div>
  );
};

export default ProfilePage;
