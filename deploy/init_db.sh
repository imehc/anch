#!/bin/bash

# PostgreSQL 数据库部署脚本
# 使用 Docker 容器执行数据库初始化和迁移

set -e

# 配置（与 server/config.yaml 保持一致）
POSTGRES_CONTAINER=${POSTGRES_CONTAINER:-"anch-postgres"}
DB_USER=${DB_USER:-"admin"}
DB_NAME=${DB_NAME:-"anch"}

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}==================================${NC}"
echo -e "${GREEN}PostgreSQL 数据库部署${NC}"
echo -e "${GREEN}==================================${NC}"
echo ""

# 检查 Docker 容器是否运行
echo -e "${YELLOW}1. 检查 PostgreSQL 容器...${NC}"
if ! docker ps --format '{{.Names}}' | grep -q "^${POSTGRES_CONTAINER}$"; then
    echo -e "${RED}错误: PostgreSQL 容器 '${POSTGRES_CONTAINER}' 未运行${NC}"
    echo "请先启动 PostgreSQL 容器，或设置正确的容器名称："
    echo "  export POSTGRES_CONTAINER=你的容器名称"
    echo ""
    echo "当前运行的容器："
    docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}"
    exit 1
fi
echo -e "${GREEN}✓ 容器运行正常${NC}"
echo ""

# 创建数据库
echo -e "${YELLOW}2. 创建数据库 '${DB_NAME}'...${NC}"
if docker exec -i ${POSTGRES_CONTAINER} psql -U ${DB_USER} -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '${DB_NAME}'" | grep -q 1; then
    echo -e "${GREEN}✓ 数据库已存在${NC}"
else
    docker exec -i ${POSTGRES_CONTAINER} psql -U ${DB_USER} -d postgres -c "CREATE DATABASE ${DB_NAME};"
    echo -e "${GREEN}✓ 数据库创建成功${NC}"
fi
echo ""

# 运行迁移
echo -e "${YELLOW}3. 运行数据库迁移...${NC}"
MIGRATION_DIR="$(dirname "$0")/migrations"

if [ ! -d "$MIGRATION_DIR" ]; then
    echo -e "${RED}错误: 迁移目录不存在: $MIGRATION_DIR${NC}"
    exit 1
fi

MIGRATION_COUNT=0
for migration in "$MIGRATION_DIR"/*.sql; do
    if [ -f "$migration" ]; then
        filename=$(basename "$migration")
        echo -e "  执行: ${filename}"
        docker exec -i ${POSTGRES_CONTAINER} psql -U ${DB_USER} -d ${DB_NAME} < "$migration"
        MIGRATION_COUNT=$((MIGRATION_COUNT + 1))
    fi
done

if [ $MIGRATION_COUNT -eq 0 ]; then
    echo -e "${YELLOW}⚠ 没有找到迁移文件${NC}"
else
    echo -e "${GREEN}✓ 执行了 ${MIGRATION_COUNT} 个迁移文件${NC}"
fi
echo ""

# 完成
echo -e "${GREEN}==================================${NC}"
echo -e "${GREEN}部署完成！${NC}"
echo -e "${GREEN}==================================${NC}"
echo ""
echo "数据库信息："
echo "  数据库名: ${DB_NAME}"
echo "  用户: ${DB_USER}"
echo ""
echo "测试用户已创建："
echo "  用户名: testuser, 密码: password123"
echo "  用户名: admin, 密码: admin123"
echo ""
