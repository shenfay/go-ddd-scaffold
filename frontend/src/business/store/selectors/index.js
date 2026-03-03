/**
 * Redux Selectors
 * 
 * 集中管理所有 Redux state 选择器
 * 便于组件中使用 useSelector 获取数据
 */

// ============ 认证相关 ============

export const selectAuthToken = (state) => state.auth.token;
export const selectIsAuthenticated = (state) => state.auth.isAuthenticated;
export const selectAuthLoading = (state) => state.auth.isLoading;
export const selectAuthError = (state) => state.auth.error;

// ============ UI 相关 ============

export const selectModalState = (state) => state.ui.modals;
export const selectIsModalOpen = (state) => state.ui.modals.isOpen;
export const selectModalType = (state) => state.ui.modals.type;
export const selectModalData = (state) => state.ui.modals.data;

export const selectSidebarOpen = (state) => state.ui.sidebar.isOpen;

export const selectNotification = (state) => state.ui.notification;
export const selectIsNotificationOpen = (state) => state.ui.notification.isOpen;

export const selectIsLoading = (state) => state.ui.isLoading;

export const selectTheme = (state) => state.ui.theme;

// ============ 用户信息相关 ============

export const selectUserProfile = (state) => state.user.profile;
export const selectUserLoading = (state) => state.user.isLoading;
export const selectUserError = (state) => state.user.error;

export const selectUserPreferences = (state) => state.user.preferences;
export const selectUserLanguage = (state) => state.user.preferences.language;
export const selectSoundEnabled = (state) => state.user.preferences.soundEnabled;
export const selectNotificationsEnabled = (state) => state.user.preferences.notificationsEnabled;

// ============ 学习进度相关 ============

export const selectCurrentCourse = (state) => state.learning.currentCourse;
export const selectCourses = (state) => state.learning.courses;
export const selectLearningProgress = (state) => state.learning.progress;
export const selectAchievements = (state) => state.learning.achievements;
export const selectKnowledgeGraph = (state) => state.learning.knowledgeGraph;
export const selectLearningLoading = (state) => state.learning.isLoading;
export const selectLearningError = (state) => state.learning.error;

// ============ 派生选择器 ============

/**
 * 获取用户完成进度百分比
 */
export const selectProgressPercentage = (state) => {
  const { totalLessons, completedLessons } = state.learning.progress;
  if (totalLessons === 0) return 0;
  return Math.round((completedLessons / totalLessons) * 100);
};

/**
 * 获取用户当前等级和经验值进度
 */
export const selectLevelProgress = (state) => {
  const { currentLevel, experiencePoints } = state.learning.progress;
  const currentLevelExp = (currentLevel - 1) * 1000;
  const nextLevelExp = currentLevel * 1000;
  const levelProgress = Math.round(
    ((experiencePoints - currentLevelExp) / (nextLevelExp - currentLevelExp)) * 100
  );
  
  return {
    currentLevel,
    experiencePoints,
    levelProgress
  };
};

export default {
  // 认证
  selectAuthToken,
  selectIsAuthenticated,
  selectAuthLoading,
  selectAuthError,
  
  // UI
  selectModalState,
  selectIsModalOpen,
  selectModalType,
  selectModalData,
  selectSidebarOpen,
  selectNotification,
  selectIsNotificationOpen,
  selectIsLoading,
  selectTheme,
  
  // 用户信息
  selectUserProfile,
  selectUserLoading,
  selectUserError,
  selectUserPreferences,
  selectUserLanguage,
  selectSoundEnabled,
  selectNotificationsEnabled,
  
  // 学习进度
  selectCurrentCourse,
  selectCourses,
  selectLearningProgress,
  selectAchievements,
  selectKnowledgeGraph,
  selectLearningLoading,
  selectLearningError,
  
  // 派生
  selectProgressPercentage,
  selectLevelProgress
};
