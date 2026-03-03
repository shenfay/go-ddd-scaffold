/**
 * Redux Hooks
 * 
 * 自定义 hooks，便于在组件中便捷地访问 Redux 状态和 dispatch 操作
 */

import { useDispatch, useSelector } from 'react-redux';
import { useCallback } from 'react';
import * as authSelectors from '../store/selectors/index.js';
import * as uiActions from '../store/slices/uiSlice.js';
import * as authActions from '../store/slices/authSlice.js';
import * as userActions from '../store/slices/userSlice.js';
import * as learningActions from '../store/slices/learningSlice.js';
import { loginUser, logoutUser } from '../store/slices/authSlice.js';
import { fetchUserProfile, updateUserProfile } from '../store/slices/userSlice.js';

// ============ 认证相关 ============

export const useAuth = () => {
  const dispatch = useDispatch();
  const token = useSelector(authSelectors.selectAuthToken);
  const isAuthenticated = useSelector(authSelectors.selectIsAuthenticated);
  const isLoading = useSelector(authSelectors.selectAuthLoading);
  const error = useSelector(authSelectors.selectAuthError);

  return {
    token,
    isAuthenticated,
    isLoading,
    error,
    login: useCallback(
      (email, password) => dispatch(loginUser({ email, password })),
      [dispatch]
    ),
    logout: useCallback(() => dispatch(logoutUser()), [dispatch]),
    clearError: useCallback(() => dispatch(authActions.clearError()), [dispatch])
  };
};

// ============ UI 相关 ============

export const useUIState = () => {
  const dispatch = useDispatch();
  const theme = useSelector(authSelectors.selectTheme);
  const isLoading = useSelector(authSelectors.selectIsLoading);
  const notification = useSelector(authSelectors.selectNotification);
  const modal = useSelector(authSelectors.selectModalState);
  const sidebarOpen = useSelector(authSelectors.selectSidebarOpen);

  return {
    theme,
    isLoading,
    notification,
    modal,
    sidebarOpen,
    setLoading: useCallback(
      (loading) => dispatch(uiActions.setLoading(loading)),
      [dispatch]
    ),
    setTheme: useCallback(
      (theme) => dispatch(uiActions.setTheme(theme)),
      [dispatch]
    ),
    toggleTheme: useCallback(() => dispatch(uiActions.toggleTheme()), [dispatch]),
    showNotification: useCallback(
      (payload) => dispatch(uiActions.showNotification(payload)),
      [dispatch]
    ),
    hideNotification: useCallback(
      () => dispatch(uiActions.hideNotification()),
      [dispatch]
    ),
    openModal: useCallback(
      (payload) => dispatch(uiActions.openModal(payload)),
      [dispatch]
    ),
    closeModal: useCallback(() => dispatch(uiActions.closeModal()), [dispatch]),
    openSidebar: useCallback(() => dispatch(uiActions.openSidebar()), [dispatch]),
    closeSidebar: useCallback(() => dispatch(uiActions.closeSidebar()), [dispatch]),
    toggleSidebar: useCallback(
      () => dispatch(uiActions.toggleSidebar()),
      [dispatch]
    )
  };
};

// ============ 用户信息相关 ============

export const useUserInfo = () => {
  const dispatch = useDispatch();
  const profile = useSelector(authSelectors.selectUserProfile);
  const preferences = useSelector(authSelectors.selectUserPreferences);
  const isLoading = useSelector(authSelectors.selectUserLoading);
  const error = useSelector(authSelectors.selectUserError);

  return {
    profile,
    preferences,
    isLoading,
    error,
    fetchProfile: useCallback(
      () => dispatch(fetchUserProfile()),
      [dispatch]
    ),
    updateProfile: useCallback(
      (userData) => dispatch(updateUserProfile(userData)),
      [dispatch]
    ),
    setLanguage: useCallback(
      (language) => dispatch(userActions.setLanguage(language)),
      [dispatch]
    ),
    toggleSound: useCallback(
      () => dispatch(userActions.toggleSound()),
      [dispatch]
    ),
    toggleNotifications: useCallback(
      () => dispatch(userActions.toggleNotifications()),
      [dispatch]
    ),
    clearError: useCallback(() => dispatch(userActions.clearError()), [dispatch])
  };
};

// ============ 学习进度相关 ============

export const useLearning = () => {
  const dispatch = useDispatch();
  const currentCourse = useSelector(authSelectors.selectCurrentCourse);
  const courses = useSelector(authSelectors.selectCourses);
  const progress = useSelector(authSelectors.selectLearningProgress);
  const achievements = useSelector(authSelectors.selectAchievements);
  const knowledgeGraph = useSelector(authSelectors.selectKnowledgeGraph);
  const isLoading = useSelector(authSelectors.selectLearningLoading);
  const error = useSelector(authSelectors.selectLearningError);

  return {
    currentCourse,
    courses,
    progress,
    achievements,
    knowledgeGraph,
    isLoading,
    error,
    setCurrentCourse: useCallback(
      (course) => dispatch(learningActions.setCurrentCourse(course)),
      [dispatch]
    ),
    setCourses: useCallback(
      (courses) => dispatch(learningActions.setCourses(courses)),
      [dispatch]
    ),
    updateProgress: useCallback(
      (progressData) => dispatch(learningActions.updateProgress(progressData)),
      [dispatch]
    ),
    completeCourse: useCallback(
      (courseId) => dispatch(learningActions.completeCourse(courseId)),
      [dispatch]
    ),
    addAchievement: useCallback(
      (achievement) => dispatch(learningActions.addAchievement(achievement)),
      [dispatch]
    ),
    earnExperience: useCallback(
      (exp) => dispatch(learningActions.earnExperience(exp)),
      [dispatch]
    ),
    setKnowledgeGraph: useCallback(
      (graph) => dispatch(learningActions.setKnowledgeGraph(graph)),
      [dispatch]
    ),
    setError: useCallback(
      (error) => dispatch(learningActions.setError(error)),
      [dispatch]
    ),
    clearError: useCallback(() => dispatch(learningActions.clearError()), [dispatch])
  };
};

export default {
  useAuth,
  useUIState,
  useUserInfo,
  useLearning
};
