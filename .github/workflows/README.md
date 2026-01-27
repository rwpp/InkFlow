# GitHub Actions 构建说明

## 概述

本项目提供了两个 GitHub Actions workflow：

1. **build-simple.yml** - Linux AMD64 构建（直接运行二进制文件）⭐
2. **build-docker.yml** - Docker 镜像构建（使用 Docker 部署）

## 选择哪个？

- **只用 Linux 服务器，直接运行** → 使用 `build-simple.yml`
- **使用 Docker 部署** → 使用 `build-docker.yml`（需要配置认证，见下方说明）

## 使用方法

1. 推送代码到 GitHub
2. 进入 GitHub 仓库的 **Actions** 标签页
3. 选择 **Build (Simple - Linux AMD64 Only)** workflow
4. 点击 **Run workflow** 手动触发（或等待自动触发）
5. 等待构建完成
6. 在构建结果页面下载 **release-linux-amd64** 产物

## 构建产物说明

构建完成后会生成以下文件：

- `ink-flow` - Go 后端可执行文件（Linux AMD64）
- `ink-flow-backend-linux-amd64.tar.gz` - 后端压缩包
- `ink-flow-frontend.tar.gz` - 前端静态文件压缩包
- `ink-flow-full.tar.gz` - 完整打包（后端+前端，推荐下载这个）

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

## Docker 镜像认证问题

如果使用 `build-docker.yml` 构建的镜像，在服务器上拉取时遇到认证错误：

```
Error response from daemon: denied: permission_denied
```

### 解决方法

1. **创建 GitHub Personal Access Token**
   - 访问：https://github.com/settings/tokens
   - 创建新 Token，勾选 `read:packages` 权限

2. **在服务器上登录**
   ```bash
   echo "your_pat_token" | docker login ghcr.io -u YOUR_USERNAME --password-stdin
   ```

3. **拉取镜像**
   ```bash
   docker pull ghcr.io/your-username/InkFlow-backend:latest
   ```

详细说明请查看：[docs/docker-auth.md](../docs/docker-auth.md)

## 故障排查

如果构建失败：

1. 检查 Go 版本是否匹配（需要 Go 1.24+）
2. 检查 Node.js 版本（需要 Node 20+）
3. 查看 Actions 日志了解具体错误
4. 确保所有依赖文件都已提交到仓库
5. Docker 构建需要仓库有 `packages: write` 权限（通常自动配置）
5. Docker 构建需要仓库有 `packages: write` 权限（通常自动配置）

