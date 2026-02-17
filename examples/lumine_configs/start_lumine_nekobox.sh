#!/system/bin/sh
# Lumine + NekoBox 一键启动脚本
# One-Click Start Script for Lumine + NekoBox

# 配置文件路径 / Config file path
CONFIG_PATH="/sdcard/lumine_config.json"

# 如果配置文件不存在，创建默认配置 / Create default config if not exists
if [ ! -f "$CONFIG_PATH" ]; then
    cat > "$CONFIG_PATH" << 'LUMINECONFIG'
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "https://1.1.1.1/dns-query",
  "dns_cache_ttl": 3600,
  "default_policy": {
    "mode": "proxy",
    "connect_timeout": "10s",
    "fragment_size": 1024,
    "fragment_sleep": "10ms",
    "fragment_type": "tls"
  },
  "domain_policies": {},
  "ip_policies": {}
}
LUMINECONFIG
    echo "✓ 已创建默认配置文件 / Default config created: $CONFIG_PATH"
fi

# 检查 Lumine 是否已在运行 / Check if Lumine is running
if pgrep -f "liblumine.so" > /dev/null; then
    echo "✓ Lumine 已在运行 / Lumine is already running"
else
    echo "→ 启动 Lumine / Starting Lumine..."
    /data/data/io.nekohasekai.sagernet/lib/liblumine.so -c "$CONFIG_PATH" &
    sleep 2
    
    if pgrep -f "liblumine.so" > /dev/null; then
        echo "✓ Lumine 启动成功 / Lumine started successfully"
    else
        echo "✗ Lumine 启动失败 / Lumine failed to start"
        exit 1
    fi
fi

# 启动 NekoBox / Start NekoBox
echo "→ 启动 NekoBox / Starting NekoBox..."
am start -n io.nekohasekai.sagernet/.ui.MainActivity

echo ""
echo "═══════════════════════════════════════"
echo "  ✓ Lumine + NekoBox 已启动"
echo "  ✓ Lumine + NekoBox Started"
echo "═══════════════════════════════════════"
echo ""
echo "在 NekoBox 中导入这个配置:"
echo "Import this config in NekoBox:"
echo ""
echo "  socks5://127.0.0.1:1080#Lumine"
echo ""
echo "═══════════════════════════════════════"
