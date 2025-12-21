@echo off
REM 构建 Linux 二进制文件的批处理脚本
REM 目标架构: linux/amd64

echo 开始构建 Linux 二进制文件...

REM 设置环境变量
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0

REM 构建二进制文件
echo 正在编译...
go build -o booksystem-linux-amd64 -ldflags="-w -s" .

if %ERRORLEVEL% EQU 0 (
    echo 构建成功！二进制文件: booksystem-linux-amd64
    dir booksystem-linux-amd64
) else (
    echo 构建失败！
    exit /b 1
)

