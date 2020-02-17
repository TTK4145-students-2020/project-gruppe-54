package network

import (
	"encoding/gob"
	"log"
	"net"
)

type Msg struct {
	A, B int
}

func InitSender(addr string) chan interface{} {
	msg_ch := make(chan interface{})
	go sendNewMessage(addr, msg_ch)
	return msg_ch
}

func sendNewMessage(addr string, msg_ch chan interface{}) {
	conn, err := net.Dial("udp", addr)
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	enc := gob.NewEncoder(conn)
	for {
		err = enc.Encode(<-msg_ch)
		if err != nil {
			log.Fatal(err)
		}
	}
}
