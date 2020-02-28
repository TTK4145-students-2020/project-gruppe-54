package network

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
)

func makeChan(msg interface{}) (MsgCh, error) {
	switch msg.(type) {
	case CostMsg:
		gob.Register(CostMsg{})
		return CostMsgCh{make(chan interface{})}, nil
	case TestMsg:
		gob.Register(TestMsg{})
		return TestMsgCh{make(chan interface{})}, nil
	default:
		return nil, errors.New("Could not make message channel!")
	}
}

func getPort(msgCh MsgCh) (string, error) {
	switch msgCh.(type) {
	case CostMsgCh:
		return COST_MSG_PORT, nil
	case TestMsgCh:
		return TEST_MSG_PORT, nil
	default:
		return "", errors.New("unknown channel type, cannot get port")
	}
}

// func InitSender(msg interface{}) MsgCh {
// 	msg_ch, err := makeChan(msg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	go sendMsgWorker(msg_ch)
// 	return msg_ch
// }

type testMsgTransmitter struct {
	Channel chan TestMsg
}

type transmitter interface {
	port() string
	channel() chan networkMsg
}

func (trans testMsgTransmitter) port() string {
	return TEST_MSG_PORT
}

func (trans testMsgTransmitter) channel() chan networkMsg {
	return trans.Channel
}

func NewTestMsgTransmitter(testMsgCh chan TestMsg) testMsgTransmitter {
	trans := testMsgTransmitter{testMsgCh}
	go msgWorker(trans)
	return trans
}

func msgWorker(trans transmitter) {
	destinationAddress, err := net.ResolveUDPAddr("udp", BROADCAST_ADDR+":"+trans.port())
	fmt.Println("Sending on", trans.port())
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
		fmt.Println("Waiting for message to send...")
		message := <-trans.Channel
		err := encoder.Encode(&message)
		fmt.Println("Sent message!")
		if err != nil {
			log.Fatalf("derp, %s\n", err)
		}
		connection.Write(buffer.Bytes())
		buffer.Reset()
	}
}

func sendMsgWorker(msgCh MsgCh) {
	port, err := getPort(msgCh)
	if err != nil {
		log.Fatal(err)
	}

	destinationAddress, err := net.ResolveUDPAddr("udp", BROADCAST_ADDR+":"+port)
	fmt.Println("Sending on", port)
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
		fmt.Println("Waiting for message to send...")
		message := <-msgCh.Ch()
		err := encoder.Encode(&message)
		fmt.Println("Sent message!")
		if err != nil {
			log.Fatalf("derp, %s\n", err)
		}
		connection.Write(buffer.Bytes())
		buffer.Reset()
	}
}
