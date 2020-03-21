package msgs

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	BROADCAST_ADDR = "255.255.255.255"
	UDP_TIMEOUT    = 50 // ms
)

// UDP Ports
const (
	COST_MSG_PORT              = "3000"
	ORDER_MSG_PORT             = "3001"
	ORDER_TENSOR_DIFF_MSG_PORT = "3002"
	TEST_MSG_PORT              = "15000"
)

type messager interface {
	Send()
	Listen() error
	port() string
}

func InitMessages() error {
	gob.Register(&CostMsg{})
	gob.Register(&OrderMsg{})
	gob.Register(&OrderTensorDiffMsg{})
	gob.Register(&TestMsg{})
	return nil
}

func send(msg messager) {
	go func() {
		destinationAddress, err := net.ResolveUDPAddr("udp", BROADCAST_ADDR+":"+msg.port())
		if err != nil {
			log.Fatal(err)
		}

		connection, err := net.DialUDP("udp", nil, destinationAddress)
		if err != nil {
			log.Fatal(err)
		}

		defer connection.Close()
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		err = encoder.Encode(&msg)
		if err != nil {
			log.Fatalf("derp, %s\n", err)
		}
		connection.Write(buffer.Bytes())
		buffer.Reset()
	}()
}

func listen(msg messager) (messager, error) {
	localAddress, _ := net.ResolveUDPAddr("udp", ":"+msg.port())
	connection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		return msg, err
	}
	defer connection.Close()
	inputBytes := make([]byte, 4096)
	fmt.Println("Listening...")
	result := make(chan error)
	timer := time.NewTimer(UDP_TIMEOUT * time.Millisecond)
	var length int
	go func() {
		length, _, err = connection.ReadFromUDP(inputBytes)
		result <- err
	}()
	go func() {
		<-timer.C
		result <- errors.New("timeout during listening")
	}()
	if err = <-result; err != nil {
		return msg, err
	}
	buffer := bytes.NewBuffer(inputBytes[:length])
	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(&msg)
	if err != nil {
		return msg, err
	}
	return msg, nil
}
