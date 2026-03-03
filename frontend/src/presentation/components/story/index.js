/**
 * 剧情体验入口文件
 * 
 * 导出所有剧情相关组件和状态管理
 */

// 导出主组件
export { default as StoryEngine } from './StoryEngine.jsx';

// 导出状态管理
export { default as storyReducer, storyActions, selectStoryState } from './storySlice.js';

// 导出子组件（用于单独测试）
export { default as SceneTransition } from './SceneTransition.jsx';
export { default as CharacterInteraction } from './CharacterInteraction.jsx';
export { default as InteractiveTask } from './InteractiveTask.jsx';
export { default as AchievementUnlock } from './AchievementUnlock.jsx';

// 剧情体验配置
export const STORY_CONFIG = {
  // 剧情ID
  PYTHAGOREAN_DISCOVERY: 'pythagorean_discovery',
  
  // 场景配置
  SCENES: {
    MODERN: 'modern',
    ANCIENT: 'ancient',
    DISCOVERY: 'discovery'
  },
  
  // 任务配置
  TASKS: {
    PYTHAGOREAN_MEASUREMENT: 'pythagorean_measurement'
  },
  
  // 成就配置
  ACHIEVEMENTS: {
    PYTHAGOREAN_DISCOVERY: 'pythagorean_discovery',
    EXPLORER: 'explorer',
    MASTER: 'master'
  }
};

// 默认导出
export default {
  StoryEngine,
  storyReducer,
  storyActions,
  selectStoryState,
  STORY_CONFIG
};