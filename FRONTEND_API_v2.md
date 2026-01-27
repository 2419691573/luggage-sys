# 前端 API 对接文档（Gin + GORM + MySQL）

本文档面向前端同学，描述接口的**用途、请求体/响应体字段含义、示例**。以实际后端实现为准。

---

## 0. 基础信息

- **Base URL**：`http://10.154.39.253:8080`
- **数据格式**：
  - JSON 接口：`Content-Type: application/json`
  - 上传接口：`multipart/form-data`
- **认证方式**：
  - 除 `GET /ping`、`POST /api/login` 外，其余 `/api/**` 均需要 Header：

```
Authorization: Bearer <token>
```

- **通用响应字段**：
  - `message`：字符串，表示本次请求结果（成功/失败）
  - `error`：字符串（可选），失败原因

---

## 1. 健康检查

### 1.1 GET `/ping`

**用途**：检查服务是否运行。

**请求**：无

**响应（200）**：

```json
{ "message": "pong" }
```

---

## 2. 登录

### 2.1 POST `/api/login`

**用途**：登录获取 JWT token。

**请求体（JSON）**：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `username` | string | 是 | 登录用户名（默认 `admin`） |
| `password` | string | 是 | 登录密码（默认 `123456`） |

**响应（200）**：

| 字段 | 类型 | 说明 |
|---|---|---|
| `message` | string | 固定：`login success` |
| `user.id` | number | 用户 ID |
| `user.username` | string | 用户名 |
| `user.role` | string | 角色（示例：`admin`） |
| `user.hotel_id` | number | 酒店 ID（后端用来做酒店隔离） |
| `token` | string | JWT token（后续请求放到 `Authorization`） |

示例：

```json
{
  "message": "login success",
  "user": { "id": 1, "username": "admin", "role": "admin", "hotel_id": 1 },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**响应（401）**：

```json
{ "message": "login failed", "error": "invalid username or password" }
```

---

## 3. 图片上传（用于 `photo_urls` / `photo_url`）

> 说明：行李表里存的是图片地址（URL），图片本体不进数据库。  
> 推荐把上传接口返回的 `relative_url` 写入 `photo_urls`（数组），便于“一个寄存单多张图”。

### 3.1 POST `/api/upload`（需要登录）

**用途**：上传图片，返回可访问的图片 URL。

**请求（multipart/form-data）**：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `file` | file | 是 | 图片文件（支持 jpg/png/webp，默认最大 5MB） |

**响应（200）**：

| 字段 | 类型 | 说明 |
|---|---|---|
| `message` | string | 固定：`upload success` |
| `url` | string | 完整 URL（可直接给 `<img src>` 用） |
| `relative_url` | string | 相对路径（推荐存库，如 `/uploads/2026/01/xxx.jpg`） |
| `content_type` | string | 识别到的文件类型（如 `image/jpeg`） |
| `size` | number | 实际写入字节数 |
| `file_name` | string | 服务器生成的文件名 |
| `max_size_byte` | number | 允许的最大字节数（固定 5MB） |

示例：

```json
{
  "message": "upload success",
  "url": "http://10.154.39.253:8080/uploads/2026/01/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.jpg",
  "relative_url": "/uploads/2026/01/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.jpg",
  "content_type": "image/jpeg",
  "size": 123456,
  "file_name": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.jpg",
  "max_size_byte": 5242880
}
```

**失败示例（400）**：

```json
{ "message": "upload failed", "error": "missing file" }
```

### 3.2 图片访问（公开）

上传后的图片可通过静态资源访问：

- `GET /uploads/...`

例如：

- `http://10.154.39.253:8080/uploads/2026/01/xxx.jpg`

---

## 4. 行李寄存单

### 4.1 POST `/api/luggage`（创建寄存单，需要登录）

**用途**：创建寄存单并生成取件码。

**请求体（JSON）**：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `guest_name` | string | 是 | 客人姓名 |
| `staff_name` | string | 否 | 经办人姓名；不传则后端自动用当前登录账号 |
| `contact_phone` | string | 否 | 联系电话 |
| `contact_email` | string | 否 | 联系邮箱 |
| `description` | string | 否 | 行李描述 |
| `quantity` | number | 否 | 数量（默认 1） |
| `special_notes` | string | 否 | 备注 |
| `photo_urls` | string[] | 否 | 图片地址数组（建议用 `/api/upload` 返回的 `relative_url` 组成数组） |
| `photo_url` | string | 否 | 单图兼容字段（如果只传它，后端会自动转成 `photo_urls=[photo_url]`） |
| `storeroom_id` | number | 是 | 寄存室 ID（先通过寄存室接口创建/查询获得） |

**响应（200）**：

| 字段 | 类型 | 说明 |
|---|---|---|
| `message` | string | 固定：`create luggage success` |
| `luggage_id` | number | 寄存单 ID |
| `retrieval_code` | string | 取件码（8 位） |
| `qrcode_url` | string | 二维码地址占位（当前仅返回字符串） |
| `photo_url` | string | 存库后的图片地址（原样返回） |
| `photo_urls` | string[] | 存库后的图片地址数组（推荐使用） |

示例：

```json
{
  "message": "create luggage success",
  "luggage_id": 1,
  "retrieval_code": "Z75BDSRH",
  "qrcode_url": "/qr/Z75BDSRH",
  "photo_url": "/uploads/2026/01/xxx.jpg",
  "photo_urls": ["/uploads/2026/01/xxx.jpg", "/uploads/2026/01/yyy.jpg"]
}
```

**失败示例（400）**：

```json
{ "message": "create luggage failed", "error": "storeroom not found" }
```

### 4.2 GET `/api/luggage/by_code`（按取件码查询，需要登录）

**请求参数（Query）**：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `code` | string | 是 | 取件码（`retrieval_code`） |

**响应（200）**：

| 字段 | 类型 | 说明 |
|---|---|---|
| `message` | string | 固定：`query luggage success` |
| `item.id` | number | 寄存单 ID |
| `item.guest_name` | string | 客人姓名 |
| `item.contact_phone` | string | 联系电话 |
| `item.storeroom_id` | number | 寄存室 ID |
| `item.retrieval_code` | string | 取件码 |
| `item.status` | string | 状态：`stored`（在存）/ `retrieved`（已取） |
| `item.photo_url` | string | 图片地址 |
| `item.photo_urls` | string[] | 图片地址数组 |

示例：

```json
{
  "message": "query luggage success",
  "item": {
    "id": 1,
    "guest_name": "张三",
    "contact_phone": "13800000000",
    "storeroom_id": 1,
    "retrieval_code": "Z75BDSRH",
    "status": "stored",
    "photo_url": "/uploads/2026/01/xxx.jpg",
    "photo_urls": ["/uploads/2026/01/xxx.jpg", "/uploads/2026/01/yyy.jpg"]
  }
}
```

### 4.3 POST `/api/luggage/{code}/checkout`（取件，需要登录）

> Path 参数名在路由里叫 `:id`，实际传取件码即可：`/api/luggage/Z75BDSRH/checkout`

**响应（200）**：

```json
{ "message": "checkout success", "luggage_id": 1 }
```

### 4.4 GET `/api/luggage/{any}/checkout`（获取当前酒店“在存”客人名单，需要登录）

> Path 参数占位即可，例如：`/api/luggage/any/checkout`

**响应（200）**：

```json
{ "message": "get checkout info success", "items": ["张三", "李四"] }
```

### 4.5 GET `/api/luggage/list/by_guest_name`（按客人姓名查在存行李，需要登录）

**Query 参数**：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `guest_name` | string | 是 | 客人姓名 |

**响应（200）**：

```json
{
  "message": "list luggage success",
  "items": [
    {
      "id": 1,
      "guest_name": "张三",
      "retrieval_code": "Z75BDSRH",
      "status": "stored",
      "photo_url": "/uploads/2026/01/xxx.jpg",
      "photo_urls": ["/uploads/2026/01/xxx.jpg", "/uploads/2026/01/yyy.jpg"]
    }
  ]
}
```

### 4.6 PUT `/api/luggage/{id}`（修改寄存信息，需要登录）

**用途**：修改寄存信息（包含 `photo_url`）。

**请求体（JSON，字段可选）**：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `guest_name` | string | 否 | 客人姓名 |
| `contact_phone` | string | 否 | 联系电话 |
| `description` | string | 否 | 行李描述 |
| `special_notes` | string | 否 | 备注 |
| `photo_urls` | string[] | 否 | 图片地址数组（推荐） |
| `photo_url` | string | 否 | 单图兼容字段 |

**响应（200）**：

```json
{ "message": "update luggage success" }
```

---

## 5. 寄存室

### 5.1 GET `/api/luggage/storerooms`（需要登录）

**响应（200）**：返回寄存室列表（后端计算并返回容量信息）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | number | 寄存室 ID |
| `hotel_id` | number | 酒店 ID |
| `name` | string | 寄存室名称 |
| `location` | string | 位置 |
| `capacity` | number | 容量（最大可存数量） |
| `is_active` | boolean | 是否启用 |
| `stored_count` | number | 已存数量（后端计算字段） |
| `remaining_capacity` | number | 剩余容量（后端计算字段） |

### 5.2 POST `/api/luggage/storerooms`（需要登录）

**请求体（JSON）**：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `name` | string | 是 | 名称 |
| `location` | string | 否 | 位置 |
| `capacity` | number | 是 | 容量 |
| `is_active` | boolean | 否 | 是否启用（不传时按后端默认值处理） |

**响应（200）**：

```json
{ "message": "create storeroom success", "item": { "id": 1, "hotel_id": 1, "name": "A区-1号", "location": "一楼A区", "capacity": 50, "is_active": true } }
```

### 5.3 PUT `/api/luggage/storerooms/{id}`（需要登录）

**用途**：启用/停用寄存室（软停用）。

**请求体（JSON）**：

```json
{ "is_active": false }
```

**响应（200）**：

```json
{ "message": "update storeroom status success" }
```

### 5.4 GET `/api/luggage/storerooms/{id}/orders`（需要登录）

**Query 参数**：

| 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|
| `status` | string | 否 | 过滤状态（如 `stored`） |

**响应（200）**：

```json
{
  "message": "list luggage success",
  "items": [
    { "id": 1, "guest_name": "张三", "retrieval_code": "Z75BDSRH", "status": "stored" }
  ]
}
```

---

## 6. 日志

### 6.1 GET `/api/luggage/logs/stored`（需要登录）

**响应（200）**：寄存记录列表（结构由后端模型返回，常用字段如下）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | number | 记录 ID |
| `hotel_id` | number | 酒店 ID |
| `luggage_id` | number | 寄存单 ID |
| `guest_name` | string | 客人姓名 |
| `status` | string | 一般为 `stored` |
| `stored_at` | string | 时间（ISO8601） |

### 6.2 GET `/api/luggage/logs/updated`（需要登录）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | number | 记录 ID |
| `hotel_id` | number | 酒店 ID |
| `luggage_id` | number | 寄存单 ID |
| `updated_by` | string | 修改人 |
| `old_data` | string | 修改前快照（JSON 字符串） |
| `new_data` | string | 修改后快照（JSON 字符串） |
| `updated_at` | string | 修改时间 |

### 6.3 GET `/api/luggage/logs/retrieved`（需要登录）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | number | 记录 ID |
| `hotel_id` | number | 酒店 ID |
| `luggage_id` | number | 寄存单 ID |
| `guest_name` | string | 客人姓名 |
| `retrieved_by` | string | 取件人（登录账号） |
| `retrieved_at` | string | 取件时间 |

---

## 7. 前端最小流程（建议照这个跑通）

1) `POST /api/login` 获取 `token`
2) `POST /api/upload` 上传图片（可多次），收集 `relative_url[]`
3) `POST /api/luggage` 创建寄存单，把 `photo_urls = relative_url[]`（或单图用 `photo_url`）
4) `GET /api/luggage/by_code?code=...` 查询并展示图片（`<img src="BaseURL + photo_url">`）

