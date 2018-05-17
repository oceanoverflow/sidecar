package utils

import (
	"log"
	"net"
)

// GetHostIP get the local ip address
func GetHostIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalln("Oops: " + err.Error())
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
