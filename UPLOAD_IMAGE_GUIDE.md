# 图片上传与 `photo_url` 对接说明（方案 A：本地存储）

本项目的 `luggages.photo_url` 字段用于保存图片可访问地址（URL 或相对路径），**图片本体不存数据库**。

实现方式：

- 前端先调用上传接口把图片发给后端
- 后端把图片保存到本地 `./uploads/` 目录
- 后端返回图片 URL
- 前端创建/更新寄存单时，把该 URL 写入 `photo_url` 或 `photo_urls`
- 查询寄存单时，后端返回 `photo_url`/`photo_urls`，前端用 `<img src="...">` 展示

---

## 1. 上传接口说明

### 1.1 POST `/api/upload`（需要登录）

- **Header**：
  - `Authorization: Bearer <token>`
- **Content-Type**：`multipart/form-data`（由浏览器/请求库自动设置，不要手动设置）
- **表单字段**：
  - `file`：图片文件（支持 `jpg/png/webp`，默认最大 5MB）

#### 成功响应示例

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

> 建议前端优先使用 `relative_url` 存到数据库（更容易换域名/IP）。

---

## 2. 静态访问规则（公开）

后端会将本地目录 `./uploads` 以静态资源方式暴露：

- `GET /uploads/...`

例如：

- `http://10.154.39.253:8080/uploads/2026/01/xxx.jpg`

---

## 3. 前端代码示例

### 3.1 JavaScript / TypeScript 示例（原生 fetch）

```javascript
// 方式 1: 使用原生 fetch API
async function uploadImage(file, token) {
  const formData = new FormData();
  formData.append('file', file); // 字段名必须是 'file'
  
  const response = await fetch('http://localhost:8080/api/upload', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`, // 注意：不要设置 Content-Type！
    },
    body: formData // 直接传 FormData，浏览器会自动设置 Content-Type
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Upload failed');
  }
  
  return await response.json();
}

// 使用示例
const fileInput = document.querySelector('input[type="file"]');
const file = fileInput.files[0];
const token = 'your-jwt-token';

uploadImage(file, token)
  .then(result => {
    console.log('Upload success:', result.relative_url);
    // 使用 result.relative_url 创建寄存单
  })
  .catch(error => {
    console.error('Upload failed:', error);
  });
```

### 3.2 Vue 3 示例

```vue
<template>
  <div>
    <input type="file" @change="handleFileChange" accept="image/jpeg,image/png,image/webp" />
    <button @click="uploadFile" :disabled="!selectedFile">上传</button>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import axios from 'axios'; // 或使用你项目中的请求库

const selectedFile = ref(null);
const token = ref('your-jwt-token'); // 从 store 或 localStorage 获取

const handleFileChange = (event) => {
  selectedFile.value = event.target.files[0];
};

const uploadFile = async () => {
  if (!selectedFile.value) return;
  
  const formData = new FormData();
  formData.append('file', selectedFile.value);
  
  try {
    // 使用 axios
    const response = await axios.post('http://localhost:8080/api/upload', formData, {
      headers: {
        'Authorization': `Bearer ${token.value}`,
        // 注意：不要手动设置 Content-Type！
        // axios 会自动为 FormData 设置正确的 Content-Type
      },
    });
    
    console.log('Upload success:', response.data.relative_url);
    return response.data.relative_url;
  } catch (error) {
    console.error('Upload failed:', error.response?.data || error.message);
    throw error;
  }
};
</script>
```

### 3.3 使用 axios 的完整示例（推荐）

```javascript
import axios from 'axios';

// 创建 axios 实例
const api = axios.create({
  baseURL: 'http://localhost:8080',
});

// 请求拦截器：自动添加 token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token'); // 或从 store 获取
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 上传图片函数
export async function uploadImage(file) {
  const formData = new FormData();
  formData.append('file', file);
  
  // 关键：不要设置 Content-Type，让 axios 自动处理
  const response = await api.post('/api/upload', formData);
  return response.data;
}

// 使用
const file = document.querySelector('input[type="file"]').files[0];
uploadImage(file)
  .then(result => {
    console.log('图片地址:', result.relative_url);
  })
  .catch(error => {
    console.error('上传失败:', error);
  });
```

---

## 4. 使用 Apifox 测试上传接口

### 4.1 基础信息

- **Method**：`POST`
- **URL**：`http://10.154.39.253:8080/api/upload`

### 4.2 认证（必须）

在 Apifox 顶部的 **Auth**（或"认证"）里选择：

- **Type**：`Bearer Token`
- **Token**：粘贴登录接口返回的 token（只粘贴 `eyJ...`，不要带 `Bearer`，不要带引号，也不要有前后空格）

> 如果你不用 Auth，也可以在 Headers 里手动加：  
> `Authorization: Bearer <token>`  
> 但不要同时用两种方式，避免出现 `Bearer Bearer ...`。

### 4.3 Body（必须）

在 **Body** 里选择：

- **类型**：`form-data`（也叫 `multipart/form-data`）

然后新增一行字段：

- **Key**：`file`
- **Type**：`File`（一定要选 File，不是 Text）
- **Value**：选择本地图片文件（`.jpg/.png/.webp`）

> 注意：上传接口不是 JSON，不要选 raw/json。

### 4.4 Headers（一般不需要手动写）

- Apifox 会自动带上 `Content-Type: multipart/form-data; boundary=...`
- 你不要手动把 `Content-Type` 写死成 `application/json`

### 4.5 发送并验证

点"发送"后，校验返回体：

- **`message`** 应为 `upload success`
- **`relative_url`** 类似 `/uploads/2026/01/xxx.jpg`
- **`url`** 类似 `http://10.154.39.253:8080/uploads/2026/01/xxx.jpg`

你可以把 `url` 复制到浏览器打开，确认图片能直接访问。

---

## 5. 前端对接流程

### 步骤 1：上传图片

使用 curl 测试（Windows / Git Bash）：

```bash
curl -X POST "http://10.154.39.253:8080/api/upload" \
  -H "Authorization: Bearer <token>" \
  -F "file=@/path/to/photo.jpg"
```

拿到响应中的 `relative_url`（或 `url`）。

### 步骤 2：创建寄存单（复用原接口，不改变结构）

`POST /api/luggage` 的 JSON 中传入图片字段：

- **单图**：传 `photo_url`
- **多图（推荐）**：传 `photo_urls`（数组）

```json
{
  "guest_name": "张三",
  "staff_name": "admin",
  "contact_phone": "13800000000",
  "description": "黑色行李箱",
  "quantity": 1,
  "photo_urls": ["/uploads/2026/01/xxx.jpg", "/uploads/2026/01/yyy.jpg"],
  "photo_url": "/uploads/2026/01/xxx.jpg",
  "storeroom_id": 1
}
```

### 步骤 3：前端展示图片

如果你保存的是相对路径（推荐）：

- `img.src = BaseURL + photo_url`
  - 例如：`BaseURL = "http://10.154.39.253:8080"`
  - `photo_url = "/uploads/2026/01/xxx.jpg"`

如果你保存的是完整 URL：

- `img.src = photo_url`

---

## 6. 常见错误与解决方案

### 6.1 认证相关错误

- **401 invalid token**：确认 `Authorization` 的格式为 `Bearer <token>`，不要多余空格/引号。
- **token is malformed / illegal base64 data at input byte 0**：检查 token 是否带了前导空格/换行；在 Apifox 的 Bearer Token 输入框里只粘贴 `eyJ...`（不要 `Bearer` 前缀）。

### 6.2 上传相关错误

- **request Content-Type isn't multipart/form-data**：说明前端发送请求时 Content-Type 设置错误。
- **上传成功但图片打不开**：确认后端已启动且 `GET /uploads/...` 可以访问；确认文件确实保存到了 `./uploads/...`。
- **file too large**：默认限制 5MB（可在代码里调整）。
- **invalid file type**：仅允许 `jpg/png/webp`。

### 6.3 错误代码示例（避免这些错误）

#### ❌ 错误示例 1: 手动设置 Content-Type

```javascript
// 错误：手动设置 Content-Type 为 application/json
const response = await fetch('/api/upload', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json', // ❌ 错误！
    'Authorization': `Bearer ${token}`,
  },
  body: JSON.stringify({ file: file }) // ❌ 错误！
});
```

#### ❌ 错误示例 2: 使用 JSON 发送文件

```javascript
// 错误：文件不能通过 JSON 发送
const response = await fetch('/api/upload', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json', // ❌ 错误！
  },
  body: JSON.stringify({ 
    file: file, // ❌ 文件对象无法序列化为 JSON
    filename: file.name 
  })
});
```

#### ❌ 错误示例 3: 在 axios 中手动设置 Content-Type

```javascript
// 错误：手动覆盖了 axios 自动设置的 Content-Type
const response = await axios.post('/api/upload', formData, {
  headers: {
    'Content-Type': 'multipart/form-data', // ❌ 错误！缺少 boundary
    'Authorization': `Bearer ${token}`,
  },
});
```

#### ✅ 正确示例

```javascript
// 正确：让浏览器/axios 自动设置 Content-Type
const formData = new FormData();
formData.append('file', file);

const response = await axios.post('/api/upload', formData, {
  headers: {
    // 只设置 Authorization，不设置 Content-Type
    'Authorization': `Bearer ${token}`,
  },
});
```

---

## 7. 关键点总结

1. ✅ **使用 FormData**：必须使用 `FormData` 对象来包装文件
2. ✅ **字段名必须是 'file'**：`formData.append('file', file)`
3. ✅ **不要手动设置 Content-Type**：让浏览器或请求库自动设置
4. ✅ **不要使用 JSON.stringify**：文件不能通过 JSON 发送
5. ✅ **确保 Authorization 头正确**：格式为 `Bearer <token>`

---

## 8. 调试技巧

如果仍然遇到问题，可以在浏览器开发者工具的 Network 标签中检查：

1. 请求的 **Content-Type** 应该是：`multipart/form-data; boundary=...`
2. 请求的 **Payload** 应该显示为 FormData，而不是 JSON
3. 检查 **Request Headers** 中是否有正确的 Authorization 头
