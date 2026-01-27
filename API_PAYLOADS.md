# API 请求体与响应体示例（对照 `Frontend_Integration_Guide.md`）

## 基础信息

- **Base URL**：`http://本地ip地址:8080`（示例中的 `本地ip地址` 请替换为实际后端地址）
- **请求体格式**：JSON
- **响应格式**：JSON
- **需要登录的接口**：除 `GET /ping` 与 `POST /api/login` 外，其他 `/api/...` 需要携带：

```
Authorization: Bearer <token>
Content-Type: application/json
```

---

## 1) GET /ping

### 请求体

- 无

### 响应体（成功）

```json
{
  "message": "pong"
}
```

---

## 2) POST /api/login

### 请求体

```json
{
  "username": "admin",
  "password": "123456"
}
```

### 响应体（成功）

```json
{
  "message": "login success",
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin",
    "hotel_id": 1
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 响应体（失败）

```json
{
  "message": "login failed",
  "error": "invalid username or password"
}
```

---

## 2.1) POST /api/upload（上传图片，获取 photo_url）

> 需要登录（`Authorization: Bearer <token>`），`multipart/form-data`。

### 请求体

- 表单字段：
  - `file`：图片文件（jpg/png/webp，默认最大 5MB）

### 响应体（成功）

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

### 响应体（失败）

```json
{
  "message": "upload failed",
  "error": "missing file"
}
```

---

## 3) POST /api/luggage（创建寄存单）

### 请求体

```json
{
  "guest_name": "张三",
  "staff_name": "admin",
  "contact_phone": "13800000000",
  "contact_email": "guest@example.com",
  "description": "黑色行李箱",
  "quantity": 1,
  "special_notes": "易碎",
  "photo_urls": ["/uploads/2026/01/xxx.jpg", "/uploads/2026/01/yyy.jpg"],
  "photo_url": "/uploads/2026/01/xxx.jpg",
  "storeroom_id": 1
}
```

### 响应体（成功）

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

### 响应体（失败）

```json
{
  "message": "create luggage failed",
  "error": "storeroom not found"
}
```

---

## 4) GET /api/luggage/by_code（按取件码查询寄存单）

### 请求体

- 无（Query 参数：`code` 必填）

### 响应体（成功）

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

### 响应体（失败）

```json
{
  "message": "query luggage failed",
  "error": "code is empty"
}
```

---

## 5) POST /api/luggage/{id}/checkout（取件）

> 文档约定：Path 参数 `id` 实际传 **取件码**（如 `Z75BDSRH`）。

### 请求体

- 无

### 响应体（成功）

```json
{
  "message": "checkout success",
  "luggage_id": 1
}
```

### 响应体（失败）

```json
{
  "message": "checkout failed",
  "error": "luggage is not in stored status"
}
```

---

## 6) GET /api/luggage/{id}/checkout（获取当前酒店有行李在存的客人名单）

> 文档约定：Path 参数 `id` 为占位即可（如 `any`）。

### 请求体

- 无

### 响应体（成功）

```json
{
  "message": "get checkout info success",
  "items": ["张三", "李四"]
}
```

### 响应体（失败）

```json
{
  "message": "missing user info"
}
```

---

## 7) GET /api/luggage/list/by_guest_name（查询某客人正在寄存的行李）

### 请求体

- 无（Query 参数：`guest_name` 必填）

### 响应体（成功）

```json
{
  "message": "list luggage success",
  "items": [
    {
      "id": 1,
      "guest_name": "张三",
      "retrieval_code": "Z75BDSRH",
      "status": "stored"
    }
  ]
}
```

### 响应体（失败）

```json
{
  "message": "list luggage failed",
  "error": "guest_name is empty"
}
```

---

## 8) GET /api/luggage/storerooms（获取寄存室列表，含容量信息）

### 请求体

- 无

### 响应体（成功）

```json
{
  "message": "list storerooms success",
  "items": [
    {
      "id": 1,
      "hotel_id": 1,
      "name": "A区-1号",
      "location": "一楼A区",
      "capacity": 50,
      "is_active": true,
      "stored_count": 12,
      "remaining_capacity": 38
    }
  ]
}
```

### 响应体（失败）

```json
{
  "message": "list storerooms failed",
  "error": "hotel_id is missing"
}
```

---

## 9) GET /api/luggage/storerooms/{id}/orders（获取某寄存室下的行李订单列表）

### 请求体

- 无（Path 参数：`id` 必填；Query 参数：`status` 可选，例如 `stored`）

### 响应体（成功）

```json
{
  "message": "list luggage success",
  "items": [
    {
      "id": 1,
      "guest_name": "张三",
      "retrieval_code": "Z75BDSRH",
      "status": "stored"
    }
  ]
}
```

### 响应体（失败）

```json
{
  "message": "invalid storeroom id"
}
```

---

## 10) POST /api/luggage/storerooms（创建寄存室）

### 请求体

```json
{
  "name": "A区-1号",
  "location": "一楼A区",
  "capacity": 50,
  "is_active": true
}
```

### 响应体（成功）

```json
{
  "message": "create storeroom success",
  "item": {
    "id": 1,
    "hotel_id": 1,
    "name": "A区-1号",
    "location": "一楼A区",
    "capacity": 50,
    "is_active": true
  }
}
```

### 响应体（失败）

```json
{
  "message": "create storeroom failed",
  "error": "invalid request"
}
```

---

## 11) PUT /api/luggage/storerooms/{id}（软删除/停用寄存室）

### 请求体

```json
{
  "is_active": false
}
```

### 响应体（成功）

```json
{
  "message": "update storeroom status success"
}
```

### 响应体（失败）

```json
{
  "message": "update storeroom status failed",
  "error": "invalid storeroom id"
}
```

---

## 12) GET /api/luggage/logs/stored（获取寄存记录）

### 请求体

- 无

### 响应体（成功）

```json
{
  "message": "list logs success",
  "items": [
    {
      "id": 1,
      "guest_name": "张三",
      "status": "stored",
      "stored_at": "2026-01-22T10:00:00+08:00"
    }
  ]
}
```

### 响应体（失败）

```json
{
  "message": "list logs failed",
  "error": "hotel_id is missing"
}
```

---

## 13) GET /api/luggage/logs/updated（获取寄存信息修改记录）

### 请求体

- 无

### 响应体（成功）

```json
{
  "message": "list logs success",
  "items": [
    {
      "id": 1,
      "hotel_id": 1,
      "luggage_id": 1,
      "updated_by": "staff_user",
      "old_data": "{\"guest_name\":\"张三\"}",
      "new_data": "{\"guest_name\":\"张三\",\"special_notes\":\"易碎\"}",
      "updated_at": "2026-01-22T11:00:00+08:00"
    }
  ]
}
```

### 响应体（失败）

```json
{
  "message": "list logs failed",
  "error": "hotel_id is missing"
}
```

---

## 14) GET /api/luggage/logs/retrieved（获取取出记录）

### 请求体

- 无

### 响应体（成功）

```json
{
  "message": "list logs success",
  "items": [
    {
      "id": 1,
      "guest_name": "张三",
      "retrieved_by": "staff_user",
      "retrieved_at": "2026-01-22T12:00:00+08:00"
    }
  ]
}
```

### 响应体（失败）

```json
{
  "message": "list logs failed",
  "error": "hotel_id is missing"
}
```

---

## 15) PUT /api/luggage/{id}（修改寄存信息）

### 请求体（字段可选）

```json
{
  "guest_name": "张三",
  "contact_phone": "13800000000",
  "description": "黑色行李箱-加锁",
  "special_notes": "易碎",
  "photo_url": "http://example.com/new.jpg"
}
```

### 响应体（成功）

```json
{
  "message": "update luggage success"
}
```

### 响应体（失败）

```json
{
  "message": "update luggage failed",
  "error": "invalid luggage id"
}
```

