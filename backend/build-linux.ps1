# 构建 Linux 二进制文件的 PowerShell 脚本
# 目标架构: linux/amd64

Write-Host "开始构建 Linux 二进制文件..." -ForegroundColor Green

# 设置环境变量
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"

# 构建二进制文件
Write-Host "正在编译..." -ForegroundColor Yellow
go build -o booksystem-linux-amd64 -ldflags="-w -s" .

if ($LASTEXITCODE -eq 0) {
    Write-Host "构建成功！二进制文件: booksystem-linux-amd64" -ForegroundColor Green
    Write-Host "文件大小:" -ForegroundColor Cyan
    Get-Item booksystem-linux-amd64 | Select-Object Name, @{Name="Size(MB)";Expression={[math]::Round($_.Length/1MB, 2)}}
} else {
    Write-Host "构建失败！" -ForegroundColor Red
    exit 1
}


