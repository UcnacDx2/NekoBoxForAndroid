# NekoBoxForAndroid Lumine链式代理配置完整解答

## 问题解答

### 问题1: 如何在NekoBoxForAndroid客户端中使用一个配置文件启用代理的VPN网络？

#### 答案：

有两种方式可以使用配置文件启用lumine代理：

#### 方式A: 使用Lumine配置文件

1. **创建Lumine配置文件** (`/sdcard/lumine_config.json`)：

```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "8.8.8.8:53",
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 10,
    "num_segs": 3,
    "send_interval": "200ms"
  }
}
```

2. **在代码中启动Lumine服务**：

```kotlin
val lumineService = Libcore.newLumineService("/sdcard/lumine_config.json")
```

#### 方式B: 在sing-box配置中集成Lumine

创建完整的sing-box配置文件 (`/sdcard/config.json`)：

```json
{
  "log": {"level": "info"},
  "inbounds": [
    {
      "type": "tun",
      "tag": "tun-in",
      "inet4_address": "172.19.0.1/28",
      "auto_route": true,
      "sniff": true
    }
  ],
  "outbounds": [
    {
      "type": "socks",
      "tag": "lumine",
      "server": "127.0.0.1",
      "server_port": 1080,
      "version": "5"
    },
    {
      "type": "shadowsocks",
      "tag": "proxy",
      "server": "your.server.com",
      "server_port": 8388,
      "method": "aes-256-gcm",
      "password": "your-password",
      "detour": "lumine"
    },
    {
      "type": "direct",
      "tag": "direct"
    }
  ],
  "route": {
    "auto_detect_interface": true,
    "final": "proxy"
  }
}
```

然后在代码中：

```kotlin
// 1. 启动lumine
val lumineService = Libcore.newLumineService("/sdcard/lumine_config.json")

// 2. 使用sing-box配置
val config = File("/sdcard/config.json").readText()
val boxInstance = Libcore.newSingBoxInstance(config, null)
boxInstance.start()
```

---

### 问题2: 如何基于Lumine作为上游实现链式代理？

#### 答案：

链式代理的核心是使用sing-box的 **`detour`** 机制。流量路径如下：

```
用户流量 → TUN接口 → sing-box路由 → lumine(detour) → 实际代理服务器 → 互联网
```

#### 实现步骤：

**步骤1: 定义Lumine出站**

在sing-box配置的 `outbounds` 数组中添加：

```json
{
  "type": "socks",
  "tag": "lumine",
  "server": "127.0.0.1",
  "server_port": 1080,
  "version": "5"
}
```

**步骤2: 为实际代理添加detour**

在您的实际代理配置中添加 `"detour": "lumine"`：

```json
{
  "type": "shadowsocks",
  "tag": "proxy",
  "server": "your.server.com",
  "server_port": 8388,
  "method": "aes-256-gcm",
  "password": "your-password",
  "detour": "lumine"          ← 关键配置
}
```

**步骤3: 配置路由规则**

```json
{
  "route": {
    "auto_detect_interface": true,
    "rules": [
      {"geoip": "cn", "outbound": "direct"}
    ],
    "final": "proxy"
  }
}
```

#### 完整的链式代理配置示例：

##### Shadowsocks链式代理

```json
{
  "outbounds": [
    {
      "type": "socks",
      "tag": "lumine",
      "server": "127.0.0.1",
      "server_port": 1080
    },
    {
      "type": "shadowsocks",
      "tag": "ss-main",
      "server": "proxy.example.com",
      "server_port": 8388,
      "method": "aes-256-gcm",
      "password": "password",
      "detour": "lumine"
    },
    {
      "type": "direct",
      "tag": "direct"
    }
  ]
}
```

##### VMess链式代理

```json
{
  "outbounds": [
    {
      "type": "socks",
      "tag": "lumine",
      "server": "127.0.0.1",
      "server_port": 1080
    },
    {
      "type": "vmess",
      "tag": "vmess-main",
      "server": "proxy.example.com",
      "server_port": 443,
      "uuid": "your-uuid",
      "security": "auto",
      "detour": "lumine",
      "tls": {
        "enabled": true,
        "server_name": "proxy.example.com"
      }
    }
  ]
}
```

##### 多代理选择器链式代理

```json
{
  "outbounds": [
    {
      "type": "socks",
      "tag": "lumine",
      "server": "127.0.0.1",
      "server_port": 1080
    },
    {
      "type": "shadowsocks",
      "tag": "ss1",
      "server": "server1.com",
      "server_port": 8388,
      "method": "aes-256-gcm",
      "password": "pass1",
      "detour": "lumine"
    },
    {
      "type": "shadowsocks",
      "tag": "ss2",
      "server": "server2.com",
      "server_port": 8388,
      "method": "aes-256-gcm",
      "password": "pass2",
      "detour": "lumine"
    },
    {
      "type": "selector",
      "tag": "proxy",
      "outbounds": ["ss1", "ss2"],
      "default": "ss1"
    }
  ]
}
```

---

## 实际应用代码示例

### 在VPN Service中集成Lumine链式代理

```kotlin
package your.package.name

import android.net.VpnService
import libcore.Libcore
import libcore.LumineService
import libcore.BoxInstance
import java.io.File

class MyVpnService : VpnService() {
    private var lumineService: LumineService? = null
    private var boxInstance: BoxInstance? = null
    
    override fun onCreate() {
        super.onCreate()
    }
    
    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        startVpnWithLumine()
        return START_STICKY
    }
    
    private fun startVpnWithLumine() {
        try {
            // 1. 启动Lumine服务
            val lumineConfigPath = "/sdcard/lumine_config.json"
            lumineService = Libcore.newLumineService(lumineConfigPath)
            
            // 2. 准备sing-box配置（包含lumine detour）
            val singboxConfig = prepareSingBoxConfig()
            
            // 3. 启动sing-box
            boxInstance = Libcore.newSingBoxInstance(singboxConfig, null)
            boxInstance?.start()
            
            Log.i(TAG, "VPN with Lumine chain proxy started successfully")
        } catch (e: Exception) {
            Log.e(TAG, "Failed to start VPN with Lumine", e)
            stopSelf()
        }
    }
    
    private fun prepareSingBoxConfig(): String {
        // 读取或生成配置
        val configFile = File("/sdcard/singbox_config.json")
        if (configFile.exists()) {
            return configFile.readText()
        }
        
        // 生成默认配置
        return """
        {
          "log": {"level": "info"},
          "inbounds": [
            {
              "type": "tun",
              "tag": "tun-in",
              "inet4_address": "172.19.0.1/28",
              "auto_route": true,
              "sniff": true
            }
          ],
          "outbounds": [
            {
              "type": "socks",
              "tag": "lumine",
              "server": "127.0.0.1",
              "server_port": 1080
            },
            {
              "type": "shadowsocks",
              "tag": "proxy",
              "server": "${getProxyServer()}",
              "server_port": ${getProxyPort()},
              "method": "aes-256-gcm",
              "password": "${getProxyPassword()}",
              "detour": "lumine"
            },
            {
              "type": "direct",
              "tag": "direct"
            }
          ],
          "route": {
            "auto_detect_interface": true,
            "final": "proxy"
          }
        }
        """.trimIndent()
    }
    
    override fun onDestroy() {
        stopVpn()
        super.onDestroy()
    }
    
    private fun stopVpn() {
        try {
            boxInstance?.close()
            boxInstance = null
            
            lumineService?.close()
            lumineService = null
            
            Log.i(TAG, "VPN stopped")
        } catch (e: Exception) {
            Log.e(TAG, "Error stopping VPN", e)
        }
    }
    
    companion object {
        private const val TAG = "MyVpnService"
    }
}
```

---

## 关键配置说明

### Detour机制工作原理

当配置了 `"detour": "lumine"` 后：

1. **流量首先到达代理出站** (例如shadowsocks)
2. **代理出站查找detour标签** ("lumine")
3. **流量被重定向到lumine出站** (127.0.0.1:1080)
4. **Lumine处理流量** (TLS分片、TTL修改等)
5. **处理后的流量发送到实际服务器**

```
客户端 → proxy (detour="lumine") → lumine (127.0.0.1:1080) → 实际服务器
```

### Lumine处理流程

```
原始TLS ClientHello (1个包)
    ↓
Lumine TLS分片 (mode="tls-rf")
    ↓
分成10个TLS记录 (num_records=10)
    ↓
每个记录分成3个TCP段 (num_segs=3)
    ↓
延迟发送 (send_interval="200ms")
    ↓
绕过DPI检测
    ↓
到达代理服务器
```

---

## 配置文件位置建议

### Android文件系统位置

- **公共存储**: `/sdcard/Download/` 或 `/sdcard/Documents/`
  - 优点: 用户可以直接编辑
  - 缺点: 需要存储权限
  
- **应用私有目录**: `context.filesDir` 或 `context.getExternalFilesDir(null)`
  - 优点: 不需要额外权限，更安全
  - 缺点: 用户不易访问

示例：

```kotlin
// 在应用私有目录创建配置
val lumineConfig = File(filesDir, "lumine_config.json")
if (!lumineConfig.exists()) {
    lumineConfig.writeText(getDefaultLumineConfig())
}

val singboxConfig = File(filesDir, "singbox_config.json")
if (!singboxConfig.exists()) {
    singboxConfig.writeText(getDefaultSingBoxConfig())
}
```

---

## 提供的示例文件

本仓库提供了以下完整配置示例：

1. **`singbox_lumine_chain_example.json`**
   - 完整的sing-box配置
   - 包含Shadowsocks + Lumine链式代理
   - 带DNS、路由规则

2. **`singbox_vmess_lumine_example.json`**
   - VMess协议 + Lumine
   - 包含TLS和WebSocket传输

3. **`lumine_config_advanced.json`**
   - 高级Lumine配置
   - 包含域名特定策略
   - 针对Google、Twitter、GitHub等不同配置

4. **`lumine_config_example.json`**
   - 基础Lumine配置
   - 适合新手使用

---

## 故障排除

### 常见问题

**问题**: Lumine无法启动，提示端口已占用

**解决**: 修改配置文件中的端口
```json
{
  "socks5_address": "127.0.0.1:10800"
}
```

**问题**: 流量没有通过Lumine

**解决**: 检查detour配置是否正确，确保：
1. lumine出站的tag为 "lumine"
2. 代理出站有 `"detour": "lumine"`

**问题**: 连接很慢

**解决**: 调整分片参数
```json
{
  "default_policy": {
    "num_records": 5,      // 减少记录数
    "send_interval": "50ms" // 减少延迟
  }
}
```

---

## 总结

使用Lumine实现链式代理的完整流程：

1. ✅ 创建Lumine配置文件 (`lumine_config.json`)
2. ✅ 创建sing-box配置文件，包含：
   - Lumine SOCKS5出站 (tag="lumine")
   - 实际代理出站 + `detour="lumine"`
3. ✅ 在代码中先启动Lumine服务
4. ✅ 然后启动sing-box服务
5. ✅ 流量自动通过Lumine预处理后到达代理服务器

参考完整文档：
- [LUMINE_CHAIN_PROXY_GUIDE.md](LUMINE_CHAIN_PROXY_GUIDE.md)
- [LUMINE_INTEGRATION.md](LUMINE_INTEGRATION.md)
