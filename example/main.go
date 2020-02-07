package main

import (
	"./client"
	"./server"
	"log"
	"net"
)

func main() {
	client.Test()
	srv, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Close()
	go client.ClientProxy()
	for {
		conn, err := srv.Accept()
		defer conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		go server.ServerHandle(&conn)
	}
}
