# 链式代理配置示例 / Chain Proxy Configuration Examples

本目录包含各种链式代理配置示例，帮助用户快速开始使用 NekoBox 的链式代理功能。

This directory contains various chain proxy configuration examples to help users quickly get started with NekoBox's chain proxy feature.

## 示例列表 / Example List

### 1. 基础双层代理 / Basic Two-Layer Proxy
**文件**: `chain_two_layer_example.md`

**场景**: 本地 SOCKS5 → Shadowsocks  
**Scenario**: Local SOCKS5 → Shadowsocks

适用于需要通过本地代理访问远程服务器的场景。  
Suitable for scenarios where you need to access remote servers through a local proxy.

### 2. Lumine 增强代理链 / Lumine Enhanced Proxy Chain
**文件**: `chain_lumine_example.md`

**场景**: Lumine (流量混淆) → Shadowsocks  
**Scenario**: Lumine (Traffic Obfuscation) → Shadowsocks

适用于需要绕过深度包检测（DPI）的场景。  
Suitable for scenarios requiring DPI bypass.

## 使用方法 / Usage Instructions

### 通过界面配置 / Configure via UI

1. 在 NekoBox 中分别添加链中的每个代理
2. 创建新的链式代理配置
3. 按示例顺序添加代理
4. 保存并连接

1. Add each proxy in the chain separately in NekoBox
2. Create a new chain proxy configuration
3. Add proxies in the order shown in examples
4. Save and connect

## 配置参数说明 / Configuration Parameter Descriptions

### 代理顺序 / Proxy Order

⚠️ **重要**: 链式代理中的顺序很重要！  
⚠️ **Important**: Order matters in chain proxy!

```
客户端 / Client
    ↓
代理1（最外层）/ Proxy 1 (Outermost)
    ↓
代理2（中间层）/ Proxy 2 (Middle)
    ↓
代理3（最内层）/ Proxy 3 (Innermost)
    ↓
目标服务器 / Destination Server
```

## 性能建议 / Performance Recommendations

### 代理层数 / Number of Layers
- ✅ 推荐 / Recommended: 1-2 层 / layers
- ⚠️ 可接受 / Acceptable: 3 层 / layers
- ❌ 不推荐 / Not Recommended: 4+ 层 / layers

### 服务器选择 / Server Selection
- 选择低延迟服务器 / Choose low-latency servers
- 避免地理位置重复 / Avoid duplicate geographic locations
- 优先使用高带宽服务器 / Prioritize high-bandwidth servers

## 故障排除 / Troubleshooting

### 连接失败 / Connection Failed
1. 检查链中每个代理是否可用 / Check if each proxy in the chain is available
2. 验证代理顺序是否正确 / Verify proxy order is correct
3. 查看 NekoBox 日志获取详细错误 / Check NekoBox logs for detailed errors

### 速度慢 / Slow Speed
1. 减少代理层数 / Reduce number of proxy layers
2. 更换高速服务器 / Replace with high-speed servers
3. 检查每个代理的延迟 / Check latency of each proxy

### DNS 问题 / DNS Issues
1. 启用远程 DNS / Enable remote DNS
2. 配置 DNS 路由 / Configure DNS routing
3. 使用支持 DNS 的代理协议 / Use proxy protocol that supports DNS

## 安全提示 / Security Tips

⚠️ **注意事项 / Cautions**:

1. 不要在配置文件中明文保存密码 / Don't save passwords in plaintext
2. 定期更换代理服务器 / Regularly change proxy servers
3. 使用加密协议（TLS/WS+TLS） / Use encrypted protocols (TLS/WS+TLS)
4. 避免使用公共代理服务器 / Avoid using public proxy servers

## 更多资源 / More Resources

- 📖 [链式代理完整指南](../../CHAIN_PROXY_GUIDE.md) / [Complete Chain Proxy Guide](../../CHAIN_PROXY_GUIDE.md)
- 📖 [Lumine 用户指南](../../LUMINE_USER_GUIDE.md) / [Lumine User Guide](../../LUMINE_USER_GUIDE.md)
- 🌐 [NekoBox 官方文档](https://matsuridayo.github.io) / [NekoBox Official Docs](https://matsuridayo.github.io)

## 许可证 / License

GPL-3.0
