package main

import (
"flag"
"fmt"
"libcore/lumine"
"os"
)

func main() {
flag.Usage = func() {
fmt.Fprintln(os.Stderr, "lumine network preprocessing tool v0.3.0")
fmt.Fprintln(os.Stderr)
flag.PrintDefaults()
}
configPath := flag.String("c", "config.json", "Config file path")
addr := flag.String("b", "", "SOCKS5 bind address (default: address from config file)")
hAddr := flag.String("hb", "", "HTTP bind address (default: address from config file)")

flag.Parse()

socks5Addr, httpAddr, err := lumine.LoadConfig(*configPath)
if err != nil {
fmt.Println("Failed to load config:", err)
return
}

if *addr != "" {
socks5Addr = *addr
}
if *hAddr != "" {
httpAddr = *hAddr
}

done := make(chan struct{}, 2)
stopChan := make(chan struct{})

go func() {
lumine.StartSOCKS5Server(socks5Addr, stopChan)
done <- struct{}{}
}()

go func() {
lumine.StartHTTPServer(httpAddr, stopChan)
done <- struct{}{}
}()

<-done
<-done
}
