# 数据库部署文件

此目录包含数据库初始化和迁移脚本。

## 目录结构

```
deploy/
├── init_db.sh           # 数据库初始化脚本
├── migrations/          # SQL 迁移文件
│   └── 001_create_users.sql
└── README.md            # 本文件
```

## 使用方法

### 初始化数据库

确保 PostgreSQL Docker 容器已启动，然后运行：

```bash
./init_db.sh
```

### 自定义配置

如果容器名称不是 `postgres`，可以设置环境变量：

```bash
export POSTGRES_CONTAINER=你的容器名称
./init_db.sh
```

完整配置示例：

```bash
export POSTGRES_CONTAINER=my-postgres
export DB_USER=admin
export DB_NAME=anch
./init_db.sh
```

## 迁移文件命名规范

迁移文件应按顺序命名：

```
001_create_users.sql
002_create_diaries.sql
003_create_bills.sql
...
```

## 添加新迁移

1. 在 `migrations/` 目录创建新的 SQL 文件
2. 使用递增的编号前缀
3. 运行 `./init_db.sh` 应用迁移

## 默认配置

- PostgreSQL 容器: `postgres`
- 数据库用户: `admin`
- 数据库名称: `anch`

这些配置与 `server/config.yaml` 保持一致。

## 测试用户

初始化后会创建以下测试用户：

| 用户名 | 密码 | 角色 |
|--------|------|------|
| testuser | password123 | user |
| admin | admin123 | admin |
