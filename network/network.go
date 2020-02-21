package network

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"net"
)

const broadcastAddr = "255.255.255.255"

func InitSender(msg interface{}, port string) interface{} {
	msg_ch, err := makeChan(msg)
	if err != nil {
		log.Fatal(err)
	}
	go sendMsgWorker(msg_ch, port)
	return msg_ch
}

func makeChan(msg interface{}) (interface{}, error) {
	switch msg.(type) {
	case CostMsg:
		return make(chan CostMsg), nil
	default:
		return nil, errors.New("Could not make message channel!")
	}
}

func sendMsgWorker(msgCh interface{}, addr string) {
	destinationAddress, err := net.ResolveUDPAddr("udp", broadcastAddr+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	connection, err := net.DialUDP("udp", nil, destinationAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	for {
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		message := <-msgCh
		err := encoder.Encode(&message)
		if err != nil {
			log.Fatal(err)
		}
		connection.Write(buffer.Bytes())
		buffer.Reset()
	}
}
