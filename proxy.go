package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func RunProxy(cfg *Config) error {
	server := &http.Server{
		Addr: cfg.ListenAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			deny := false
			if len(cfg.Users) > 0 {
				user, pass, ok := parseProxyBasicAuth(r)
				if ok {
					for _, u := range cfg.Users {
						if u.UserName == user && u.Password == pass {
							deny = true
							break
						}
					}
				}
			} else {
				deny = true
			}
			if !deny {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			r.Header.Del("Proxy-Authorization")
			if r.Method == http.MethodConnect {
				handleTunneling(w, r, cfg.BindIp)
			} else {
				handleHTTP(w, r, cfg.BindIp)
			}
		}),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Printf("Listen: %s", cfg.ListenAddr)
	if cfg.BindIp != "" {
		log.Printf("BindIp: %s", cfg.BindIp)
	} else {
		log.Println("Bind default ip")
	}
	return server.ListenAndServe()
}

func handleTunneling(w http.ResponseWriter, r *http.Request, bindIp string) {
	dialler := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 90 * time.Second,
	}
	if bindIp != "" {
		localAddr, err := net.ResolveIPAddr("ip", bindIp)
		if err != nil {
			log.Fatal(fmt.Errorf("incorrect localAddr %s %s", bindIp, err))
		}
		dialler.LocalAddr = &net.TCPAddr{
			IP: localAddr.IP,
		}
	}
	destConn, err := dialler.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}
func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
func handleHTTP(w http.ResponseWriter, req *http.Request, bindIp string) {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	if bindIp != "" {
		localAddr, err := net.ResolveIPAddr("ip", bindIp)
		if err != nil {
			log.Fatal(fmt.Errorf("incorrect localAddr %s %s", bindIp, err))
		}
		dialer.LocalAddr = &net.TCPAddr{
			IP: localAddr.IP,
		}
	}

	def := http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	resp, err := def.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
func parseProxyBasicAuth(r *http.Request) (username, password string, ok bool) {
	auth := r.Header.Get("Proxy-Authorization")
	if auth == "" {
		return
	}
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}
