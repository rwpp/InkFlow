# Docker 镜像认证指南

## 问题说明

GitHub Container Registry (ghcr.io) 默认是**私有**的，需要认证才能拉取镜像。如果直接拉取会报错：
```
Error response from daemon: denied: permission_denied
```

## 解决方案

### 方法一：使用 GitHub Personal Access Token (PAT) - 推荐

#### 1. 创建 Personal Access Token

1. 访问：https://github.com/settings/tokens
2. 点击 **Generate new token** → **Generate new token (classic)**
3. 设置名称，例如：`docker-pull-token`
4. 选择过期时间
5. **重要**：勾选权限 `read:packages`（读取包）
6. 点击 **Generate token**
7. **复制 token**（只显示一次，请保存好）

#### 2. 在服务器上登录

```bash
# 方式一：交互式登录
docker login ghcr.io -u YOUR_GITHUB_USERNAME
# 输入密码时，粘贴刚才创建的 PAT token

# 方式二：使用环境变量（推荐，适合脚本）
export GITHUB_TOKEN=your_pat_token_here
echo $GITHUB_TOKEN | docker login ghcr.io -u YOUR_GITHUB_USERNAME --password-stdin

# 方式三：一行命令
echo "your_pat_token" | docker login ghcr.io -u YOUR_GITHUB_USERNAME --password-stdin
```

#### 3. 验证登录

```bash
# 查看登录信息
cat ~/.docker/config.json

# 测试拉取镜像
docker pull ghcr.io/YOUR_USERNAME/InkFlow-backend:latest
```

#### 4. 持久化登录（可选）

登录信息会保存在 `~/.docker/config.json`，下次不需要重新登录。

### 方法二：将仓库设置为公开

如果镜像不需要保密，可以将仓库设置为公开：

1. 进入 GitHub 仓库
2. Settings → General → Danger Zone
3. 点击 **Change visibility** → **Make public**
4. 这样任何人都可以拉取镜像，无需认证

### 方法三：使用 GitHub Actions 自动登录（CI/CD）

在 GitHub Actions 中，使用 `GITHUB_TOKEN` 自动登录：

```yaml
- name: Log in to Container Registry
  uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}
```

## 使用 Digest 拉取镜像

### 获取 Digest

1. 在 GitHub Actions 构建完成后，查看 **Summary**
2. 复制显示的 Digest，例如：
   ```
   ghcr.io/username/InkFlow-backend@sha256:abc123def456...
   ```

### 使用 Digest 拉取

```bash
# 先登录（必须）
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# 使用 Digest 拉取（最精确）
docker pull ghcr.io/username/InkFlow-backend@sha256:abc123def456...

# 或使用标签拉取
docker pull ghcr.io/username/InkFlow-backend:latest
```

## 在 docker-compose.yaml 中使用

### 方式一：先登录，再使用

```bash
# 1. 登录
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# 2. 使用 docker-compose
docker-compose pull
docker-compose up -d
```

### 方式二：在 docker-compose.yaml 中配置认证

```yaml
services:
  backend:
    image: ghcr.io/username/InkFlow-backend:latest
    # 或者使用 Digest
    # image: ghcr.io/username/InkFlow-backend@sha256:abc123...
```

然后在服务器上先登录 Docker。

## 常见问题

### Q: 为什么需要认证？

A: GitHub Container Registry 默认是私有的，即使是公开仓库的镜像也需要认证才能拉取。

### Q: PAT Token 过期了怎么办？

A: 重新生成一个新的 PAT，然后重新登录即可。

### Q: 可以在多个服务器上使用同一个 Token 吗？

A: 可以，但建议为不同环境创建不同的 Token，便于管理和撤销。

### Q: 如何查看镜像的 Digest？

```bash
# 拉取镜像后查看
docker images --digests

# 或使用 docker inspect
docker inspect ghcr.io/username/InkFlow-backend:latest | grep -i digest
```

## 安全建议

1. ✅ 使用 PAT Token 而不是密码
2. ✅ Token 只授予最小必要权限（`read:packages`）
3. ✅ 定期轮换 Token
4. ✅ 不要在代码中硬编码 Token
5. ✅ 使用环境变量存储 Token

