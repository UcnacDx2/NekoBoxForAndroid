#!/system/bin/sh
# 检查 Lumine 状态
# Check Lumine Status

echo "═══════════════════════════════════════"
echo "  Lumine 状态检查 / Status Check"
echo "═══════════════════════════════════════"
echo ""

# 检查进程 / Check process
if pgrep -f "liblumine.so" > /dev/null; then
    PID=$(pgrep -f "liblumine.so")
    echo "✓ Lumine 正在运行 / Running"
    echo "  PID: $PID"
else
    echo "✗ Lumine 未运行 / Not running"
fi

echo ""

# 检查端口 / Check port
if command -v netstat > /dev/null; then
    echo "端口监听状态 / Port Listening:"
    netstat -an | grep 1080 | grep LISTEN || echo "  端口 1080 未监听 / Port 1080 not listening"
fi

echo ""
echo "═══════════════════════════════════════"
