/**
 * 剧情引擎主组件
 * 
 * 管理整个剧情体验的核心逻辑和状态
 */

import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { storyActions, selectStoryState } from './storySlice.js';
import SceneTransition from './SceneTransition.jsx';
import CharacterInteraction from './CharacterInteraction.jsx';
import InteractiveTask from './InteractiveTask.jsx';
import AchievementUnlock from './AchievementUnlock.jsx';

const StoryEngine = ({ children }) => {
  const dispatch = useDispatch();
  const storyState = useSelector(selectStoryState);

  // 初始化剧情状态
  useEffect(() => {
    // 从用户数据中获取年龄信息
    const userAge = 8; // TODO: 从实际用户数据获取
    dispatch(storyActions.initializeStory({ userAge }));
  }, [dispatch]);

  // 处理剧情触发
  const handleStoryTrigger = () => {
    dispatch(storyActions.startStory());
  };

  // 处理场景切换
  const handleSceneChange = (scene) => {
    dispatch(storyActions.setScene(scene));
  };

  // 处理任务完成
  const handleTaskComplete = (taskId, result) => {
    dispatch(storyActions.completeTask({ taskId, result }));
  };

  // 处理成就解锁
  const handleAchievementUnlock = (achievementId) => {
    dispatch(storyActions.unlockAchievement(achievementId));
  };

  // 处理剧情结束
  const handleStoryEnd = () => {
    dispatch(storyActions.endStory());
  };

  // 如果剧情未激活，显示触发按钮
  if (!storyState.isActive) {
    return (
      <div>
        {children}
        <button
          onClick={handleStoryTrigger}
          style={{
            position: 'fixed',
            bottom: '20px',
            right: '20px',
            padding: '12px 24px',
            backgroundColor: '#4f46e5',
            color: 'white',
            border: 'none',
            borderRadius: '8px',
            fontSize: '16px',
            cursor: 'pointer',
            boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
            zIndex: 1000
          }}
        >
          时空穿越 🚀
        </button>
      </div>
    );
  }

  return (
    <div style={{ 
      position: 'fixed',
      top: 0,
      left: 0,
      width: '100vw',
      height: '100vh',
      zIndex: 2000,
      backgroundColor: '#f0f9ff'
    }}>
      {/* 场景切换组件 */}
      <SceneTransition 
        currentScene={storyState.currentScene}
        onSceneChange={handleSceneChange}
      />
      
      {/* NPC交互组件 */}
      <CharacterInteraction 
        scene={storyState.currentScene}
        onTaskStart={(taskId) => console.log('Task started:', taskId)}
      />
      
      {/* 交互任务组件 */}
      <InteractiveTask 
        currentTask={storyState.currentTask}
        onComplete={handleTaskComplete}
      />
      
      {/* 成就解锁组件 */}
      <AchievementUnlock 
        recentUnlocks={storyState.recentAchievements}
        onAchievementShown={() => dispatch(storyActions.clearRecentAchievements())}
      />
      
      {/* 剧情控制按钮 */}
      <div style={{
        position: 'absolute',
        bottom: '20px',
        left: '20px',
        display: 'flex',
        gap: '10px'
      }}>
        <button
          onClick={handleStoryEnd}
          style={{
            padding: '8px 16px',
            backgroundColor: '#ef4444',
            color: 'white',
            border: 'none',
            borderRadius: '6px',
            cursor: 'pointer'
          }}
        >
          退出剧情
        </button>
        
        <div style={{
          padding: '8px 16px',
          backgroundColor: 'rgba(0, 0, 0, 0.1)',
          borderRadius: '6px',
          fontSize: '14px'
        }}>
          进度: {storyState.progress}%
        </div>
      </div>
    </div>
  );
};

export default StoryEngine;