package lumine

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

var httpConnID uint32

// LoadConfig loads the configuration from the specified file
func LoadConfig(configPath string) (string, string, error) {
	return loadConfig(configPath)
}

// StartSOCKS5Server starts the SOCKS5 server on the specified address
func StartSOCKS5Server(serverAddr string, stopChan chan struct{}) {
	if serverAddr == "" {
		fmt.Println("SOCKS5 bind address not specified")
		return
	}
	if serverAddr == "none" {
		return
	}

	ln, err := net.Listen("tcp", serverAddr)
	if err != nil {
		fmt.Println("SOCKS5 Listen error:", err)
		return
	}
	defer ln.Close()

	listenAddr := serverAddr
	if listenAddr[0] == ':' {
		listenAddr = "0.0.0.0" + listenAddr
	}
	fmt.Println("Listening on", "socks5://"+listenAddr)

	var connID uint32
	// Accept loop with stop channel
	acceptDone := make(chan struct{})
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-stopChan:
					// Server is stopping
					return
				default:
					log.Printf("SOCKS5 accept error: %s", err)
					continue
				}
			}
			connID += 1
			if connID > 0xFFFFF {
				connID = 0
			}
			go handleSOCKS5(conn, connID)
		}
	}()

	<-stopChan
	close(acceptDone)
}

// StartHTTPServer starts the HTTP proxy server on the specified address
func StartHTTPServer(serverAddr string, stopChan chan struct{}) {
	if serverAddr == "" {
		fmt.Println("HTTP bind address not specified")
		return
	}
	if serverAddr == "none" {
		return
	}

	srv := &http.Server{
		Addr:              serverAddr,
		Handler:           http.HandlerFunc(handleHTTP),
		ReadHeaderTimeout: 10 * time.Second,
	}

	listenAddr := serverAddr
	if listenAddr[0] == ':' {
		listenAddr = "0.0.0.0" + listenAddr
	}
	fmt.Println("Listening on", "http://"+listenAddr)

	// Start server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("HTTP ListenAndServe:", err)
		}
	}()

	<-stopChan
	srv.Close()
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

func handleSOCKS5(clientConn net.Conn, id uint32) {
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
		dstHost, ipPolicy, err = ipRedirect(logger, originHost)
		if err != nil {
			logger.Println("IP redirect error:", err)
			sendReply(logger, clientConn, 0x01)
			return
		}
		if ipPolicy == nil {
			policy = &defaultPolicy
		} else {
			policy = mergePolicies(*ipPolicy, defaultPolicy)
		}
	case 0x04: // IPv6 address
		ipBytes, err := readN(clientConn, 16)
		if err != nil {
			logger.Println("Read IPv6 fail", err)
			return
		}
		originHost = net.IP(ipBytes).String()
		var ipPolicy *Policy
		dstHost, ipPolicy, err = ipRedirect(logger, originHost)
		if err != nil {
			logger.Println("IP redirect error:", err)
			sendReply(logger, clientConn, 0x01)
			return
		}
		if ipPolicy == nil {
			policy = &defaultPolicy
		} else {
			policy = mergePolicies(*ipPolicy, defaultPolicy)
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
		dstHost, policy, fail, block = genPolicy(logger, originHost)
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

	handleTunnel(policy, replyFirst, dstConn, clientConn,
		logger, target, originHost, closeBoth)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	connID := atomic.AddUint32(&httpConnID, 1)
	if connID > 0xFFFFF {
		atomic.StoreUint32(&httpConnID, 0)
		connID = 0
	}

	logger := log.New(os.Stdout, fmt.Sprintf("[H%05x] ", connID), log.LstdFlags)
	logger.Printf("%s - \"%s %s %s\"", req.RemoteAddr, req.Method, req.RequestURI, req.Proto)

	if req.Method == http.MethodConnect {
		handleConnect(logger, w, req)
		return
	}

	if !req.URL.IsAbs() {
		logger.Println("URI not fully qualified")
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}

	forwardHTTPRequest(logger, w, req)
}

func handleConnect(logger *log.Logger, w http.ResponseWriter, req *http.Request) {
	const (
		status500 = "500 Internal Server Error"
		status403 = "403 Forbidden"
	)

	oldDest := req.Host
	if oldDest == "" {
		logger.Println("Empty host")
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	originHost, dstPort, err := net.SplitHostPort(oldDest)
	if err != nil {
		logger.Println("SplitHostPort fail:", err)
		return
	}

	dstHost, policy, fail, block := genPolicy(logger, originHost)
	if fail {
		http.Error(w, status500, http.StatusInternalServerError)
		return
	}
	if block {
		logger.Println("Connection blocked")
		http.Error(w, status403, http.StatusForbidden)
		return
	}

	logger.Println("Policy:", policy)

	if policy.Mode == ModeBlock {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	if policy.Port != 0 && policy.Port != -1 {
		dstPort = fmt.Sprintf("%d", policy.Port)
	}

	dest := net.JoinHostPort(dstHost, dstPort)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		logger.Println("Hijacking not supported")
		http.Error(w, status500, http.StatusInternalServerError)
		return
	}
	cliConn, _, err := hijacker.Hijack()
	if err != nil {
		logger.Println("Hijack fail:", err)
		http.Error(w, status500, http.StatusInternalServerError)
		return
	}

	var (
		once    sync.Once
		dstConn net.Conn
	)
	closeBoth := func() {
		once.Do(func() {
			cliConn.Close()
			if dstConn != nil {
				dstConn.Close()
			}
		})
	}
	defer closeBoth()

	replyFirst := policy.ReplyFirst == BoolTrue
	if replyFirst {
		_, err = cliConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
		if err != nil {
			logger.Println("Write 200 error:", err)
			return
		}
	} else {
		dstConn, err = net.Dial("tcp", dest)
		if err != nil {
			logger.Println("Connection failed:", err)
			_, err = cliConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
			if err != nil {
				logger.Println("Write 502 error:", err)
			}
			return
		}
		_, err = cliConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
		if err != nil {
			logger.Println("Write 200 error:", err)
			return
		}
	}

	handleTunnel(policy, replyFirst, dstConn, cliConn,
		logger, dest, originHost, closeBoth)
}

func forwardHTTPRequest(logger *log.Logger, w http.ResponseWriter, originReq *http.Request) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	outReq := originReq.Clone(originReq.Context())
	outReq.Host = outReq.URL.Host
	outReq.Header.Del("Proxy-Authorization")
	outReq.Header.Del("Proxy-Connection")
	if outReq.Header.Get("Connection") == "" {
		outReq.Header.Set("Connection", "close")
	}

	resp, err := transport.RoundTrip(outReq)
	if err != nil {
		logger.Println("Transport error:", err)
		http.Error(w, "502 Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy headers
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	if _, err = io.Copy(w, resp.Body); err != nil {
		logger.Println("Error copying response body:", err)
	}
}
