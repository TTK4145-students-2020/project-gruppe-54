package msgs

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	c "../../configuration"
)

// Needs a local channel to access metadata
var metaDataChanLocal <-chan c.MetaData

// Channels for testing purposes
var testMetaDataSenderChanLocal <-chan c.MetaData
var testMetaDataReceiverChanLocal <-chan c.MetaData

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
	ACK_MSG_PORT               = "15001"
)

type messager interface {
	Send()
	Listen() error
	port() string
}

func InitTestMessage(testMetaDataSenderChan <-chan c.MetaData, testMetaDataReceiverChan <-chan c.MetaData) error {
	gob.Register(&TestMsg{})
	testMetaDataSenderChanLocal = testMetaDataSenderChan
	testMetaDataReceiverChanLocal = testMetaDataReceiverChan
	return nil
}

func InitMessages(metaDataChan <-chan c.MetaData) error {
	gob.Register(&CostMsg{})
	gob.Register(&OrderMsg{})
	gob.Register(&OrderTensorDiffMsg{})
	metaDataChanLocal = metaDataChan
	return nil
}

func sendTest(msg TestMsg) {
	go func() {
		destinationAddress, err := net.ResolveUDPAddr("udp", BROADCAST_ADDR+":"+msg.port())
		if err != nil {
			log.Fatal(err)
		}
		metaData := <-testMetaDataSenderChanLocal
		msg.setId(metaData.Id)

		connection, err := net.DialUDP("udp", nil, destinationAddress)
		if err != nil {
			log.Fatal(err)
		}

		defer connection.Close()
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		err = encoder.Encode(&msg)
		if err != nil {
			// TODO: Have better error handling here?
			log.Fatalf("derp, %s\n", err)
		}
		connection.Write(buffer.Bytes())
		buffer.Reset()
	}()
}

func listenTest(msg TestMsg) (TestMsg, error) {
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

func send(msg messager) {
	go func() {
		destinationAddress, err := net.ResolveUDPAddr("udp", BROADCAST_ADDR+":"+msg.port())
		if err != nil {
			log.Fatal(err)
		}
		// metaData := <-metaDataChanLocal
		// msg.setId(metaData.Id)

		connection, err := net.DialUDP("udp", nil, destinationAddress)
		if err != nil {
			log.Fatal(err)
		}
		defer connection.Close()
		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		err = encoder.Encode(&msg)
		if err != nil {
			// TODO: Have better error handling here?
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
	var length int
	connection.SetReadDeadline(time.Now().Add(time.Millisecond * UDP_TIMEOUT))
	length, _, err = connection.ReadFromUDP(inputBytes)
	if err != nil {
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
