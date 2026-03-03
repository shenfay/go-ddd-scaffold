// 简单的声音管理器（模拟实现，实际项目中可以使用Web Audio API）
class SoundManager {
  constructor() {
    // 在实际实现中，这里会加载音效文件
    this.sounds = {};
    this.enabled = true;
  }

  // 模拟播放音效
  playSound(soundName) {
    if (!this.enabled) return;

    // 在实际项目中，这里会播放真实的音频文件
    console.log(`Playing sound: ${soundName}`);
    
    // 可以通过Web Audio API或HTML5 Audio实现真实音效
    // 为简化演示，此处仅输出日志
  }

  // 控制音效开关
  toggle(enabled) {
    this.enabled = !!enabled;
  }

  // 预加载音效
  preloadSound(soundName, soundPath) {
    // 实际实现中会预加载音频文件
    this.sounds[soundName] = soundPath;
  }
}

export default new SoundManager();