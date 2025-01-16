#!/bin/bash

# 设置 Go 环境变量
GO111MODULE=on
GOOS=windows   # 修改为 Windows 系统
GOARCH=amd64   # 64 位架构

# 项目根目录
PROJECT_DIR=$(pwd)

# 清理之前的编译产物
echo "清理之前的编译文件..."
rm -f "$PROJECT_DIR/memory-monitor.exe"  # Windows 下使用 .exe 扩展名

# 安装依赖
echo "安装项目依赖..."
go mod tidy

# 编译项目
echo "编译项目..."
go build -o memory-monitor.exe "$PROJECT_DIR/cmd"

# 输出构建结果
if [ -f "$PROJECT_DIR/memory-monitor.exe" ]; then
  echo "构建成功！"
else
  echo "构建失败！"
  exit 1
fi
