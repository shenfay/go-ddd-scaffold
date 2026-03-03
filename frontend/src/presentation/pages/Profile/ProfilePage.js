import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../../shared/hooks/useRedux';
import Button from '../../ui/Button/Button';

/**
 * 个人中心页面 - 通用版本
 */
const ProfilePage = () => {
  const navigate = useNavigate();
  const { isAuthenticated, logout, user } = useAuth();
  const [isEditing, setIsEditing] = useState(false);

  // 模拟用户数据（实际应从 API 获取）
  const [profileData, setProfileData] = useState({
    nickname: '测试用户',
    email: 'user@example.com',
    avatar: '',
    phone: '',
    bio: ''
  });

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleSave = () => {
    // TODO: 调用 API 保存资料
    setIsEditing(false);
  };

  if (!isAuthenticated) {
    // 未登录时跳转到登录页
    navigate('/login');
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* 头部信息卡片 */}
        <div className="bg-white rounded-xl shadow-md overflow-hidden mb-6">
          <div className="bg-gradient-to-r from-blue-500 to-blue-600 h-32"></div>
          <div className="px-6 pb-6">
            <div className="flex justify-between items-end -mt-12 mb-4">
              <div className="w-24 h-24 bg-white rounded-full p-1 shadow-lg">
                <div className="w-full h-full bg-gray-200 rounded-full flex items-center justify-center text-4xl">
                  {profileData.avatar || '👤'}
                </div>
              </div>
              <Button 
                variant="outline" 
                size="sm"
                onClick={() => setIsEditing(!isEditing)}
              >
                {isEditing ? '取消' : '编辑资料'}
              </Button>
            </div>
            
            <h1 className="text-2xl font-bold text-gray-900">{profileData.nickname}</h1>
            <p className="text-gray-600 mt-1">{profileData.email}</p>
          </div>
        </div>

        {/* 个人信息表单 */}
        {isEditing && (
          <div className="bg-white rounded-xl shadow-md p-6 mb-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">编辑个人资料</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  昵称
                </label>
                <input
                  type="text"
                  value={profileData.nickname}
                  onChange={(e) => setProfileData({...profileData, nickname: e.target.value})}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  手机号
                </label>
                <input
                  type="tel"
                  value={profileData.phone}
                  onChange={(e) => setProfileData({...profileData, phone: e.target.value})}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  placeholder="请输入手机号"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  个人简介
                </label>
                <textarea
                  value={profileData.bio}
                  onChange={(e) => setProfileData({...profileData, bio: e.target.value})}
                  rows={3}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  placeholder="介绍一下自己..."
                />
              </div>

              <div className="flex space-x-3">
                <Button variant="primary" onClick={handleSave}>
                  保存
                </Button>
                <Button variant="outline" onClick={() => setIsEditing(false)}>
                  取消
                </Button>
              </div>
            </div>
          </div>
        )}

        {/* 账号安全 */}
        <div className="bg-white rounded-xl shadow-md mb-6">
          <div className="divide-y divide-gray-200">
            <div className="px-6 py-4 flex items-center justify-between">
              <div>
                <h3 className="text-sm font-medium text-gray-900">邮箱地址</h3>
                <p className="text-sm text-gray-500 mt-1">{profileData.email}</p>
              </div>
              <Button variant="outline" size="sm">修改</Button>
            </div>

            <div className="px-6 py-4 flex items-center justify-between">
              <div>
                <h3 className="text-sm font-medium text-gray-900">密码</h3>
                <p className="text-sm text-gray-500 mt-1">••••••••</p>
              </div>
              <Button variant="outline" size="sm">修改</Button>
            </div>

            <div className="px-6 py-4 flex items-center justify-between">
              <div>
                <h3 className="text-sm font-medium text-gray-900">手机绑定</h3>
                <p className="text-sm text-gray-500 mt-1">
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
        <div className="bg-white rounded-xl shadow-md mb-6">
          <div className="px-6 py-4">
            <h2 className="text-lg font-semibold text-gray-900">设置</h2>
          </div>
          <div className="divide-y divide-gray-200">
            <div className="px-6 py-4 flex items-center justify-between">
              <span className="text-sm text-gray-700">通知提醒</span>
              <label className="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" defaultChecked className="sr-only peer" />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>

            <div className="px-6 py-4 flex items-center justify-between">
              <span className="text-sm text-gray-700">深色模式</span>
              <label className="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" className="sr-only peer" />
                <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
              </label>
            </div>
          </div>
        </div>

        {/* 退出登录按钮 */}
        <div className="mt-6">
          <Button 
            variant="danger" 
            fullWidth 
            size="lg"
            onClick={handleLogout}
          >
            退出登录
          </Button>
        </div>
      </div>
    </div>
  );
};

export default ProfilePage;
