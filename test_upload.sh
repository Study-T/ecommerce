#!/bin/bash

# 测试脚本：验证空文件夹上传功能

echo "======================================"
echo "测试空文件夹上传功能"
echo "======================================"

# 1. 创建测试目录结构（包含空文件夹）
echo -e "\n[1] 创建测试目录结构..."
mkdir -p test-folder/{src,docs,logs,config}
echo "package main" > test-folder/src/main.go
echo "# Test Project" > test-folder/README.md

echo "创建的目录结构:"
tree test-folder/ || find test-folder/ -print

# 2. 打包成 ZIP（保留空目录）
echo -e "\n[2] 打包成 ZIP 文件..."
cd test-folder
zip -r ../test-folder.zip . 
cd ..

echo "ZIP 文件内容:"
unzip -l test-folder.zip

# 3. 启动服务器（后台运行）
echo -e "\n[3] 启动服务器..."
cd server
go run main.go &
SERVER_PID=$!
cd ..

# 等待服务器启动
sleep 3

# 4. 测试健康检查
echo -e "\n[4] 健康检查..."
curl -s http://localhost:8080/health | jq .

# 5. 上传文件夹
echo -e "\n[5] 上传包含空文件夹的 ZIP..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/upload/folder \
  -F "file=@test-folder.zip")

echo "上传响应:"
echo "$RESPONSE" | jq .

# 6. 验证上传结果
echo -e "\n[6] 验证上传结果..."
if [ -d "server/uploads/test-folder" ]; then
    echo "✓ 文件夹上传成功"
    echo -e "\n上传后的目录结构:"
    tree server/uploads/test-folder/ || find server/uploads/test-folder/ -print
    
    # 检查空文件夹是否包含 .gitkeep
    echo -e "\n检查空文件夹中的 .gitkeep 文件:"
    for dir in server/uploads/test-folder/{docs,logs,config}; do
        if [ -f "$dir/.gitkeep" ]; then
            echo "✓ $dir/.gitkeep 存在"
        else
            echo "✗ $dir/.gitkeep 不存在"
        fi
    done
else
    echo "✗ 文件夹上传失败"
fi

# 7. 清理
echo -e "\n[7] 清理测试文件..."
kill $SERVER_PID 2>/dev/null
rm -rf test-folder test-folder.zip

echo -e "\n======================================"
echo "测试完成"
echo "======================================"
