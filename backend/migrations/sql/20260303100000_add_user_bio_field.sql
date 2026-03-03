-- +goose Up
-- SQL section 'Up' is executed when this migration is applied

-- ============================================
-- 为 users 表添加 bio 字段
-- ============================================

-- 添加 bio 字段（个人简介）
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS bio VARCHAR(500);

-- 添加注释
COMMENT ON COLUMN users.bio IS '用户个人简介';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

-- 移除 bio 字段
ALTER TABLE users 
DROP COLUMN IF EXISTS bio;
