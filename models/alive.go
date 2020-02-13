package models

import (
	"log"
	"net"
	"os"
	"time"
)

//IsAliveResponse  -
type IsAliveResponse struct {
	Hostname  string `json:"hostname"`
	IP        string `json:"ip"`
	Timestemp int64  `json:"timestemp"`
	VER       string `json:"version"`
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

//GetIsAliveResponse - returns isalive resposne
func GetIsAliveResponse() IsAliveResponse {

	myip := getOutboundIP()
	hname, _ := os.Hostname()

	return IsAliveResponse{hname, myip.String(), time.Now().UnixNano(), "v1.0.6"}
}
