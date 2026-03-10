import React, { useState, useEffect} from 'react';
import { useNavigate } from 'react-router-dom';
import userService from '../../../data/api/services/userService';
import Button from '../../components/ui/Button/Button';
import Input from '../../components/ui/Input/Input';
import Card from '../../components/ui/Card';

/**
 * 租户管理页面
 */
const TenantManagementPage = () => {
  const navigate = useNavigate();
  const [tenants, setTenants] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    maxMembers: 100
  });

  // 加载租户列表
  useEffect(() => {
   loadTenants();
  }, []);

  const loadTenants = async () => {
    try {
     setIsLoading(true);
     setError(null);
     const data = await userService.getUserTenants();
     setTenants(data || []);
    } catch (err) {
     setError('加载租户列表失败');
     console.error('加载租户失败:', err);
    } finally {
     setIsLoading(false);
    }
  };

  const handleInputChange = (e) => {
   setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleCreateTenant = async (e) => {
    e.preventDefault();
    try {
     setError(null);
     const newTenant= await userService.createTenant(formData);
     setTenants([...tenants, newTenant]);
     setShowCreateForm(false);
     setFormData({
        name: '',
        description: '',
        maxMembers: 100
      });
    } catch (err) {
     setError('创建租户失败');
     console.error('创建租户失败:', err);
    }
  };

  const handleSelectTenant = (tenantId) => {
   userService.selectTenant(tenantId);
    // 跳转到首页或其他页面
    navigate('/');
  };

  if (isLoading) {
   return (
      <div className="min-h-screen bg-secondary flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-text-secondary">加载中...</p>
        </div>
      </div>
    );
  }

 return (
    <div className="min-h-screen bg-secondary py-12 px-page">
      <div className="max-w-6xl mx-auto">
        {/* 页面头部 */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-text-primary">租户管理</h1>
            <p className="mt-2 text-text-secondary">管理和选择您的租户</p>
          </div>
          <Button
            variant="primary"
            onClick={() => setShowCreateForm(!showCreateForm)}
          >
            {showCreateForm ? '取消创建' : '+ 创建租户'}
          </Button>
        </div>

        {/* 错误提示 */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-800">{error}</p>
          </div>
        )}

        {/* 创建租户表单 */}
        {showCreateForm && (
          <Card className="mb-8">
            <h2 className="text-xl font-semibold mb-4">创建新租户</h2>
            <form onSubmit={handleCreateTenant} className="space-y-4">
              <Input
                label="租户名称"
                name="name"
                value={formData.name}
                onChange={handleInputChange}
                placeholder="请输入租户名称"
               required
              />
              
              <Input
                label="描述"
                name="description"
               type="textarea"
                value={formData.description}
                onChange={handleInputChange}
                placeholder="请输入租户描述"
                rows={3}
              />
              
              <Input
                label="最大成员数"
                name="maxMembers"
               type="number"
                value={formData.maxMembers}
                onChange={handleInputChange}
                min={1}
               required
              />
              
              <div className="flex gap-4">
                <Button type="submit" variant="primary">
                  创建租户
                </Button>
                <Button type="button" variant="secondary" onClick={() => setShowCreateForm(false)}>
                  取消
                </Button>
              </div>
            </form>
          </Card>
        )}

        {/* 租户列表 */}
        {tenants.length === 0 ? (
          <Card className="text-center py-12">
            <p className="text-text-secondary">暂无租户，请创建一个新租户</p>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {tenants.map((tenant) => (
              <Card key={tenant.id} className="hover:shadow-lg transition-shadow">
                <div className="flex justify-between items-start mb-4">
                  <h3 className="text-lg font-semibold text-text-primary">{tenant.name}</h3>
                  <span className={`px-2 py-1 rounded text-xs font-medium ${
                    tenant.role === 'owner' 
                      ? 'bg-purple-100 text-purple-800' 
                      : 'bg-blue-100 text-blue-800'
                  }`}>
                    {tenant.role === 'owner' ? '创建者' : '成员'}
                  </span>
                </div>
                
                {tenant.description && (
                  <p className="text-sm text-text-secondary mb-4">{tenant.description}</p>
                )}
                
                <div className="space-y-2 text-sm text-text-secondary mb-4">
                  <div className="flex justify-between">
                    <span>成员数：</span>
                    <span>{tenant.memberCount} / {tenant.maxMembers}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>创建时间：</span>
                    <span>{new Date(tenant.createdAt).toLocaleDateString()}</span>
                  </div>
                </div>
                
                <Button
                  variant="primary"
                  fullWidth
                  onClick={() => handleSelectTenant(tenant.id)}
                >
                  选择此租户
                </Button>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default TenantManagementPage;
