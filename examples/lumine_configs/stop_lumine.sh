#!/system/bin/sh
# 停止 Lumine 服务
# Stop Lumine Service

echo "→ 停止 Lumine / Stopping Lumine..."

# 查找并终止 Lumine 进程 / Find and kill Lumine process
PID=$(pgrep -f "liblumine.so")

if [ -z "$PID" ]; then
    echo "✗ Lumine 未运行 / Lumine is not running"
    exit 0
fi

kill $PID
sleep 1

# 确认已停止 / Confirm stopped
if pgrep -f "liblumine.so" > /dev/null; then
    echo "→ 强制停止 / Force killing..."
    kill -9 $PID
    sleep 1
fi

if pgrep -f "liblumine.so" > /dev/null; then
    echo "✗ 停止失败 / Failed to stop"
    exit 1
else
    echo "✓ Lumine 已停止 / Lumine stopped"
fi
