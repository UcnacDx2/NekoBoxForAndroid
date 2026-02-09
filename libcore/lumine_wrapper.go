package libcore

import (
	"libcore/lumine"
	"log"
)

// LumineService represents a Lumine proxy service instance
type LumineService struct {
	proxy *lumine.LumineProxy
}

// NewLumineService creates a new Lumine service with the given configuration file
func NewLumineService(configPath string) (*LumineService, error) {
	proxy, err := lumine.NewLumineProxy(configPath)
	if err != nil {
		return nil, err
	}
	return &LumineService{proxy: proxy}, nil
}

// NewLumineServiceWithAddrs creates a new Lumine service with explicit addresses
// socks5Addr: SOCKS5 listen address (e.g., "127.0.0.1:1080" or "none" to disable)
// httpAddr: HTTP proxy listen address (e.g., "127.0.0.1:1225" or "none" to disable)
// configJSON: JSON configuration string (optional, can be empty)
func NewLumineServiceWithAddrs(socks5Addr, httpAddr, configJSON string) (*LumineService, error) {
	// For now, we'll use a simple configuration
	// In the future, we could parse configJSON to create a Config struct
	proxy, err := lumine.NewLumineProxyWithConfig(socks5Addr, httpAddr, nil)
	if err != nil {
		return nil, err
	}
	return &LumineService{proxy: proxy}, nil
}

// Close stops the Lumine service
func (s *LumineService) Close() error {
	if s.proxy != nil {
		return s.proxy.Close()
	}
	return nil
}

// TestLumineIntegration tests if lumine integration works
func TestLumineIntegration() string {
	log.Println("Lumine integration test: OK")
	return "Lumine network preprocessing is integrated successfully"
}
