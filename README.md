# 酒店行李寄存系统后端

基于 Gin + GORM + MySQL 的酒店行李寄存管理系统后端。

## 功能特性

- 用户认证（JWT Token）
- 行李寄存单管理
- 寄存室管理
- 取件功能
- 操作日志记录
- RESTful API 接口

## 技术栈

- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置数据库

设置环境变量 `DB_DSN`：

```bash
# Windows CMD
set "DB_DSN=root:你的密码@tcp(127.0.0.1:3306)/hotel_luggage?charset=utf8mb4&parseTime=True&loc=Local"

# Windows PowerShell
$env:DB_DSN="root:你的密码@tcp(127.0.0.1:3306)/hotel_luggage?charset=utf8mb4&parseTime=True&loc=Local"

# Linux/Mac
export DB_DSN="root:你的密码@tcp(127.0.0.1:3306)/hotel_luggage?charset=utf8mb4&parseTime=True&loc=Local"
```

### 3. 创建数据库

在 MySQL 中创建数据库：

```sql
CREATE DATABASE IF NOT EXISTS hotel_luggage CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 运行项目

```bash
go run .
```

看到 `Listening and serving HTTP on :8080` 说明启动成功。

### 5. 初始化数据

系统会自动创建表结构。默认测试账号：
- 用户名: `admin`
- 密码: `123456`

## API 接口

详细接口文档请参考 `Frontend_Integration_Guide.md`

### 公开接口
- `GET /ping` - 健康检查
- `POST /api/login` - 用户登录

### 需要认证的接口（需要 Authorization Header）
- `POST /api/luggage` - 创建寄存单
- `GET /api/luggage/by_code` - 按取件码查询
- `POST /api/luggage/{id}/checkout` - 取件
- `GET /api/luggage/{id}/checkout` - 获取客人名单
- `GET /api/luggage/list/by_guest_name` - 查询客人行李
- `GET /api/luggage/storerooms` - 获取寄存室列表
- `POST /api/luggage/storerooms` - 创建寄存室
- `PUT /api/luggage/storerooms/{id}` - 更新寄存室
- `GET /api/luggage/storerooms/{id}/orders` - 获取寄存室订单
- `PUT /api/luggage/{id}` - 修改寄存信息
- `GET /api/luggage/logs/stored` - 获取寄存记录
- `GET /api/luggage/logs/updated` - 获取修改记录
- `GET /api/luggage/logs/retrieved` - 获取取出记录

## 项目结构

```
.
├── main.go              # 主程序入口
├── internal/
│   ├── config/          # 配置管理
│   ├── models/          # 数据模型
│   ├── middleware/      # 中间件
│   ├── handlers/        # 请求处理器
│   ├── services/        # 业务逻辑层
│   └── utils/           # 工具函数
├── go.mod
├── go.sum
└── README.md
```

## 开发说明

- 数据库表会在首次运行时自动创建
- JWT Secret 默认使用 "your-secret-key"，生产环境请修改
- 密码使用 bcrypt 加密存储
