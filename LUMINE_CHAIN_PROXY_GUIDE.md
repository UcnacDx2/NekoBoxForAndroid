# Lumine Chain Proxy Configuration Guide

## Overview

This guide explains how to configure NekoBoxForAndroid to use lumine as an upstream proxy (chain proxy), providing TLS fragmentation and censorship bypass for your VPN connections.

## What is Chain Proxy?

Chain proxy (also called upstream proxy or proxy chain) means routing your VPN traffic through multiple proxies in sequence:

```
Your Device → NekoBoxForAndroid → Lumine (preprocessing) → Your Proxy Server → Internet
```

Benefits:
- **TLS Fragmentation**: Lumine fragments TLS handshakes before reaching your proxy server
- **Bypass DPI**: Helps bypass Deep Packet Inspection at the local network level
- **Policy-Based**: Different fragmentation policies for different domains/IPs
- **Transparent**: Works with any proxy protocol (SOCKS5, HTTP, Shadowsocks, VMess, etc.)

## Quick Start

### Method 1: Use Lumine Configuration File

1. **Create a lumine configuration file** (e.g., `/sdcard/lumine_config.json`):

```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "8.8.8.8:53",
  "udp_minsize": 4096,
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 10,
    "num_segs": 3,
    "send_interval": "200ms"
  },
  "domain_policies": {},
  "ip_policies": {}
}
```

2. **In your NekoBoxForAndroid configuration**, add lumine as a SOCKS5 outbound:

```json
{
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
    }
  ]
}
```

The `"detour": "lumine"` makes the proxy route through lumine first.

### Method 2: Programmatic Configuration

```kotlin
// 1. Start lumine service
val lumineConfig = File("/sdcard/lumine_config.json").readText()
val lumineService = Libcore.newLumineService("/sdcard/lumine_config.json")

// 2. Configure your sing-box config to use lumine as detour
val config = """
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
      "tag": "proxy",
      "server": "your.server.com",
      "server_port": 8388,
      "method": "aes-256-gcm",
      "password": "your-password",
      "detour": "lumine"
    }
  ]
}
"""

// 3. Use the config with NekoBox
// ... create box instance with config ...

// 4. When done, close lumine
lumineService.close()
```

## Configuration Examples

### Example 1: Shadowsocks with Lumine Preprocessing

```json
{
  "log": {
    "level": "info"
  },
  "dns": {
    "servers": [
      {
        "tag": "dns_proxy",
        "address": "tls://8.8.8.8",
        "detour": "proxy"
      }
    ]
  },
  "inbounds": [
    {
      "type": "tun",
      "tag": "tun-in",
      "inet4_address": "172.19.0.1/28"
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
      "server": "example.com",
      "server_port": 8388,
      "method": "aes-256-gcm",
      "password": "password",
      "detour": "lumine"
    },
    {
      "type": "direct",
      "tag": "direct"
    },
    {
      "type": "block",
      "tag": "block"
    }
  ],
  "route": {
    "auto_detect_interface": true,
    "rules": [
      {
        "geoip": "cn",
        "outbound": "direct"
      },
      {
        "geosite": "cn",
        "outbound": "direct"
      }
    ],
    "final": "proxy"
  }
}
```

### Example 2: VMess with Lumine Preprocessing

```json
{
  "outbounds": [
    {
      "type": "socks",
      "tag": "lumine",
      "server": "127.0.0.1",
      "server_port": 1080,
      "version": "5"
    },
    {
      "type": "vmess",
      "tag": "proxy",
      "server": "example.com",
      "server_port": 443,
      "uuid": "your-uuid-here",
      "security": "auto",
      "alter_id": 0,
      "detour": "lumine",
      "tls": {
        "enabled": true,
        "server_name": "example.com"
      }
    },
    {
      "type": "direct",
      "tag": "direct"
    }
  ]
}
```

### Example 3: Multiple Proxies with Selective Lumine Usage

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
      "tag": "ss-with-lumine",
      "server": "censored-region.com",
      "server_port": 8388,
      "method": "aes-256-gcm",
      "password": "password",
      "detour": "lumine"
    },
    {
      "type": "shadowsocks",
      "tag": "ss-direct",
      "server": "free-region.com",
      "server_port": 8388,
      "method": "aes-256-gcm",
      "password": "password"
    },
    {
      "type": "selector",
      "tag": "proxy",
      "outbounds": ["ss-with-lumine", "ss-direct"],
      "default": "ss-with-lumine"
    }
  ]
}
```

## Lumine Configuration Options

### Basic Configuration

```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "127.0.0.1:1225",
  "dns_addr": "8.8.8.8:53"
}
```

### Advanced TLS Fragmentation

```json
{
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 10,
    "num_segs": 3,
    "send_interval": "200ms"
  }
}
```

- `mode`: `"tls-rf"` (TLS record fragmentation), `"ttl-d"` (TTL desync), `"direct"` (no manipulation)
- `num_records`: Number of TLS records to split into (higher = more fragmentation)
- `num_segs`: Number of TCP segments per record
- `send_interval`: Delay between sending fragments

### TTL Desynchronization

```json
{
  "default_policy": {
    "mode": "ttl-d",
    "fake_ttl": 17,
    "fake_sleep": "200ms"
  }
}
```

- `fake_ttl`: TTL value for fake packets (0 = auto-detect)
- `fake_sleep`: Delay after sending fake packet

### Domain-Specific Policies

```json
{
  "domain_policies": {
    "example.com;*.example.com": {
      "mode": "tls-rf",
      "num_records": 20
    },
    "google.com": {
      "mode": "ttl-d",
      "fake_ttl": 15
    },
    "github.com": {
      "mode": "direct"
    }
  }
}
```

### IP-Specific Policies

```json
{
  "ip_policies": {
    "1.1.1.1": {
      "mode": "tls-rf"
    },
    "8.8.8.8/32": {
      "mode": "ttl-d"
    }
  }
}
```

## Integration with NekoBoxForAndroid

### Starting Lumine Service

Add this code to your VPN service initialization:

```kotlin
class MyVpnService : VpnService() {
    private var lumineService: LumineService? = null
    
    override fun onCreate() {
        super.onCreate()
        
        // Start lumine with config file
        val configFile = File(filesDir, "lumine_config.json")
        if (configFile.exists()) {
            lumineService = Libcore.newLumineService(configFile.absolutePath)
        } else {
            // Or use default addresses
            lumineService = Libcore.newLumineServiceWithAddrs(
                "127.0.0.1:1080",
                "none",
                ""
            )
        }
    }
    
    override fun onDestroy() {
        lumineService?.close()
        lumineService = null
        super.onDestroy()
    }
}
```

### Updating sing-box Configuration

Modify your config builder to add lumine outbound:

```kotlin
fun buildConfigWithLumine(originalConfig: String): String {
    val config = gson.fromJson(originalConfig, SingBoxOptions::class.java)
    
    // Add lumine outbound
    val lumineOutbound = Outbound_SocksOptions().apply {
        type = "socks"
        tag = "lumine"
        server = "127.0.0.1"
        server_port = 1080
        version = "5"
    }
    config.outbounds.add(0, lumineOutbound)
    
    // Set detour for main proxy
    config.outbounds.firstOrNull { it.tag == "proxy" }?.let {
        it._hack_config_map["detour"] = "lumine"
    }
    
    return gson.toJson(config)
}
```

## Troubleshooting

### Lumine Not Starting

Check if the port is available:
```bash
netstat -an | grep 1080
```

Try a different port in config:
```json
{
  "socks5_address": "127.0.0.1:10800"
}
```

### Connection Timeout

- Check lumine is running: `Libcore.testLumineIntegration()`
- Verify SOCKS5 address in config matches lumine address
- Check firewall/SELinux settings

### No Fragmentation Effect

- Verify `mode` is set to `"tls-rf"` or `"ttl-d"`
- Increase `num_records` for more fragmentation
- Try different `send_interval` values

### Performance Issues

- Reduce `num_records` (try 5-10 instead of 20+)
- Reduce `send_interval` (try "50ms" or "100ms")
- Use `"direct"` mode for trusted domains

## Best Practices

1. **Start with defaults**: Use the example config and adjust based on results
2. **Test incrementally**: Enable lumine for one proxy first, then expand
3. **Monitor performance**: Fragmentation adds latency, balance security vs speed
4. **Use domain policies**: Apply fragmentation only where needed
5. **Keep backups**: Save working configurations
6. **Update regularly**: Check for lumine updates and new features

## Performance Tuning

### Low Latency Profile
```json
{
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 5,
    "num_segs": 2,
    "send_interval": "50ms"
  }
}
```

### High Security Profile
```json
{
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 20,
    "num_segs": 5,
    "send_interval": "300ms"
  }
}
```

### Balanced Profile
```json
{
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 10,
    "num_segs": 3,
    "send_interval": "200ms"
  }
}
```

## See Also

- [LUMINE_INTEGRATION.md](LUMINE_INTEGRATION.md) - Basic integration guide
- [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) - Technical details
- [lumine_config_example.json](lumine_config_example.json) - Example configuration
- Original lumine project: https://github.com/UcnacDx2/lumine
