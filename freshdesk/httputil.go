package freshdesk

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	connectTimeOut   = time.Duration(10 * time.Second)
	readWriteTimeout = time.Duration(20 * time.Second)
)

func timeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		if rwTimeout > 0 {
			conn.SetDeadline(time.Now().Add(rwTimeout))
		}
		return conn, nil
	}
}

// apps will set two OS variables:
// freshdesk_sslcert - location of the http ssl cert
// freshdesk_sslkey - location of the http ssl key
func NewTimeoutClient(cTimeout time.Duration, rwTimeout time.Duration) *http.Client {
	certLocation := os.Getenv("freshdesk_sslcert")
	keyLocation := os.Getenv("freshdesk_sslkey")
	// default
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	if len(certLocation) > 0 && len(keyLocation) > 0 {
		// Load client cert if available
		cert, err := tls.LoadX509KeyPair(certLocation, keyLocation)
		if err == nil {
			tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		}
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
			Dial:            timeoutDialer(cTimeout, rwTimeout),
		},
	}
}

func DefaultTimeoutClient() *http.Client {
	return NewTimeoutClient(connectTimeOut, readWriteTimeout)
}
