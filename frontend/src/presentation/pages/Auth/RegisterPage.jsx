import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../../shared/hooks/useRedux';
import Button from '../../ui/Button/Button';
import Input from '../../ui/Input/Input';
import userService from '../../../data/api/services/userService';

/**
 * 注册页面组件
 */
const RegisterPage = () => {
  const navigate = useNavigate();
  const { login, isLoading: authLoading, error: authError, clearError } = useAuth();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState('');
  
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    confirmPassword: '',
    nickname: ''
  });

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
    if (authError || submitError) {
      clearError();
      setSubmitError('');
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // 验证密码匹配
    if (formData.password !== formData.confirmPassword) {
      setSubmitError('两次输入的密码不一致');
      return;
    }

    // 验证密码长度
    if (formData.password.length < 6) {
      setSubmitError('密码长度至少为 6 位');
      return;
    }

    setIsSubmitting(true);
    setSubmitError('');

    try {
      // 调用注册 API
      await userService.register({
        email: formData.email,
        password: formData.password,
        nickname: formData.nickname
      });
      
      // 注册成功后自动登录
      await login(formData.email, formData.password);
      
      // 跳转到个人中心
      navigate('/profile');
    } catch (error) {
      setSubmitError(error.message || '注册失败，请稍后重试');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
          创建新账号
        </h2>
        <p className="mt-2 text-center text-sm text-gray-600">
          已有账号？{' '}
          <Link to="/login" className="font-medium text-blue-600 hover:text-blue-500">
            立即登录
          </Link>
        </p>
      </div>

      <div className="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
        <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
          <form className="space-y-6" onSubmit={handleSubmit}>
            {(authError || submitError) && (
              <div className="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-md text-sm">
                {authError || submitError}
              </div>
            )}

            <Input
              label="昵称"
              name="nickname"
              type="text"
              required
              value={formData.nickname}
              onChange={handleChange}
              placeholder="请输入昵称"
            />

            <Input
              label="邮箱"
              name="email"
              type="email"
              autoComplete="email"
              required
              value={formData.email}
              onChange={handleChange}
              placeholder="请输入邮箱"
            />

            <Input
              label="密码"
              name="password"
              type="password"
              autoComplete="new-password"
              required
              minLength={6}
              value={formData.password}
              onChange={handleChange}
              placeholder="至少 6 位密码"
            />

            <Input
              label="确认密码"
              name="confirmPassword"
              type="password"
              autoComplete="new-password"
              required
              minLength={6}
              value={formData.confirmPassword}
              onChange={handleChange}
              placeholder="再次输入密码"
            />

            <div className="flex items-center">
              <input
                id="terms"
                name="terms"
                type="checkbox"
                required
                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <label htmlFor="terms" className="ml-2 block text-sm text-gray-900">
                我同意{' '}
                <a href="#" className="font-medium text-blue-600 hover:text-blue-500">
                  服务条款和隐私政策
                </a>
              </label>
            </div>

            <Button
              type="submit"
              variant="primary"
              size="lg"
              fullWidth
              disabled={isSubmitting || authLoading}
            >
              {isSubmitting || authLoading ? '注册中...' : '注册'}
            </Button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default RegisterPage;
