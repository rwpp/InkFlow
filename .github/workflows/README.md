# GitHub Actions 构建说明

## 概述

本项目提供了两个 GitHub Actions workflow 用于编译代码：

1. **build.yml** - 多平台构建（支持 Linux/Windows/macOS，AMD64/ARM64）
2. **build-simple.yml** - 简化版（仅 Linux AMD64，推荐使用）

## 使用方法

### 方法一：使用简化版（推荐）

1. 推送代码到 GitHub
2. 进入 GitHub 仓库的 **Actions** 标签页
3. 选择 **Build (Simple - Linux AMD64 Only)** workflow
4. 点击 **Run workflow** 手动触发（或等待自动触发）
5. 等待构建完成
6. 在构建结果页面下载 **release-linux-amd64** 产物

### 方法二：使用多平台构建

1. 推送代码到 GitHub
2. 进入 **Actions** 标签页
3. 选择 **Build** workflow
4. 等待所有平台构建完成
5. 下载对应平台的构建产物

## 构建产物说明

### 简化版构建产物

- `ink-flow` - Go 后端可执行文件（Linux AMD64）
- `ink-flow-backend-linux-amd64.tar.gz` - 后端压缩包
- `ink-flow-frontend.tar.gz` - 前端静态文件压缩包
- `ink-flow-full.tar.gz` - 完整打包（后端+前端）

### 多平台构建产物

- `backend-{os}-{arch}` - 各平台的后端可执行文件
- `frontend-dist` - 前端静态文件
- `release-linux-amd64` - Linux AMD64 完整发布包

## 部署到服务器

### 1. 下载构建产物

从 GitHub Actions 下载 `ink-flow-full.tar.gz`

### 2. 上传到服务器

```bash
scp release/ink-flow-full.tar.gz user@your-server:/path/to/deploy/
```

### 3. 解压并部署

```bash
# 在服务器上
cd /path/to/deploy
tar -xzf ink-flow-full.tar.gz

# 设置执行权限
chmod +x ink-flow

# 创建配置文件
cp config/config.temp.yaml config/config.yaml
# 编辑 config/config.yaml 配置数据库等

# 运行后端
./ink-flow

# 前端文件部署到 Nginx
# 将 dist 目录内容复制到 Nginx 静态文件目录
```

## 自动触发条件

- 推送到 `main`、`master`、`develop` 分支
- 创建 Pull Request 到上述分支
- 手动触发（workflow_dispatch）

## 注意事项

1. 构建产物会保留 30-90 天，请及时下载
2. 首次构建可能需要较长时间（下载依赖）
3. 后续构建会使用缓存，速度更快
4. 确保 `config/config.yaml` 配置文件存在且正确配置

## 故障排查

如果构建失败：

1. 检查 Go 版本是否匹配（需要 Go 1.24+）
2. 检查 Node.js 版本（需要 Node 20+）
3. 查看 Actions 日志了解具体错误
4. 确保所有依赖文件都已提交到仓库

