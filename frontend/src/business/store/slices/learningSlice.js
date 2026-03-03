/**
 * 学习进度状态 Slice
 * 
 * 管理学习相关的全局状态：课程、进度、成就等
 */

import { createSlice } from '@reduxjs/toolkit';

const learningSlice = createSlice({
  name: 'learning',
  initialState: {
    // 当前课程
    currentCourse: null,
    
    // 课程列表
    courses: [],
    
    // 学习进度
    progress: {
      totalLessons: 0,
      completedLessons: 0,
      currentLevel: 1,
      experiencePoints: 0
    },
    
    // 成就
    achievements: [],
    
    // 知识图谱数据
    knowledgeGraph: null,
    
    // 加载状态
    isLoading: false,
    
    // 错误状态
    error: null
  },
  
  reducers: {
    // 设置当前课程
    setCurrentCourse: (state, action) => {
      state.currentCourse = action.payload;
    },
    
    // 设置课程列表
    setCourses: (state, action) => {
      state.courses = action.payload;
    },
    
    // 更新学习进度
    updateProgress: (state, action) => {
      state.progress = {
        ...state.progress,
        ...action.payload
      };
    },
    
    // 完成课程
    completeCourse: (state, action) => {
      const courseId = action.payload;
      const course = state.courses.find(c => c.id === courseId);
      
      if (course) {
        course.completed = true;
        state.progress.completedLessons++;
      }
    },
    
    // 添加成就
    addAchievement: (state, action) => {
      const achievement = action.payload;
      
      if (!state.achievements.find(a => a.id === achievement.id)) {
        state.achievements.push(achievement);
      }
    },
    
    // 获得经验值
    earnExperience: (state, action) => {
      state.progress.experiencePoints += action.payload;
      
      // 检查升级（每 1000 经验值升一级）
      const newLevel = Math.floor(state.progress.experiencePoints / 1000) + 1;
      if (newLevel > state.progress.currentLevel) {
        state.progress.currentLevel = newLevel;
      }
    },
    
    // 设置知识图谱
    setKnowledgeGraph: (state, action) => {
      state.knowledgeGraph = action.payload;
    },
    
    // 设置加载状态
    setLoading: (state, action) => {
      state.isLoading = action.payload;
    },
    
    // 设置错误
    setError: (state, action) => {
      state.error = action.payload;
    },
    
    // 清除错误
    clearError: (state) => {
      state.error = null;
    }
  }
});

export const {
  setCurrentCourse,
  setCourses,
  updateProgress,
  completeCourse,
  addAchievement,
  earnExperience,
  setKnowledgeGraph,
  setLoading,
  setError,
  clearError
} = learningSlice.actions;

export default learningSlice.reducer;
