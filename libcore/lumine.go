package libcore

import (
	"fmt"
	"libcore/lumine"
	"log"
	"os"
)

// LumineServer represents a running lumine server instance
type LumineServer struct {
	configPath  string
	socks5Addr  string
	httpAddr    string
	done        chan struct{}
	stopChan    chan struct{}
}

// NewLumineServer creates a new lumine server instance
func NewLumineServer(configPath string) (*LumineServer, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config path cannot be empty")
	}
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}
	
	return &LumineServer{
		configPath: configPath,
		done:       make(chan struct{}, 2),
		stopChan:   make(chan struct{}),
	}, nil
}

// Start starts the lumine server
func (s *LumineServer) Start() error {
	socks5Addr, httpAddr, err := lumine.LoadConfig(s.configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	s.socks5Addr = socks5Addr
	s.httpAddr = httpAddr
	
	log.Printf("Starting lumine server with SOCKS5: %s, HTTP: %s", socks5Addr, httpAddr)
	
	// Start SOCKS5 and HTTP servers
	go s.startSOCKS5(socks5Addr)
	go s.startHTTP(httpAddr)
	
	return nil
}

// Stop stops the lumine server
func (s *LumineServer) Stop() {
	close(s.stopChan)
	// Wait for both servers to stop
	<-s.done
	<-s.done
}

func (s *LumineServer) startSOCKS5(addr string) {
	defer func() { s.done <- struct{}{} }()
	
	if addr == "" || addr == "none" {
		log.Println("SOCKS5 server disabled")
		return
	}
	
	// Use lumine's internal function to start SOCKS5 server
	lumine.StartSOCKS5Server(addr, s.stopChan)
}

func (s *LumineServer) startHTTP(addr string) {
	defer func() { s.done <- struct{}{} }()
	
	if addr == "" || addr == "none" {
		log.Println("HTTP server disabled")
		return
	}
	
	// Use lumine's internal function to start HTTP server
	lumine.StartHTTPServer(addr, s.stopChan)
}

// GetSOCKS5Address returns the SOCKS5 server address
func (s *LumineServer) GetSOCKS5Address() string {
	return s.socks5Addr
}

// GetHTTPAddress returns the HTTP server address
func (s *LumineServer) GetHTTPAddress() string {
	return s.httpAddr
}
