# 文件夹上传功能说明（支持空文件夹）

## 问题背景

Git 默认不跟踪空文件夹，这导致很多项目在上传文件夹时会丢失空目录结构。本项目实现了完整的文件夹上传功能，**完全支持空文件夹的上传和保存**。

## 解决方案

### 1. 核心功能

- ✅ 支持 ZIP 文件上传并自动解压
- ✅ **完整保留空文件夹结构**
- ✅ 自动在空文件夹中创建 `.gitkeep` 文件，确保 Git 可以跟踪
- ✅ 支持单个文件上传
- ✅ 支持嵌套目录结构

### 2. API 端点

#### 上传文件夹（支持空文件夹）
```
POST /api/upload/folder
Content-Type: multipart/form-data
```

**请求参数:**
- `file`: ZIP 格式的压缩文件

**响应示例:**
```json
{
  "code": 200,
  "message": "文件夹上传成功（包括空文件夹）",
  "data": {
    "path": "./uploads/my-folder",
    "filename": "my-folder.zip",
    "size": 1024
  }
}
```

#### 上传单个文件
```
POST /api/upload/file
Content-Type: multipart/form-data
```

**请求参数:**
- `file`: 任意格式的文件

### 3. 使用示例

#### 使用 curl 上传文件夹
```bash
# 1. 将文件夹打包成 ZIP（保留空目录）
zip -r my-folder.zip my-folder/

# 2. 上传 ZIP 文件
curl -X POST http://localhost:8080/api/upload/folder \
  -F "file=@my-folder.zip"
```

#### 使用 curl 上传文件
```bash
curl -X POST http://localhost:8080/api/upload/file \
  -F "file=@document.pdf"
```

#### 前端 JavaScript 示例
```javascript
// 上传文件夹（ZIP格式）
async function uploadFolder(zipFile) {
  const formData = new FormData();
  formData.append('file', zipFile);
  
  const response = await fetch('http://localhost:8080/api/upload/folder', {
    method: 'POST',
    body: formData
  });
  
  const result = await response.json();
  console.log('上传结果:', result);
}

// 上传单个文件
async function uploadFile(file) {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await fetch('http://localhost:8080/api/upload/file', {
    method: 'POST',
    body: formData
  });
  
  const result = await response.json();
  console.log('上传结果:', result);
}
```

### 4. 空文件夹处理机制

当上传包含空文件夹的 ZIP 文件时：

1. **检测空目录**: 系统自动识别 ZIP 中的空目录项
2. **创建目录**: 在目标位置创建相应的空目录
3. **添加 .gitkeep**: 自动在空目录中创建 `.gitkeep` 文件
4. **Git 跟踪**: `.gitkeep` 文件确保 Git 可以跟踪空目录

**示例目录结构:**
```
uploads/
└── my-project/
    ├── src/
    │   └── main.go
    ├── docs/          # 空文件夹
    │   └── .gitkeep   # 自动创建
    ├── logs/          # 空文件夹
    │   └── .gitkeep   # 自动创建
    └── README.md
```

### 5. 配置说明

在 `main.go` 中可以配置：

```go
// 修改上传文件大小限制（默认 32MB）
r.MaxMultipartMemory = 32 << 20  // 32MB

// 修改上传目录路径
uploadPath := filepath.Join(".", "uploads")
```

### 6. 运行服务器

```bash
cd server
go run main.go
```

服务器将在 `http://localhost:8080` 启动

### 7. 健康检查

```bash
curl http://localhost:8080/health
```

**响应:**
```json
{
  "status": "healthy",
  "message": "服务器运行正常，支持空文件夹上传"
}
```

## 技术实现

### 关键代码逻辑

```go
// 1. 检测目录项
if file.FileInfo().IsDir() {
    // 创建空目录
    os.MkdirAll(filePath, file.Mode())
    
    // 2. 添加 .gitkeep 确保 Git 跟踪
    gitkeepPath := filepath.Join(filePath, ".gitkeep")
    os.WriteFile(gitkeepPath, []byte(""), 0644)
}
```

### 优势

1. **完整性**: 保留原始文件夹结构，包括空目录
2. **兼容性**: 使用标准 ZIP 格式，兼容所有压缩工具
3. **Git 友好**: 自动添加 `.gitkeep` 确保版本控制
4. **安全性**: 路径验证防止目录遍历攻击
5. **易用性**: 简单的 REST API，支持多种客户端

## 常见问题

### Q: 为什么要用 ZIP 格式？
A: ZIP 格式原生支持目录结构信息，包括空目录，是最通用的压缩格式。

### Q: .gitkeep 文件是什么？
A: `.gitkeep` 是一个约定俗成的空文件，用于让 Git 跟踪空目录（Git 默认不跟踪空目录）。

### Q: 可以上传多大的文件？
A: 默认限制为 32MB，可以在 `main.go` 中修改 `MaxMultipartMemory` 配置。

### Q: 支持哪些压缩格式？
A: 目前仅支持 ZIP 格式，因为它对空目录的支持最好。

## 项目结构

```
server/
├── service/
│   └── upload_service.go    # 上传服务（空文件夹处理逻辑）
├── router/
│   └── upload_router.go     # 路由处理器
├── uploads/                  # 上传文件存储目录
├── main.go                   # 服务器入口
└── go.mod                    # 依赖管理
```
