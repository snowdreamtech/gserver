package net

import (
	"net"
	"strconv"

	"github.com/snowdreamtech/gserver/pkg/tools"
)

// IsPortAvailable Is Port Available
func IsPortAvailable(port string) bool {
	listener, err := net.Listen("tcp", ":"+port)

	if err == nil {
		listener.Close()
		return true
	}

	return false
}

// GetAvailablePort Get Available Port
func GetAvailablePort(startPort int) string {
	port := ""

	// Get the free port from startPort
	for i := startPort; i < 65535; i++ {
		listener, err := net.Listen("tcp", ":"+strconv.Itoa(i))

		if err == nil {
			listener.Close()
			port = strconv.Itoa(i)
			break
		}
	}

	return port
}

// GetAvailableIPS Get Available IP
func GetAvailableIPS() []string {
	ips := []string{}

	// get list of available addresses
	addr, err := net.InterfaceAddrs()
	if err != nil {
		tools.DebugPrintF("[ERROR] " + err.Error())
		return ips
	}

	for _, addr := range addr {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// check if IPv4 or IPv6 is not nil
			if ipnet.IP.To4() != nil || ipnet.IP.To16() != nil {
				// print available addresses
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	return ips
}
