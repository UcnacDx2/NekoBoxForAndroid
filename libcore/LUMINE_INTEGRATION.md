# Lumine Plugin Integration

This document describes the integration of the lumine network preprocessing tool into NekoBoxForAndroid as a native plugin.

## Overview

Lumine is a network preprocessing tool that provides:
- SOCKS5 proxy with traffic manipulation capabilities
- HTTP/HTTPS proxy support
- DNS manipulation and fake TTL support
- Packet fragmentation and desynchronization
- Policy-based traffic routing

## Integration Approach

Lumine is integrated as a **standalone native binary** (similar to hysteria-plugin and naive-plugin) rather than as a Go library. This approach:
- Maintains separation of concerns
- Allows lumine to run as an independent process
- Enables easier updates and debugging
- Follows the existing plugin architecture pattern

## Directory Structure

```
libcore/
├── lumine/                    # Lumine library package
│   ├── config.go              # Configuration loading
│   ├── server.go              # SOCKS5/HTTP server implementation
│   ├── policy.go              # Traffic policy management
│   ├── dns.go                 # DNS handling
│   ├── fragment.go            # Packet fragmentation
│   ├── desync_*.go            # Platform-specific desync
│   └── utils.go               # Utility functions
├── lumine_cmd/                # CLI wrapper for standalone binary
│   └── main.go                # Entry point for lumine executable
├── lumine.go                  # (Optional) Libcore wrapper for programmatic access
├── build_lumine.sh            # Build script for lumine binaries
└── build.sh                   # Main build script (updated to build lumine)
```

## Build Process

### Building Lumine

The `build_lumine.sh` script:
1. Sets up Android NDK cross-compilation environment
2. Builds lumine for multiple Android architectures:
   - ARM64 (arm64-v8a)
   - ARMv7 (armeabi-v7a)
   - x86
   - x86_64
3. Outputs `liblumine.so` binaries to `app/src/main/jniLibs/`
4. Each binary is a standalone executable (named .so for Android compatibility)

### Build Commands

```bash
cd libcore

# Build only lumine
./build_lumine.sh

# Build entire libcore (includes lumine)
./build.sh
```

## Plugin Registration

### In PluginManager.kt

Lumine is registered as an internal plugin in `initNativeInternal()`:

```kotlin
private fun initNativeInternal(pluginId: String): String? {
    fun soIfExist(soName: String): String? {
        val f = File(SagerNet.application.applicationInfo.nativeLibraryDir, soName)
        if (f.canExecute()) {
            return f.absolutePath
        }
        return null
    }
    return when (pluginId) {
        "hysteria-plugin" -> soIfExist("libhysteria.so")
        "hysteria2-plugin" -> soIfExist("libhysteria2.so")
        "lumine-plugin" -> soIfExist("liblumine.so")  // <-- Added
        else -> null
    }
}
```

### In Executable.kt

Lumine is added to the list of managed executables:

```kotlin
private val EXECUTABLES = setOf(
    "libtrojan.so", "libtrojan-go.so", "libnaive.so", 
    "libtuic.so", "libhysteria.so", "liblumine.so"  // <-- Added
)
```

## Usage

### Configuration

Lumine requires a JSON configuration file. Example structure:

```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "127.0.0.1:1081",
  "dns_addr": "8.8.8.8:53",
  "default_policy": {
    "mode": "proxy",
    "fragment_size": 1024
  },
  "domain_policies": {
    "example.com": {
      "mode": "direct"
    }
  }
}
```

### Running Lumine

From Android app:

```kotlin
// Get plugin path
val pluginPath = PluginManager.init("lumine-plugin")?.path

// Execute with config
ProcessBuilder(pluginPath, "-c", configPath).start()
```

Command-line options:
- `-c <path>`: Configuration file path (default: config.json)
- `-b <addr>`: SOCKS5 bind address (overrides config)
- `-hb <addr>`: HTTP bind address (overrides config)

## Architecture

### Package Structure

- **`lumine` package**: Core library code
  - Provides reusable functions for SOCKS5/HTTP server
  - Policy management and traffic manipulation
  - DNS handling and caching
  
- **`lumine_cmd` package**: CLI wrapper
  - Simple main() that calls lumine library functions
  - Handles command-line arguments
  - Manages server lifecycle

### Process Flow

1. App calls PluginManager.init("lumine-plugin")
2. PluginManager locates liblumine.so in native library directory
3. App launches liblumine.so with configuration
4. Lumine starts SOCKS5 and/or HTTP proxy servers
5. Traffic is routed through lumine for preprocessing
6. Lumine applies policies, fragmentation, DNS manipulation
7. Processed traffic is forwarded to destination

## Dependencies

Lumine requires the following Go packages:
- `github.com/miekg/dns` - DNS library
- `github.com/moi-si/addrtrie` - IP/domain matching
- `golang.org/x/net` - Network utilities
- `golang.org/x/sys` - System calls

These are added to libcore's `go.mod` file.

## Future Enhancements

Possible improvements:
1. **JNI Integration**: Convert to gomobile-bound library for tighter integration
2. **Configuration UI**: Add Android UI for lumine configuration
3. **Protocol Support**: Extend to support additional protocols (VMess, Trojan, etc.)
4. **Performance Tuning**: Optimize for mobile battery and performance
5. **Logging**: Integrate with Android logging system

## Troubleshooting

### Build Failures

- **NDK not found**: Ensure `ANDROID_NDK_HOME` or `ANDROID_HOME` is set
- **Go compilation errors**: Run `go mod tidy` in libcore directory
- **Missing dependencies**: Check that submodules are initialized

### Runtime Issues

- **Plugin not found**: Verify `liblumine.so` exists in app's native library directory
- **Configuration errors**: Check JSON configuration file syntax
- **Permission denied**: Ensure binary has execute permissions

### Debugging

Enable verbose logging by modifying lumine source:
```go
// In server.go
fmt.Println("Verbose log message")
```

Check Android logcat:
```bash
adb logcat | grep lumine
```

## License

Lumine is integrated from [moi-si/lumine](https://github.com/moi-si/lumine) and maintains its original license.
