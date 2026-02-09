# Lumine Network Preprocessing Integration

## Overview

Lumine network preprocessing has been successfully integrated into NekoBoxForAndroid to provide TLS fragmentation and censorship bypass capabilities.

## What is Lumine?

Lumine is a lightweight local HTTP/SOCKS5 proxy server that protects TLS connections over TCP using:
- TLS record fragmentation (tls-rf mode)
- TTL-based desynchronization with fake packets (ttl-d mode)
- TCP segmentation
- DNS query caching and policy-based routing

Based on [TlsFragment](https://github.com/maoist2009/TlsFragment) technology, it helps bypass network censorship and DPI (Deep Packet Inspection).

## Integration Architecture

The lumine library has been integrated at the libcore level:

```
NekoBoxForAndroid (Android App)
    ↓
libcore.aar (Gomobile-generated Android library)
    ↓
libcore/lumine (Go package)
    ↓
Lumine network preprocessing (SOCKS5/HTTP proxy with TLS protection)
```

## Usage from Android/Kotlin

### 1. Create a Lumine Service with Configuration File

```kotlin
import libcore.LumineService

// Create service with config file path
val lumineService = Libcore.newLumineService("/path/to/config.json")

// Use the service...

// Close when done
lumineService.close()
```

### 2. Create a Lumine Service with Explicit Addresses

```kotlin
// Create service with specific addresses
// Parameters: socks5Address, httpAddress, configJSON
val lumineService = Libcore.newLumineServiceWithAddrs(
    "127.0.0.1:1080",  // SOCKS5 address
    "127.0.0.1:1225",  // HTTP proxy address
    ""                 // Optional config JSON (empty for now)
)

// Use the service...

// Close when done
lumineService.close()
```

### 3. Test Integration

```kotlin
// Test if lumine is available
val testResult = Libcore.testLumineIntegration()
println(testResult) // "Lumine network preprocessing is integrated successfully"
```

## Configuration File Format

The lumine configuration uses JSON format. Here's an example:

```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "127.0.0.1:1225",
  "dns_addr": "8.8.8.8:53",
  "udp_minsize": 4096,
  "socks5_for_doh": "",
  "max_jump": 30,
  "fake_ttl_rules": "0-1;3=3;5-1;8-2;13-3;20=18",
  "transmit_file_limit": 2,
  "dns_cache_ttl": 259200,
  "fake_ttl_cache_ttl": 259200,
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 10,
    "num_segs": 3,
    "send_interval": "200ms"
  },
  "domain_policies": {
    "example.com": {
      "mode": "ttl-d",
      "fake_ttl": 17
    }
  },
  "ip_policies": {}
}
```

### Configuration Fields

#### Top-Level Fields
- `socks5_address`: SOCKS5 bind address (use "none" to disable)
- `http_address`: HTTP proxy bind address (use "none" to disable)
- `dns_addr`: DNS server address for resolution
- `udp_minsize`: Minimum UDP packet size for DNS queries
- `socks5_for_doh`: SOCKS5 proxy for DoH (empty to disable)
- `max_jump`: Maximum redirect chain length for IP mapping
- `fake_ttl_rules`: TTL calculation rules for fake packets
- `transmit_file_limit`: Maximum concurrent TransmitFile operations
- `dns_cache_ttl`: DNS answer cache duration in seconds (-1 = forever, 0 = disabled)
- `fake_ttl_cache_ttl`: TTL cache duration in seconds (-1 = forever, 0 = disabled)
- `default_policy`: Default policy for all connections
- `domain_policies`: Domain-specific policies
- `ip_policies`: IP/CIDR-specific policies

#### Policy Fields
- `connect_timeout`: Maximum connection establishment time (e.g., "10s")
- `reply_first`: Send SOCKS5 reply SUCCESS before connecting (boolean)
- `host`: Override target host
- `port`: Override target port
- `mode`: Traffic manipulation mode
  - `raw`: Raw TCP forwarding
  - `direct`: Pass-through without manipulation
  - `tls-rf`: TLS record fragmentation
  - `ttl-d`: TTL-based desynchronization
  - `block`: Block connection
  - `tls-alert`: Send TLS alert and terminate
- `num_records`: Number of TLS records for fragmentation
- `num_segs`: Number of segments for TCP fragmentation
- `send_interval`: Interval between sending segments
- `fake_ttl`: TTL value for fake packets (0 = auto detection)
- `fake_sleep`: Sleep time after sending fake packet

## How It Works

### TLS Record Fragmentation (tls-rf mode)
1. Client sends TLS ClientHello
2. Lumine splits the ClientHello into multiple TLS records
3. Each record is sent separately with configurable delays
4. This bypasses DPI systems that analyze complete ClientHello packets

### TTL-Based Desynchronization (ttl-d mode)
1. Lumine sends fake packets with low TTL that won't reach destination
2. Follows with real packets that have normal TTL
3. DPI systems may be confused by the fake packets
4. Auto-detection can find optimal TTL values

### DNS Caching
- Reduces DNS queries
- Configurable TTL for cached entries
- Supports both A and AAAA records
- Optional retry with dual query

## Integration into NekoBoxForAndroid Workflow

To integrate lumine into your VPN workflow:

1. **Start Lumine service** when VPN connects
2. **Route traffic** through local SOCKS5/HTTP proxy
3. **Apply policies** based on domain/IP patterns
4. **Stop service** when VPN disconnects

Example:
```kotlin
class VpnService {
    private var lumineService: LumineService? = null
    
    fun onVpnStart() {
        // Start lumine proxy
        lumineService = Libcore.newLumineServiceWithAddrs(
            "127.0.0.1:1080",
            "127.0.0.1:1225",
            ""
        )
        
        // Configure your VPN to route through 127.0.0.1:1080
    }
    
    fun onVpnStop() {
        // Clean up
        lumineService?.close()
        lumineService = null
    }
}
```

## Build Information

The lumine library is compiled into libcore.aar during the build process:

```bash
cd libcore
bash init.sh      # Install gomobile-matsuri (first time only)
bash build.sh     # Build libcore.aar with lumine
```

The resulting libcore.aar includes:
- LumineService Java class
- Native libraries for all Android architectures (arm, arm64, x86, x86_64)
- All lumine functionality accessible from Android/Java/Kotlin

## Source Code

The lumine source code is located in:
- `libcore/lumine/` - Core lumine Go package
- `libcore/lumine_wrapper.go` - Android integration wrapper

## Credits

- Original lumine project: https://github.com/UcnacDx2/lumine
- Based on TlsFragment: https://github.com/maoist2009/TlsFragment
- Integrated into NekoBoxForAndroid: https://github.com/MatsuriDayo/NekoBoxForAndroid
