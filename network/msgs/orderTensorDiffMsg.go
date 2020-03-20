package msgs

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"../../hardware/driver-go/elevio"
)

const (
	ORDER_TENSOR_DIFF_MSG_PORT = "3002"
)

type OrderTensorDiff int

const (
	ADD    OrderTensorDiff = 0
	REMOVE                 = 1
)

type OrderTensorDiffMsg struct {
	Order elevio.ButtonEvent
	Diff  OrderTensorDiff
}

func (msg OrderTensorDiffMsg) port() string {
	return ORDER_TENSOR_DIFF_MSG_PORT
}

func (msg OrderTensorDiffMsg) Send() {
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

func (msg *OrderTensorDiffMsg) Listen() error {
	localAddress, _ := net.ResolveUDPAddr("udp", ":"+msg.port())
	connection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		return err
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
		return err
	}
	buffer := bytes.NewBuffer(inputBytes[:length])
	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(&msg)
	if err != nil {
		return err
	}
	return nil
}