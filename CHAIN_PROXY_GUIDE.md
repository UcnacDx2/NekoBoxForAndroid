# NekoBox 链式代理使用指南 / Chain Proxy User Guide

[中文](#中文文档) | [English](#english-documentation)

---

## 中文文档

### 什么是链式代理？

链式代理（Proxy Chain）是一种将多个代理服务器串联使用的技术。流量会按顺序通过每个代理服务器，最终到达目标服务器。这种方式可以：

- **增强隐私保护**：多层代理使得追踪源头更加困难
- **绕过多层限制**：通过不同地区的代理服务器突破地理限制
- **负载均衡**：将流量分散到多个代理服务器
- **灵活路由**：根据不同需求组合不同类型的代理

### 链式代理工作原理

```
客户端 → 代理1 → 代理2 → 代理3 → 目标服务器
```

在 NekoBox 中，链式代理的数据流向是：
1. 客户端发送请求
2. 请求首先到达链中的第一个代理（最外层）
3. 依次通过链中的每个代理
4. 最后一个代理（最内层）连接到目标服务器

### 如何创建链式代理

#### 方法一：通过界面创建

1. **打开 NekoBox 应用**

2. **添加链式代理配置**
   - 点击右下角的 "+" 按钮
   - 选择 "Proxy Chain"（链式代理）

3. **配置链式代理**
   - 输入配置名称（例如："三层代理链"）
   - 点击 "+" 按钮添加代理
   - 选择要添加到链中的代理（按顺序选择）
   - 可以通过拖拽调整代理顺序
   - 向左滑动可以删除链中的代理

4. **保存配置**
   - 点击右上角的 "✓" 保存

5. **使用链式代理**
   - 在主界面选择刚创建的链式代理
   - 点击连接按钮启动 VPN

#### 重要说明：

- **代理顺序**：链中的第一个代理是最外层（客户端直接连接的代理），最后一个是最内层（连接目标服务器的代理）
- **避免循环**：不能将链式代理添加到自己或包含自己的链中
- **性能影响**：代理层数越多，延迟越高，速度可能越慢

### 配置示例

#### 示例 1：双层代理（本地 → 国内 → 国外）

**使用场景**：通过国内跳板访问国外服务

```
配置名称：双层代理-CN-US
代理链：
1. 国内 SOCKS5 代理（例如：本地服务器）
2. 国外 Shadowsocks 服务器（例如：美国服务器）
```

**创建步骤**：
1. 先添加两个独立的代理配置：
   - 代理A：国内 SOCKS5（例如：192.168.1.100:1080）
   - 代理B：国外 Shadowsocks（例如：ss://xxx@us-server.com:8388）

2. 创建链式代理：
   - 名称：双层代理-CN-US
   - 添加代理A（国内 SOCKS5）
   - 添加代理B（国外 Shadowsocks）

#### 示例 2：三层代理（本地 → 中转 → 跳板 → 目标）

**使用场景**：高度匿名访问，多层保护

```
配置名称：三层匿名代理
代理链：
1. 本地 Trojan 服务器
2. 中转 VMess 服务器
3. 跳板 Shadowsocks 服务器
```

#### 示例 3：结合 Lumine 的链式代理

**使用场景**：使用 Lumine 进行流量混淆，然后通过代理服务器

```
配置名称：Lumine + 代理链
代理链：
1. Lumine 插件（SOCKS5: 127.0.0.1:1080）
2. Shadowsocks 服务器
```

**Lumine 配置文件** (`lumine_config.json`):
```json
{
    "socks5_address": "127.0.0.1:1080",
    "http_address": "none",
    "dns_addr": "https://1.1.1.1/dns-query",
    "default_policy": {
        "mode": "tls-rf",
        "num_records": 10,
        "num_segs": 3
    }
}
```

**创建步骤**：
1. 启动 Lumine 插件（作为本地 SOCKS5 代理）
2. 创建 SOCKS5 代理配置指向 Lumine（127.0.0.1:1080）
3. 创建链式代理，依次添加：
   - Lumine SOCKS5 代理
   - 你的 Shadowsocks/VMess 等服务器

### 将链式代理作为上游代理

NekoBox 支持在分组（Group）级别设置前置代理和落地代理：

#### 前置代理（Front Proxy）

前置代理会添加到该分组所有代理的最外层。

**设置方法**：
1. 进入分组管理
2. 选择一个分组
3. 设置 "Front Proxy"（前置代理）
4. 选择要作为前置的代理或链式代理

**使用场景**：
- 所有流量都先通过本地加密代理
- 统一的流量入口控制

#### 落地代理（Landing Proxy）

落地代理会添加到该分组所有代理的最内层。

**设置方法**：
1. 进入分组管理
2. 选择一个分组
3. 设置 "Landing Proxy"（落地代理）
4. 选择要作为落地的代理

**使用场景**：
- 最终统一通过某个特定出口
- 需要固定落地IP的场景

### 配置文件方式

虽然 NekoBox 主要通过界面配置，但你也可以导出配置文件：

**导出配置**：
1. 长按链式代理配置
2. 选择 "导出" 或 "分享"
3. 保存配置文件

**导入配置**：
1. 点击右上角菜单
2. 选择 "从文件导入"
3. 选择配置文件

### 最佳实践

#### 1. 代理选择原则
- **速度优先**：选择延迟低的代理
- **稳定性优先**：选择长期稳定的服务器
- **地理位置**：根据需求选择合适地区的服务器

#### 2. 链长度建议
- **一般使用**：1-2层代理即可
- **高匿名需求**：可使用3层代理
- **不推荐**：超过3层（性能显著下降）

#### 3. 性能优化
- 使用高速、低延迟的服务器
- 避免在同一地区重复代理
- 定期测试链式代理速度

#### 4. 安全建议
- 定期更换代理服务器
- 使用不同提供商的服务器
- 启用加密协议（如 TLS）

### 故障排除

#### 问题 1：链式代理无法连接

**可能原因**：
- 链中某个代理不可用
- 代理顺序配置错误
- 防火墙阻止连接

**解决方法**：
- 单独测试链中的每个代理
- 检查代理顺序是否正确
- 检查防火墙和网络设置

#### 问题 2：速度很慢

**可能原因**：
- 代理链太长
- 某个代理服务器速度慢
- 路由不优

**解决方法**：
- 减少代理层数
- 更换高速服务器
- 优化代理地理位置

#### 问题 3：DNS 解析失败

**可能原因**：
- DNS 配置错误
- 代理不支持 DNS

**解决方法**：
- 在设置中配置远程 DNS
- 启用 "DNS 路由" 功能
- 使用支持 DNS 的代理协议

### 高级技巧

#### 1. 动态链式代理

结合 NekoBox 的路由功能，可以为不同网站使用不同的代理链：

1. 创建多个链式代理配置
2. 在路由设置中添加规则
3. 为不同域名指定不同的链

#### 2. 负载均衡

使用 Balancer（均衡器）功能：

1. 创建多个相似的链式代理
2. 使用均衡器将它们组合
3. 流量会自动分配到不同的链

#### 3. 故障转移

配置备用链式代理：

1. 创建主链式代理和备用链
2. 使用选择器（Selector）功能
3. 主链不可用时自动切换

---

## English Documentation

### What is Chain Proxy?

Chain Proxy (Proxy Chain) is a technique that connects multiple proxy servers in sequence. Traffic passes through each proxy server in order before reaching the destination. This approach provides:

- **Enhanced Privacy**: Multiple proxy layers make it harder to trace the origin
- **Multi-layer Bypass**: Break through geographical restrictions using proxies from different regions
- **Load Balancing**: Distribute traffic across multiple proxy servers
- **Flexible Routing**: Combine different types of proxies for different needs

### How Chain Proxy Works

```
Client → Proxy1 → Proxy2 → Proxy3 → Destination Server
```

In NekoBox, the data flow for chain proxy is:
1. Client sends a request
2. Request first reaches the first proxy in the chain (outermost)
3. Passes through each proxy in the chain sequentially
4. The last proxy (innermost) connects to the destination server

### How to Create Chain Proxy

#### Method 1: Create via UI

1. **Open NekoBox App**

2. **Add Chain Proxy Configuration**
   - Tap the "+" button in the bottom right
   - Select "Proxy Chain"

3. **Configure Chain Proxy**
   - Enter configuration name (e.g., "Three-Layer Proxy Chain")
   - Tap "+" button to add proxies
   - Select proxies to add to the chain (in order)
   - Drag to reorder proxies
   - Swipe left to remove a proxy from the chain

4. **Save Configuration**
   - Tap "✓" in the top right to save

5. **Use Chain Proxy**
   - Select the newly created chain proxy on the main screen
   - Tap connect button to start VPN

#### Important Notes:

- **Proxy Order**: The first proxy in the chain is the outermost (client connects directly), the last is innermost (connects to destination)
- **Avoid Loops**: Cannot add a chain proxy to itself or to a chain containing itself
- **Performance Impact**: More proxy layers mean higher latency and potentially slower speed

### Configuration Examples

#### Example 1: Two-Layer Proxy (Local → Domestic → Foreign)

**Use Case**: Access foreign services through domestic relay

```
Config Name: Two-Layer-CN-US
Proxy Chain:
1. Domestic SOCKS5 Proxy (e.g., local server)
2. Foreign Shadowsocks Server (e.g., US server)
```

**Setup Steps**:
1. First add two independent proxy configurations:
   - Proxy A: Domestic SOCKS5 (e.g., 192.168.1.100:1080)
   - Proxy B: Foreign Shadowsocks (e.g., ss://xxx@us-server.com:8388)

2. Create chain proxy:
   - Name: Two-Layer-CN-US
   - Add Proxy A (Domestic SOCKS5)
   - Add Proxy B (Foreign Shadowsocks)

#### Example 2: Three-Layer Proxy (Local → Relay → Jump → Destination)

**Use Case**: Highly anonymous access with multi-layer protection

```
Config Name: Three-Layer-Anonymous
Proxy Chain:
1. Local Trojan Server
2. Relay VMess Server
3. Jump Shadowsocks Server
```

#### Example 3: Chain Proxy with Lumine

**Use Case**: Use Lumine for traffic obfuscation, then through proxy server

```
Config Name: Lumine + Proxy Chain
Proxy Chain:
1. Lumine Plugin (SOCKS5: 127.0.0.1:1080)
2. Shadowsocks Server
```

**Lumine Config File** (`lumine_config.json`):
```json
{
    "socks5_address": "127.0.0.1:1080",
    "http_address": "none",
    "dns_addr": "https://1.1.1.1/dns-query",
    "default_policy": {
        "mode": "tls-rf",
        "num_records": 10,
        "num_segs": 3
    }
}
```

**Setup Steps**:
1. Start Lumine plugin (as local SOCKS5 proxy)
2. Create SOCKS5 proxy config pointing to Lumine (127.0.0.1:1080)
3. Create chain proxy, add in sequence:
   - Lumine SOCKS5 proxy
   - Your Shadowsocks/VMess server

### Using Chain Proxy as Upstream

NekoBox supports setting front proxy and landing proxy at the Group level:

#### Front Proxy

Front proxy will be added to the outermost layer of all proxies in the group.

**Setup Method**:
1. Enter group management
2. Select a group
3. Set "Front Proxy"
4. Choose a proxy or chain proxy as front proxy

**Use Cases**:
- All traffic first goes through local encrypted proxy
- Unified traffic entry control

#### Landing Proxy

Landing proxy will be added to the innermost layer of all proxies in the group.

**Setup Method**:
1. Enter group management
2. Select a group
3. Set "Landing Proxy"
4. Choose a proxy as landing proxy

**Use Cases**:
- Final unified exit through a specific proxy
- Scenarios requiring fixed landing IP

### Configuration File Method

While NekoBox primarily uses UI configuration, you can also export configuration files:

**Export Config**:
1. Long press on chain proxy config
2. Select "Export" or "Share"
3. Save configuration file

**Import Config**:
1. Tap menu in top right
2. Select "Import from File"
3. Choose configuration file

### Best Practices

#### 1. Proxy Selection Principles
- **Speed Priority**: Choose low-latency proxies
- **Stability Priority**: Choose long-term stable servers
- **Geographic Location**: Choose appropriate regions based on needs

#### 2. Chain Length Recommendations
- **General Use**: 1-2 layers sufficient
- **High Anonymity Needs**: Can use 3 layers
- **Not Recommended**: More than 3 layers (significant performance degradation)

#### 3. Performance Optimization
- Use high-speed, low-latency servers
- Avoid duplicate proxies in the same region
- Regularly test chain proxy speed

#### 4. Security Recommendations
- Regularly change proxy servers
- Use servers from different providers
- Enable encryption protocols (such as TLS)

### Troubleshooting

#### Issue 1: Chain Proxy Cannot Connect

**Possible Causes**:
- One proxy in the chain is unavailable
- Incorrect proxy order configuration
- Firewall blocking connection

**Solutions**:
- Test each proxy in the chain individually
- Check if proxy order is correct
- Check firewall and network settings

#### Issue 2: Very Slow Speed

**Possible Causes**:
- Proxy chain too long
- One proxy server is slow
- Poor routing

**Solutions**:
- Reduce number of proxy layers
- Replace with high-speed servers
- Optimize proxy geographic locations

#### Issue 3: DNS Resolution Failure

**Possible Causes**:
- Incorrect DNS configuration
- Proxy doesn't support DNS

**Solutions**:
- Configure remote DNS in settings
- Enable "DNS Routing" feature
- Use proxy protocol that supports DNS

### Advanced Tips

#### 1. Dynamic Chain Proxy

Combine with NekoBox routing to use different proxy chains for different websites:

1. Create multiple chain proxy configurations
2. Add rules in routing settings
3. Specify different chains for different domains

#### 2. Load Balancing

Use Balancer feature:

1. Create multiple similar chain proxies
2. Use balancer to combine them
3. Traffic automatically distributed to different chains

#### 3. Failover

Configure backup chain proxies:

1. Create main chain proxy and backup chains
2. Use Selector feature
3. Automatically switch when main chain unavailable

---

## 技术原理 / Technical Principles

### 链式代理的实现

NekoBox 使用 sing-box 核心实现链式代理功能：

1. **ChainBean**: 存储代理链配置（代理 ID 列表）
2. **ConfigBuilder**: 将代理链转换为 sing-box 配置
3. **resolveChain**: 递归解析代理链，支持嵌套链

### 配置生成流程

```
ChainBean (代理ID列表)
    ↓
resolveChainInternal() - 递归解析
    ↓
添加 frontProxy (前置代理)
    ↓
添加 landingProxy (落地代理)
    ↓
生成 sing-box outbound 配置
```

### Chain Proxy Implementation

NekoBox uses sing-box core to implement chain proxy functionality:

1. **ChainBean**: Stores chain proxy configuration (list of proxy IDs)
2. **ConfigBuilder**: Converts proxy chain to sing-box configuration
3. **resolveChain**: Recursively resolves proxy chain, supports nested chains

### Configuration Generation Flow

```
ChainBean (List of Proxy IDs)
    ↓
resolveChainInternal() - Recursive resolution
    ↓
Add frontProxy (Front Proxy)
    ↓
Add landingProxy (Landing Proxy)
    ↓
Generate sing-box outbound config
```

---

## 相关资源 / Related Resources

- **NekoBox 项目主页**: https://matsuridayo.github.io
- **sing-box 文档**: https://sing-box.sagernet.org
- **Lumine 用户指南**: [LUMINE_USER_GUIDE.md](./LUMINE_USER_GUIDE.md)
- **NekoBox GitHub**: https://github.com/MatsuriDayo/NekoBoxForAndroid

---

## 许可证 / License

本文档遵循 GPL-3.0 许可证

This documentation follows GPL-3.0 License
