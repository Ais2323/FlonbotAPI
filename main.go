package FlonBotAPI

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

// 主程式
func main() {
	var ipAddress = "127.0.0.1"
	var postAddress = 8080
	var fullIPAddr = fmt.Sprintf("%s:%d", ipAddress, postAddress)
	var addr = fmt.Sprintf(":%d", postAddress)
	fmt.Printf("Local IP Address: %s\n", fullIPAddr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}
