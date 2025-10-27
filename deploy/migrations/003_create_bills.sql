-- 创建账单表
CREATE TABLE IF NOT EXISTS bills (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    original_amount DECIMAL(12, 2) NOT NULL,
    discount_amount DECIMAL(12, 2) DEFAULT 0,
    actual_amount DECIMAL(12, 2) NOT NULL,
    discount_type VARCHAR(50),
    payment_method VARCHAR(50) NOT NULL,
    category VARCHAR(50),
    notes TEXT,
    related_diary_id INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_bills_user_id ON bills(user_id);
CREATE INDEX IF NOT EXISTS idx_bills_created_at ON bills(created_at);
CREATE INDEX IF NOT EXISTS idx_bills_category ON bills(category);
CREATE INDEX IF NOT EXISTS idx_bills_payment_method ON bills(payment_method);

-- 创建账单文件关联表
CREATE TABLE IF NOT EXISTS bill_files (
    bill_id INTEGER NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
    file_id INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    PRIMARY KEY (bill_id, file_id)
);

-- 为账单表添加更新时间触发器
DROP TRIGGER IF EXISTS update_bills_updated_at ON bills;
CREATE TRIGGER update_bills_updated_at
    BEFORE UPDATE ON bills
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加外键约束（可选，如果需要强制关联）
-- ALTER TABLE diaries ADD CONSTRAINT fk_diaries_related_bill
--     FOREIGN KEY (related_bill_id) REFERENCES bills(id) ON DELETE SET NULL;
-- ALTER TABLE bills ADD CONSTRAINT fk_bills_related_diary
--     FOREIGN KEY (related_diary_id) REFERENCES diaries(id) ON DELETE SET NULL;
