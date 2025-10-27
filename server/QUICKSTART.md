# 快速使用指南

## 前置条件

确保已启动 PostgreSQL Docker 容器，配置如下：
- Host: localhost
- Port: 5432
- User: admin
- Password: admin2025

## 🚀 快速启动（3 步）

### 1. 初始化数据库

```bash
cd ../deploy
./init_db.sh
```

输出示例：
```
==================================
PostgreSQL 数据库部署
==================================

1. 检查 PostgreSQL 容器...
✓ 容器运行正常

2. 创建数据库 'anch'...
✓ 数据库创建成功

3. 运行数据库迁移...
  执行: 001_create_users.sql
✓ 执行了 1 个迁移文件

==================================
部署完成！
==================================

测试用户已创建：
  用户名: testuser, 密码: password123
  用户名: admin, 密码: admin123
```

### 2. 编译并运行服务器

```bash
cd ../server
make build
make run
```

服务器将在 `http://localhost:6020` 启动。

## 🧪 测试登录

### 使用 curl 测试

```bash
# 登录
curl -X POST http://localhost:6020/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' \
  | jq '.'

# 保存 token
TOKEN=$(curl -s -X POST http://localhost:6020/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' \
  | jq -r '.access_token')

# 获取用户信息
curl -X GET http://localhost:6020/api/auth/me \
  -H "Authorization: Bearer $TOKEN" \
  | jq '.'
```

### 预期响应

**登录成功：**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 7200
}
```

**获取用户信息：**
```json
{
  "id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "role": "user",
  "status": "active",
  "created_at": "2025-10-27T15:00:00Z",
  "updated_at": "2025-10-27T15:00:00Z"
}
```

## 📝 常用命令

### 服务器（在 server 目录）

```bash
make help           # 显示所有命令
make gen-apis       # 重新生成 API 代码
make build          # 编译项目
make run            # 运行服务器
make clean          # 清理编译文件
```

### 数据库（在 deploy 目录）

```bash
./init_db.sh        # 初始化数据库和运行迁移
```

## 🔧 配置

### 数据库容器名称

如果你的 PostgreSQL 容器名称不是 `postgres`：

```bash
export POSTGRES_CONTAINER=你的容器名称
cd ../deploy
./init_db.sh
```

### 应用配置

修改 `server/config.yaml`：

```yaml
server:
  port: 6020

database:
  host: localhost
  port: 5432
  user: admin
  password: admin2025
  dbname: anch

jwt:
  secret_key: your-secret-key  # 生产环境必须修改
  access_token_duration: 2h
```

## 🐛 故障排除

### 数据库连接失败

```bash
# 检查 PostgreSQL 是否运行
docker ps | grep postgres

# 查看容器日志
docker logs 容器名称

# 测试连接
docker exec -it 容器名称 psql -U admin -d postgres
```

### 重置数据库

```bash
# 删除数据库（在容器内执行）
docker exec -it postgres psql -U admin -c "DROP DATABASE anch;"

# 重新初始化
cd ../deploy
./init_db.sh
```

## 📁 项目结构

```
anch/
├── deploy/              # 数据库部署文件
│   ├── init_db.sh       # 初始化脚本
│   ├── migrations/      # SQL 迁移文件
│   └── README.md
└── server/              # 应用服务器
    ├── cmd/
    ├── config/
    ├── service/
    ├── Makefile
    └── config.yaml
```

## 📚 下一步

- 阅读完整文档：`server/README.md`
- 数据库部署文档：`deploy/README.md`
- 查看 API 规范：`server/openapi.yaml`
- 实现业务逻辑：`server/service/` 目录
