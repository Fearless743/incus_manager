#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# ========================
# 检查是否以 root 运行
# ========================
if [ "$(id -u)" -ne 0 ]; then
    log_error "请使用 root 用户运行此脚本"
    exit 1
fi

# ========================
# 1. 安装 Incus
# ========================
log_info "=== 1/4 安装 Incus ==="

if command -v incus &> /dev/null; then
    log_warn "Incus 已安装: $(incus version)"
    read -p "是否重新安装？(y/N) " confirm
    if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
        log_info "跳过 Incus 安装"
    else
        apt-get remove -y incus incus-client 2>/dev/null || true
    fi
else
    apt-get update -y
    apt install curl -y

    mkdir -p /etc/apt/keyrings/
    curl -fsSL https://pkgs.zabbly.com/key.asc -o /etc/apt/keyrings/zabbly.asc
    sh -c 'cat <<EOF > /etc/apt/sources.list.d/zabbly-incus-stable.sources
Enabled: yes
Types: deb
URIs: https://pkgs.zabbly.com/incus/stable
Suites: $(. /etc/os-release && echo ${VERSION_CODENAME})
Components: main
Architectures: $(dpkg --print-architecture)
Signed-By: /etc/apt/keyrings/zabbly.asc
EOF'

    apt-get update
    apt install btrfs-progs -y
    apt-get install incus incus-client -y
fi

# ========================
# 2. 初始化 Incus
# ========================
log_info "=== 2/4 初始化 Incus ==="

if incus list &> /dev/null 2>&1; then
    log_warn "Incus 已初始化"
else
    log_info "正在初始化 Incus..."
    incus admin init --auto --network-address="[::]" --network-port=8443 --storage-backend=btrfs --storage-create-loop=5
    sleep 5
    log_info "Incus 初始化完成"
fi

# 确保 root 在 incus-admin 组
usermod -aG incus-admin root 2>/dev/null || true
newgrp incus-admin 2>/dev/null || true

# ========================
# 3. 配置 Incus 远程和镜像
# ========================
log_info "=== 3/4 配置远程镜像源 ==="

# 添加简单streams远程源
if ! incus remote list | grep -q "^simplestreams"; then
    log_info "添加 simplestreams 远程源..."
    incus remote add simplestreams https://simplestreams.util.eu.org --protocol simplestreams
    sleep 2
fi

# 根据架构拉取常用镜像
ARCH=$(uname -m)
if [ "$ARCH" == "aarch64" ]; then
    log_info "拉取 ARM64 镜像..."
    incus image copy simplestreams:alpine/3.18/arm64/default local: --project=default || true
    incus image copy simplestreams:debian/bullseye/arm64/default local: --project=default || true
else
    log_info "拉取 AMD64 镜像..."
    incus image copy simplestreams:alpine/3.18/amd64/default local: --project=default || true
    incus image copy simplestreams:debian/bullseye/amd64/default local: --project=default || true
    incus image copy simplestreams:debian/bullseye/amd64/cloud local: --project=default --vm || true
fi
sleep 3

# ========================
# 4. 安全加固
# ========================
log_info "=== 4/4 安全加固 ==="

# 防止容器名称泄露
chmod 400 /proc/sched_debug 2>/dev/null || true
chmod 700 /sys/kernel/slab/ 2>/dev/null || true

# 生成用于面板认证的 TLS 凭证
CERT_DIR="/etc/incus-manager"
mkdir -p "$CERT_DIR"

if [ ! -f "$CERT_DIR/server.pem" ] || [ ! -f "$CERT_DIR/client.pem" ]; then
    log_info "生成 TLS 凭证..."
    
    # 生成自签名证书
    openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
        -keyout "$CERT_DIR/server.key" \
        -out "$CERT_DIR/server.pem" \
        -subj "/CN=incus-manager" \
        -addext "subjectAltName=DNS:*,IP:0.0.0.0" 2>/dev/null
    
    # 从 server cert 提取客户端凭证
    cp "$CERT_DIR/server.pem" "$CERT_DIR/client.pem"
    
    log_info "凭证已保存到 $CERT_DIR/"
fi

# 将凭证添加到 Incus 信任列表
TRUST_TOKEN=$(cat "$CERT_DIR/server.pem" "$CERT_DIR/server.key" 2>/dev/null | base64 -w0)
INCUS_CERT_NAME="incus-manager-cert"

if ! incus config trust list | grep -q "$INCUS_CERT_NAME"; then
    incus config trust add "$INCUS_CERT_NAME" <<< "$TRUST_TOKEN" 2>/dev/null || true
fi

# ========================
# 完成
# ========================
echo ""
log_info "=========================================="
log_info "   Incus 安装配置完成"
log_info "=========================================="
echo ""
log_info "Incus 地址: https://$(hostname -I | awk '{print $1}'):8443"
log_info "凭证文件: $CERT_DIR/"
echo ""
log_info "下一步：部署本面板"
log_info "  docker compose up --build -d"
echo ""
log_info "在面板中添加主机时，使用以下信息："
log_info "  地址: https://$(hostname -I | awk '{print $1}'):8443"
log_info "  凭证: $(cat $CERT_DIR/client.pem 2>/dev/null | head -3)..."
echo ""
