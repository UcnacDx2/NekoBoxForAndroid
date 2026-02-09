package lumine

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"slices"
	"sync"
)

// LumineProxy represents a running Lumine proxy instance
type LumineProxy struct {
	socks5Listener net.Listener
	httpListener   net.Listener
	cancel         chan struct{}
	wg             sync.WaitGroup
}

// NewLumineProxy creates a new Lumine proxy instance with the given configuration
func NewLumineProxy(configPath string) (*LumineProxy, error) {
	socks5Addr, httpAddr, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	proxy := &LumineProxy{
		cancel: make(chan struct{}),
	}

	// Start SOCKS5 listener if configured
	if socks5Addr != "" && socks5Addr != "none" {
		ln, err := net.Listen("tcp", socks5Addr)
		if err != nil {
			return nil, fmt.Errorf("SOCKS5 listen error: %v", err)
		}
		proxy.socks5Listener = ln
		proxy.wg.Add(1)
		go proxy.serveSocks5(ln)
	}

	// Start HTTP listener if configured
	if httpAddr != "" && httpAddr != "none" {
		ln, err := net.Listen("tcp", httpAddr)
		if err != nil {
			if proxy.socks5Listener != nil {
				proxy.socks5Listener.Close()
			}
			return nil, fmt.Errorf("HTTP listen error: %v", err)
		}
		proxy.httpListener = ln
		proxy.wg.Add(1)
		go proxy.serveHTTP(ln)
	}

	return proxy, nil
}

// NewLumineProxyWithConfig creates a new Lumine proxy with explicit addresses
func NewLumineProxyWithConfig(socks5Addr, httpAddr string, config *Config) (*LumineProxy, error) {
	if config != nil {
		if err := applyConfig(config); err != nil {
			return nil, fmt.Errorf("failed to apply config: %v", err)
		}
	}

	proxy := &LumineProxy{
		cancel: make(chan struct{}),
	}

	// Start SOCKS5 listener if configured
	if socks5Addr != "" && socks5Addr != "none" {
		ln, err := net.Listen("tcp", socks5Addr)
		if err != nil {
			return nil, fmt.Errorf("SOCKS5 listen error: %v", err)
		}
		proxy.socks5Listener = ln
		proxy.wg.Add(1)
		go proxy.serveSocks5(ln)
	}

	// Start HTTP listener if configured
	if httpAddr != "" && httpAddr != "none" {
		ln, err := net.Listen("tcp", httpAddr)
		if err != nil {
			if proxy.socks5Listener != nil {
				proxy.socks5Listener.Close()
			}
			return nil, fmt.Errorf("HTTP listen error: %v", err)
		}
		proxy.httpListener = ln
		proxy.wg.Add(1)
		go proxy.serveHTTP(ln)
	}

	return proxy, nil
}

// Close stops the Lumine proxy
func (p *LumineProxy) Close() error {
	close(p.cancel)
	
	if p.socks5Listener != nil {
		p.socks5Listener.Close()
	}
	if p.httpListener != nil {
		p.httpListener.Close()
	}
	
	p.wg.Wait()
	return nil
}

func (p *LumineProxy) serveSocks5(ln net.Listener) {
	defer p.wg.Done()
	
	var connID uint32
	for {
		select {
		case <-p.cancel:
			return
		default:
		}

		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-p.cancel:
				return
			default:
				log.Printf("SOCKS5 accept error: %s", err)
				continue
			}
		}
		
		connID++
		if connID > 0xFFFFF {
			connID = 0
		}
		go HandleSOCKS5(conn, connID)
	}
}

func (p *LumineProxy) serveHTTP(ln net.Listener) {
	defer p.wg.Done()
	
	var connID uint32
	for {
		select {
		case <-p.cancel:
			return
		default:
		}

		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-p.cancel:
				return
			default:
				log.Printf("HTTP accept error: %s", err)
				continue
			}
		}
		
		connID++
		if connID > 0xFFFFF {
			connID = 0
		}
		go HandleHTTP(conn, connID)
	}
}

func readN(conn net.Conn, n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := io.ReadFull(conn, buf)
	return buf, err
}

func sendReply(logger *log.Logger, conn net.Conn, rep byte) {
	resp := []byte{0x05, rep, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	if _, err := conn.Write(resp); err != nil {
		logger.Println("Send SOCKS5 reply fail:", err)
	}
}

// HandleSOCKS5 handles a SOCKS5 connection
func HandleSOCKS5(clientConn net.Conn, id uint32) {
	var (
		once    sync.Once
		dstConn net.Conn
	)
	closeBoth := func() {
		once.Do(func() {
			clientConn.Close()
			if dstConn != nil {
				dstConn.Close()
			}
		})
	}
	defer closeBoth()

	logger := log.New(os.Stdout, fmt.Sprintf("[S%05x] ", id), log.LstdFlags)
	logger.Println("Connection from", clientConn.RemoteAddr().String())

	header, err := readN(clientConn, 2)
	if err != nil {
		logger.Println("Read method selection fail:", err)
		return
	}
	if header[0] != 0x05 {
		logger.Println("Not SOCKS5:", header[0])
		return
	}
	nMethods := int(header[1])
	methods, err := readN(clientConn, nMethods)
	if err != nil {
		logger.Println("Read methods fail:", err)
		return
	}
	var authMethod byte = 0xFF
	if slices.Contains(methods, 0x00) {
		authMethod = 0x00
	}
	if _, err = clientConn.Write([]byte{0x05, authMethod}); err != nil {
		logger.Println("Method write fail:", err)
		return
	}
	if authMethod == 0xFF {
		logger.Println("No `no auth` method")
		return
	}

	header, err = readN(clientConn, 4)
	if err != nil {
		logger.Println("Read req header fail:", err)
		return
	}
	if header[0] != 0x05 {
		logger.Println("Ver err:", header[0])
		return
	}
	if header[1] != 0x01 {
		logger.Println("Not CONNECT:", header[1])
		sendReply(logger, clientConn, 0x07)
		return
	}

	var (
		originHost, dstHost string
		policy              *Policy
	)
	switch header[3] {
	case 0x01: // IPv4 address
		ipBytes, err := readN(clientConn, 4)
		if err != nil {
			logger.Println("Read IPv4 fail:", err)
			return
		}
		originHost = net.IP(ipBytes).String()
		var ipPolicy *Policy
		dstHost, ipPolicy, err = IpRedirect(logger, originHost)
		if err != nil {
			logger.Println("IP redirect error:", err)
			sendReply(logger, clientConn, 0x01)
			return
		}
		if ipPolicy == nil {
			policy = GetDefaultPolicy()
		} else {
			policy = MergePolicies(*ipPolicy, *GetDefaultPolicy())
		}
	case 0x04: // IPv6 address
		ipBytes, err := readN(clientConn, 16)
		if err != nil {
			logger.Println("Read IPv6 fail", err)
			return
		}
		originHost = net.IP(ipBytes).String()
		var ipPolicy *Policy
		dstHost, ipPolicy, err = IpRedirect(logger, originHost)
		if err != nil {
			logger.Println("IP redirect error:", err)
			sendReply(logger, clientConn, 0x01)
			return
		}
		if ipPolicy == nil {
			policy = GetDefaultPolicy()
		} else {
			policy = MergePolicies(*ipPolicy, *GetDefaultPolicy())
		}
	case 0x03: // Domain name
		lenByte, err := readN(clientConn, 1)
		if err != nil {
			logger.Println("Read domain len fail:", err)
			return
		}
		domainBytes, err := readN(clientConn, int(lenByte[0]))
		if err != nil {
			logger.Println("Read domain fail:", err)
		}
		originHost = string(domainBytes)
		var fail, block bool
		dstHost, policy, fail, block = GenPolicy(logger, originHost)
		if fail {
			sendReply(logger, clientConn, 0x01)
			return
		}
		if block {
			logger.Printf("Blocked connection to %s", originHost)
			sendReply(logger, clientConn, 0x02)
			return
		}
	default:
		logger.Println("Invalid address type:", header[3])
		sendReply(logger, clientConn, 0x08)
		return
	}
	portBytes, err := readN(clientConn, 2)
	if err != nil {
		logger.Println("Read port fail:", err)
		return
	}
	dstPort := binary.BigEndian.Uint16(portBytes)
	oldTarget := net.JoinHostPort(originHost, fmt.Sprintf("%d", dstPort))
	logger.Println("CONNECT", oldTarget)
	logger.Println("Policy:", policy)
	if policy.Mode == ModeBlock {
		sendReply(logger, clientConn, 0x02)
		return
	}
	if policy.Port != 0 && policy.Port != -1 {
		dstPort = uint16(policy.Port)
	}
	target := net.JoinHostPort(dstHost, fmt.Sprintf("%d", dstPort))

	replyFirst := policy.ReplyFirst == BoolTrue
	if !replyFirst {
		dstConn, err = net.DialTimeout("tcp", target, policy.ConnectTimeout)
		if err != nil {
			logger.Println("Connection failed:", err)
			sendReply(logger, clientConn, 0x01)
			return
		}
	}
	sendReply(logger, clientConn, 0x00)

	HandleTunnel(policy, replyFirst, dstConn, clientConn,
		logger, target, originHost, closeBoth)
}
