# bookmananer

图书管理系统

## 项目简介

这是一个基于 Vue 3 + Go 的图书管理系统，支持图书管理、借阅、归还等功能。

## 技术栈

- 前端：Vue 3 + Vite + Element Plus
- 后端：Go + Hertz + GORM + SQLite

## 功能特性

- 图书管理
- 位置管理（区域、书架、层架）
- 图书借阅
- 图书归还
- 借阅记录查询

## 快速开始

### 后端启动

```bash
cd backend
go run main.go
```

后端服务运行在 `http://localhost:8089`

### 前端启动

```bash
cd frontend
npm install
npm run dev
```

前端服务运行在 `http://localhost:3000`

## 项目结构

```
booksystem/
├── backend/          # 后端代码
│   ├── internal/     # 内部包
│   │   ├── db/       # 数据库模型
│   │   ├── handler/  # 请求处理器
│   │   └── middleware/ # 中间件
│   └── main.go       # 入口文件
├── frontend/         # 前端代码
│   ├── src/
│   │   ├── api/      # API接口
│   │   ├── components/ # 组件
│   │   ├── views/    # 页面
│   │   └── router/   # 路由
│   └── package.json
└── database/         # 数据库初始化脚本
```

