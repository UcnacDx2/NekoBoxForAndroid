package libcore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"libcore/lumine"
	"os"
)

// LumineChainConfig helps configure lumine as a chain proxy
type LumineChainConfig struct {
	lumineService *lumine.LumineProxy
	configPath    string
	socks5Port    int
	httpPort      int
}

// NewLumineChainConfig creates a new chain proxy configuration helper
// socks5Port: port for SOCKS5 proxy (0 to disable, default 1080)
// httpPort: port for HTTP proxy (0 to disable, default 0)
// configPath: path to lumine config file (empty for defaults)
func NewLumineChainConfig(socks5Port, httpPort int, configPath string) (*LumineChainConfig, error) {
	if socks5Port == 0 && httpPort == 0 {
		return nil, fmt.Errorf("at least one of SOCKS5 or HTTP port must be specified")
	}

	config := &LumineChainConfig{
		configPath: configPath,
		socks5Port: socks5Port,
		httpPort:   httpPort,
	}

	return config, nil
}

// Start starts the lumine proxy service
func (c *LumineChainConfig) Start() error {
	socks5Addr := "none"
	if c.socks5Port > 0 {
		socks5Addr = fmt.Sprintf("127.0.0.1:%d", c.socks5Port)
	}

	httpAddr := "none"
	if c.httpPort > 0 {
		httpAddr = fmt.Sprintf("127.0.0.1:%d", c.httpPort)
	}

	var err error
	if c.configPath != "" {
		c.lumineService, err = lumine.NewLumineProxy(c.configPath)
	} else {
		c.lumineService, err = lumine.NewLumineProxyWithConfig(socks5Addr, httpAddr, nil)
	}

	return err
}

// Stop stops the lumine proxy service
func (c *LumineChainConfig) Stop() error {
	if c.lumineService != nil {
		return c.lumineService.Close()
	}
	return nil
}

// IsRunning returns true if lumine service is running
func (c *LumineChainConfig) IsRunning() bool {
	return c.lumineService != nil
}

// GetSOCKS5Address returns the SOCKS5 proxy address (e.g., "127.0.0.1:1080")
func (c *LumineChainConfig) GetSOCKS5Address() string {
	if c.socks5Port > 0 {
		return fmt.Sprintf("127.0.0.1:%d", c.socks5Port)
	}
	return ""
}

// GetHTTPAddress returns the HTTP proxy address (e.g., "127.0.0.1:1225")
func (c *LumineChainConfig) GetHTTPAddress() string {
	if c.httpPort > 0 {
		return fmt.Sprintf("127.0.0.1:%d", c.httpPort)
	}
	return ""
}

// GenerateDefaultConfig generates a default lumine configuration file
func GenerateDefaultLumineConfig(outputPath string) error {
	defaultConfig := map[string]interface{}{
		"socks5_address":       "127.0.0.1:1080",
		"http_address":         "none",
		"dns_addr":             "8.8.8.8:53",
		"udp_minsize":          4096,
		"max_jump":             30,
		"transmit_file_limit":  2,
		"dns_cache_ttl":        259200,
		"fake_ttl_cache_ttl":   259200,
		"default_policy": map[string]interface{}{
			"connect_timeout": "10s",
			"mode":            "tls-rf",
			"num_records":     10,
			"num_segs":        3,
			"oob":             false,
			"send_interval":   "200ms",
		},
		"domain_policies": map[string]interface{}{},
		"ip_policies":     map[string]interface{}{},
	}

	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputPath, data, 0644)
}

// ValidateLumineConfig validates a lumine configuration file
func ValidateLumineConfig(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", configPath)
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("invalid JSON format: %v", err)
	}

	// Basic validation
	if socks5Addr, ok := config["socks5_address"].(string); ok && socks5Addr == "" {
		return fmt.Errorf("socks5_address cannot be empty")
	}

	return nil
}

// CreateChainProxyOutbound creates a sing-box outbound configuration for using lumine as detour
// Example usage:
//   lumineOutbound := CreateLumineSOCKS5Outbound("lumine", "127.0.0.1", 1080)
//   proxyOutbound := CreateChainProxyOutbound("proxy", "shadowsocks", proxyConfig, "lumine")
func CreateLumineSOCKS5Outbound(tag, server string, port int) string {
	config := map[string]interface{}{
		"type":        "socks",
		"tag":         tag,
		"server":      server,
		"server_port": port,
		"version":     "5",
	}

	data, _ := json.Marshal(config)
	return string(data)
}

// AddLumineDetourToOutbound adds lumine as a detour to an existing outbound configuration
func AddLumineDetourToOutbound(outboundJSON, lumineTag string) (string, error) {
	var outbound map[string]interface{}
	if err := json.Unmarshal([]byte(outboundJSON), &outbound); err != nil {
		return "", fmt.Errorf("invalid outbound JSON: %v", err)
	}

	outbound["detour"] = lumineTag

	data, err := json.Marshal(outbound)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Example: Complete chain proxy setup
type ChainProxySetup struct {
	LumineConfig *LumineChainConfig
	ConfigPath   string
}

// NewChainProxySetup creates a complete chain proxy setup
func NewChainProxySetup(lumineConfigPath string) (*ChainProxySetup, error) {
	// Default to SOCKS5 on port 1080
	lumineConfig, err := NewLumineChainConfig(1080, 0, lumineConfigPath)
	if err != nil {
		return nil, err
	}

	return &ChainProxySetup{
		LumineConfig: lumineConfig,
		ConfigPath:   lumineConfigPath,
	}, nil
}

// Start starts the chain proxy
func (s *ChainProxySetup) Start() error {
	return s.LumineConfig.Start()
}

// Stop stops the chain proxy
func (s *ChainProxySetup) Stop() error {
	return s.LumineConfig.Stop()
}

// GetProxyAddress returns the address to use in sing-box config
func (s *ChainProxySetup) GetProxyAddress() string {
	return s.LumineConfig.GetSOCKS5Address()
}
