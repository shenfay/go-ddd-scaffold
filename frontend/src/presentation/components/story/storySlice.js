/**
 * 剧情状态管理 Slice
 * 
 * 使用 Redux Toolkit 管理剧情相关的状态
 */

import { createSlice } from '@reduxjs/toolkit';

// 初始状态
const initialState = {
  isActive: false,
  currentScene: 'modern', // 'modern' | 'ancient' | 'discovery'
  userAge: 0,
  progress: 0,
  currentTask: null,
  achievements: [],
  recentAchievements: [],
  storyHistory: []
};

// 剧情 slice
const storySlice = createSlice({
  name: 'story',
  initialState,
  reducers: {
    // 初始化剧情
    initializeStory: (state, action) => {
      state.userAge = action.payload.userAge;
      state.storyHistory = [];
    },
    
    // 开始剧情
    startStory: (state) => {
      state.isActive = true;
      state.currentScene = 'modern';
      state.progress = 0;
      state.currentTask = null;
      state.recentAchievements = [];
    },
    
    // 结束剧情
    endStory: (state) => {
      state.isActive = false;
      state.currentScene = 'modern';
      state.currentTask = null;
      state.recentAchievements = [];
    },
    
    // 设置当前场景
    setScene: (state, action) => {
      state.currentScene = action.payload;
      // 根据场景更新进度
      const sceneProgress = {
        'modern': 0,
        'ancient': 30,
        'discovery': 70
      };
      state.progress = Math.max(state.progress, sceneProgress[action.payload] || 0);
    },
    
    // 设置当前任务
    setCurrentTask: (state, action) => {
      state.currentTask = action.payload;
    },
    
    // 完成任务
    completeTask: (state, action) => {
      const { taskId, result } = action.payload;
      
      // 记录任务完成历史
      state.storyHistory.push({
        type: 'task_completed',
        taskId,
        result,
        timestamp: Date.now()
      });
      
      // 更新进度
      state.progress = Math.min(100, state.progress + 20);
      
      // 清除当前任务
      state.currentTask = null;
      
      // 检查是否解锁成就
      if (state.progress >= 100) {
        const achievementId = 'pythagorean_discovery';
        if (!state.achievements.includes(achievementId)) {
          state.achievements.push(achievementId);
          state.recentAchievements.push({
            id: achievementId,
            name: '毕达哥拉斯发现者',
            description: '成功体验了毕达哥拉斯定理的发现过程',
            icon: '🎓',
            unlockedAt: Date.now()
          });
        }
      }
    },
    
    // 解锁成就
    unlockAchievement: (state, action) => {
      const achievementId = action.payload;
      if (!state.achievements.includes(achievementId)) {
        state.achievements.push(achievementId);
        
        // 添加到最近解锁列表
        const achievementData = {
          'pythagorean_discovery': {
            name: '毕达哥拉斯发现者',
            description: '成功体验了毕达哥拉斯定理的发现过程',
            icon: '🎓'
          }
        };
        
        state.recentAchievements.push({
          id: achievementId,
          ...achievementData[achievementId],
          unlockedAt: Date.now()
        });
      }
    },
    
    // 清除最近解锁的成就
    clearRecentAchievements: (state) => {
      state.recentAchievements = [];
    },
    
    // 重置剧情状态
    resetStory: () => initialState
  }
});

// 导出 actions
export const storyActions = storySlice.actions;

// 导出 selector
export const selectStoryState = (state) => state.story;

// 导出 reducer
export default storySlice.reducer;