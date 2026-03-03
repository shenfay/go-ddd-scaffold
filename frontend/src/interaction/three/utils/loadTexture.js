import { TextureLoader } from 'three';

// Example utility function to load a texture
// Note: This assumes the image file is placed in public/images/ or src/assets/images/
// and imported directly, or served statically.
// For dynamic loading from public/, you'd use the full path relative to public folder.
const loadTexture = (path) => {
  const loader = new TextureLoader();
  // Path should be relative to the 'public' folder, e.g., '/images/my_texture.jpg'
  return loader.load(path);
};

export default loadTexture;