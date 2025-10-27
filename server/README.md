# Diary & Bill API Server

基于 OpenAPI 规范自动生成的 Go 服务端应用，使用 PostgreSQL 数据库和 JWT 认证。

## 功能特性

- ✅ 基于 OpenAPI 3.1 自动生成 API 代码
- ✅ PostgreSQL 数据库（原生 SQL，不使用 ORM）
- ✅ Docker Compose 支持
- ✅ JWT 认证和授权
- ✅ Bcrypt 密码加密
- ✅ YAML 配置文件或环境变量配置
- ✅ 用户登录和认证
- 🚧 日记管理（待实现）
- 🚧 账单管理（待实现）
- 🚧 统计分析（待实现）

## 项目结构

```
server/
├── cmd/
│   └── main.go              # 应用入口点
├── config/
│   └── config.go            # 配置加载器
├── db/
│   └── postgres.go          # PostgreSQL 连接管理
├── repository/
│   └── user.go              # 用户数据访问层（原生 SQL）
├── service/
│   ├── auth.go              # 认证服务（已实现登录）
│   ├── diary.go             # 日记服务
│   ├── bill.go              # 账单服务
│   └── stats.go             # 统计服务
├── util/
│   ├── jwt.go               # JWT 工具
│   └── password.go          # 密码加密工具
├── generated/               # OpenAPI 生成的代码
├── migrations/              # 数据库迁移文件
│   └── 001_create_users.sql
├── docker-compose.yml       # Docker Compose 配置
├── config.yaml              # 配置文件
├── Makefile                 # Make 命令
└── go.mod                   # Go 模块配置
```

## 快速开始

### 前置要求

- Go 1.24+
- Docker 和 Docker Compose
- Make

### 1. 启动数据库

```bash
make db-up
```

这会启动一个 PostgreSQL 15 Docker 容器。

### 2. 初始化数据库

```bash
make db-init
```

这会：
- 创建 `anch` 数据库
- 运行所有迁移脚本
- 创建测试用户

### 3. 编译并运行服务器

```bash
# 编译
make build

# 运行
make run
```

服务器将在 `http://localhost:6020` 启动。

## Make 命令

```bash
make help           # 显示所有可用命令
make gen-apis       # 根据 OpenAPI 规范生成代码
make build          # 编译项目
make run            # 运行服务器
make db-up          # 启动 PostgreSQL Docker 容器
make db-down        # 停止 PostgreSQL 容器
make db-init        # 初始化数据库
make db-migrate     # 运行数据库迁移
make db-reset       # 重置数据库（删除并重新创建）
make clean          # 清理编译文件
```

## 配置

### 配置文件 (config.yaml)

```yaml
server:
  host: localhost
  port: 6020

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: anch
  sslmode: disable

jwt:
  secret_key: your-super-secret-jwt-key-change-this-in-production
  access_token_duration: 2h
  refresh_token_duration: 168h  # 7 days
```

### 环境变量

也可以使用环境变量配置：

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=anch
export JWT_SECRET=your-secret-key

./bin/server -env
```

## API 使用

### 测试用户

数据库初始化后会自动创建测试用户：

| 用户名 | 密码 | 角色 |
|--------|------|------|
| testuser | password123 | user |
| admin | admin123 | admin |

### 用户登录

```bash
curl -X POST http://localhost:6020/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

响应：
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 7200
}
```

### 获取当前用户信息

```bash
curl -X GET http://localhost:6020/api/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## API 端点

### 认证 (Auth)
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/me` - 获取当前用户信息

### 日记 (Diary)
- `GET /api/diary` - 获取日记列表
- `POST /api/diary` - 创建日记
- `GET /api/diary/{id}` - 查询单条日记
- `PUT /api/diary/{id}` - 更新日记
- `DELETE /api/diary/{id}` - ���除日记

### 账单 (Bill)
- `GET /api/bill` - 获取账单列表
- `POST /api/bill` - 创建账单
- `GET /api/bill/{id}` - 查询单条账单
- `PUT /api/bill/{id}` - 更新账单
- `DELETE /api/bill/{id}` - 删除账单

### 统计 (Stats)
- `GET /api/stats/monthly?month=2025-10` - 获取月度统计
- `GET /api/stats/category?month=2025-10` - 获取分类支出占比
- `GET /api/stats/discount?month=2025-10` - 获取优惠类型占比
- `GET /api/stats/trend?month=2025-10` - 获取趋势数据

## 开发指南

### 数据库操作

所有数据库操作都使用原生 SQL，位于 `repository/` 目录：

```go
// 示例：查询用户
user, err := userRepo.GetByUsername(ctx, "testuser")
if err != nil {
    return err
}
```

### 添加新的数据库迁移

1. 在 `migrations/` 目录创建新的 SQL 文件：
   ```
   002_create_diaries.sql
   003_create_bills.sql
   ```

2. 运行迁移：
   ```bash
   make db-migrate
   ```

### 实现服务逻辑

在 `service/` 目录下实现业务逻辑：

```go
func (s *DiaryService) CreateDiary(ctx context.Context, req api.DiaryCreate) (api.ImplResponse, error) {
    // 1. 验证输入
    // 2. 调用 repository
    // 3. 返回响应
    return api.Response(http.StatusCreated, diary), nil
}
```

### 重新生成 API 代码

修改 `openapi.yaml` 后：

```bash
make gen-apis
```

## Docker

### Docker Compose 配置

```yaml
services:
  postgres:
    image: postgres:15-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data
```

### 常用命令

```bash
# 启动服务
docker compose up -d postgres

# 查看日志
docker compose logs -f postgres

# 停止服务
docker compose down

# 进入 PostgreSQL 容器
docker compose exec postgres psql -U postgres
```

## 依赖

- **github.com/go-chi/chi/v5** - HTTP 路由器
- **github.com/lib/pq** - PostgreSQL 驱动
- **github.com/golang-jwt/jwt/v5** - JWT 实现
- **golang.org/x/crypto** - 密码加密
- **gopkg.in/yaml.v3** - YAML 配置解析

## 安全建议

1. **生产环境**：
   - 修改 `jwt.secret_key` 为强密码
   - 使用环境变量而非配置文件存储敏感信息
   - 启用 SSL/TLS (`sslmode: require`)
   - 定期轮换 JWT 密钥

2. **密码策略**：
   - 使用 bcrypt 加密（已实现）
   - 建议设置密码复杂度要求
   - 考虑添加登录失败次数限制

3. **数据库**：
   - 使用专用数据库用户，不要使用 postgres 超级用户
   - 限制数据库用户权限
   - 定期备份数据

## 故障排除

### 数据库连接失败

```bash
# 检查容器是否运行
docker compose ps

# 查看容器日志
docker compose logs postgres

# 重启容器
docker compose restart postgres
```

### 端口冲突

如果 5432 端口已被占用，修改 `docker-compose.yml`：

```yaml
ports:
  - "5433:5432"  # 使用 5433 端口
```

然后更新 `config.yaml` 中的 `database.port`。

## License

MIT
