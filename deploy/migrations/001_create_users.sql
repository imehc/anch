-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar_url VARCHAR(255),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    disabled_reason TEXT,
    disabled_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);

-- 插入测试用户 (密码是 "password123" 的 bcrypt hash)
-- 密码: password123
INSERT INTO users (username, email, password_hash, role, status)
VALUES (
    'testuser',
    'test@example.com',
    '$2a$10$GzoTFT2X7Ty5wBPUD.0GXe6mmFGfKe6DQK49in7n461Vm4TQP9Mbq',
    'user',
    'active'
) ON CONFLICT (username) DO NOTHING;

-- 插入管理员用户
-- 密码: admin123
INSERT INTO users (username, email, password_hash, role, status)
VALUES (
    'admin',
    'admin@example.com',
    '$2a$10$yBfKvmTLK7FOLgcbvcfVO.hfPVlw0alFKfTnRm3VrKVHt./9Hsbia',
    'admin',
    'active'
) ON CONFLICT (username) DO NOTHING;

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为用户表添加更新时间触发器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
