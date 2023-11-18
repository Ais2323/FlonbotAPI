package main

import (
	Method "FlonBotApi/APIMethod"
	db "FlonBotApi/Database"
	komica "FlonBotApi/KomicaReply"
	"fmt"
	"log"
	"net/http"
)

// 主程式
func main() {
	db.Init()
	komica.StartRequest()
	InitRequest()

	var ipAddress = "127.0.0.1"
	var postAddress = 80 // http:80
	var fullIPAddr = fmt.Sprintf("%s:%d", ipAddress, postAddress)
	var addr = fmt.Sprintf(":%d", postAddress)
	fmt.Printf("Local IP Address: %s\n", fullIPAddr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func InitRequest() {
	http.HandleFunc("/reply", Method.RequestReply)
}
