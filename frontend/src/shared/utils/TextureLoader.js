import { TextureLoader } from 'three';

class TextureManager {
  constructor() {
    this.loader = new TextureLoader();
    this.textures = {};
  }

  // 加载纹理
  loadTexture(texturePath) {
    if (this.textures[texturePath]) {
      return this.textures[texturePath];
    }

    // 在实际项目中，这里会加载真实的纹理文件
    // 由于我们没有实际的纹理文件，返回null作为占位符
    // 在真实项目中，应使用loader.load(texturePath)来加载纹理
    return null;
  }

  // 预加载纹理
  preloadTexture(name, texturePath) {
    return new Promise((resolve, reject) => {
      this.loader.load(
        texturePath,
        (texture) => {
          this.textures[name] = texture;
          resolve(texture);
        },
        undefined,
        (err) => {
          console.error(`Error loading texture ${texturePath}:`, err);
          reject(err);
        }
      );
    });
  }

  // 获取已加载的纹理
  getTexture(name) {
    return this.textures[name];
  }
}

export default new TextureManager();