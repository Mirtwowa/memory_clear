#!/bin/bash

# 项目目录
PROJECT_DIR=$(pwd)

# 删除二进制文件
echo "清理编译产物..."
rm -f "$PROJECT_DIR/memory-monitor"

echo "清理完成！"
