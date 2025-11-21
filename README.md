# ecommerce

电商平台后端服务 - 支持文件和文件夹上传（包括空文件夹）

## 核心功能

✅ **支持空文件夹上传** - 完整保留目录结构  
✅ ZIP 文件自动解压  
✅ 自动创建 .gitkeep 确保 Git 跟踪空目录  
✅ REST API 接口  

## 快速开始

### 1. 启动服务器

```bash
cd server
go run main.go
```

服务器将在 `http://localhost:8080` 启动

### 2. 上传文件夹（支持空文件夹）

```bash
# 打包文件夹（保留空目录）
zip -r my-folder.zip my-folder/

# 上传到服务器
curl -X POST http://localhost:8080/api/upload/folder \
  -F "file=@my-folder.zip"
```

### 3. API 端点

- `POST /api/upload/folder` - 上传文件夹（ZIP格式，支持空文件夹）
- `POST /api/upload/file` - 上传单个文件
- `GET /health` - 健康检查

## 详细文档

查看 [UPLOAD_GUIDE.md](./UPLOAD_GUIDE.md) 获取完整的使用说明和 API 文档。

## 为什么能上传空文件夹？

Git 默认不跟踪空目录，但本项目通过以下方式解决：

1. **ZIP 格式支持**: ZIP 文件原生支持空目录信息
2. **自动解压**: 解压时完整还原目录结构
3. **.gitkeep 文件**: 自动在空目录中创建 `.gitkeep` 文件
4. **Git 跟踪**: 有了 `.gitkeep` 文件，Git 就能跟踪空目录

## 项目结构

```
ecommerce/
└── server/              # 后端服务模块（项目根目录）
    ├── service/         # 应用服务（包含上传逻辑）
    │   └── upload_service.go  # 空文件夹处理核心逻辑
    ├── router/          # 路由定义
    │   └── upload_router.go   # 上传 API 路由
    ├── uploads/         # 上传文件存储目录
    │   └── .gitkeep     # 确保目录被 Git 跟踪
    ├── config.yaml      # 配置文件
    ├── go.mod           # 模块依赖管理
    └── main.go          # 启动 HTTP 服务
```

## 测试

运行测试脚本验证空文件夹上传功能：

```bash
./test_upload.sh
```

## 技术栈

- **语言**: Go 1.21+
- **框架**: Gin Web Framework
- **功能**: 文件上传、ZIP 解压、目录管理