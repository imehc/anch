-- 创建日记表
CREATE TABLE IF NOT EXISTS diaries (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    mood INTEGER CHECK (mood >= 1 AND mood <= 5),
    tags TEXT[], -- PostgreSQL 数组类型
    related_bill_id INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_diaries_user_id ON diaries(user_id);
CREATE INDEX IF NOT EXISTS idx_diaries_created_at ON diaries(created_at);
CREATE INDEX IF NOT EXISTS idx_diaries_mood ON diaries(mood);
CREATE INDEX IF NOT EXISTS idx_diaries_tags ON diaries USING GIN(tags); -- GIN 索引用于数组查询

-- 创建文件表（用于存储日记和账单的附件）
CREATE TABLE IF NOT EXISTS files (
    id SERIAL PRIMARY KEY,
    file_type VARCHAR(20) NOT NULL CHECK (file_type IN ('image', 'audio', 'video')),
    file_url VARCHAR(500) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建日记文件关联表
CREATE TABLE IF NOT EXISTS diary_files (
    diary_id INTEGER NOT NULL REFERENCES diaries(id) ON DELETE CASCADE,
    file_id INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    PRIMARY KEY (diary_id, file_id)
);

-- 为日记表添加更新时间触发器
DROP TRIGGER IF EXISTS update_diaries_updated_at ON diaries;
CREATE TRIGGER update_diaries_updated_at
    BEFORE UPDATE ON diaries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
