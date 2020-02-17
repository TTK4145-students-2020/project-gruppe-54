package main

// import (
// 	"bytes"
// 	"encoding/gob"
// 	"fmt"
// 	"log"
// 	"net"
// 	"time"

// 	"./client"
// 	"./msgs"
// 	"./server"
// )
import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

const broadcast_addr = "255.255.255.255"

type Packet struct {
	ID, Response string
	Content      []byte
	A            int
	B            int
}

func main() {
	rec_ch, send_ch := Init("3000", "3000")
	a, b := 1, 10
	for {
		send_msg := Packet{A: a, B: b}
		fmt.Printf("Sent: %+v\n", send_msg)
		send_ch <- send_msg
		fmt.Printf("Received: %+v\n", <-rec_ch)
		a, b = a+1, b+1
		time.Sleep(500 * time.Millisecond)
	}
}

func Init(readPort string, writePort string) (<-chan Packet, chan<- Packet) {
	receive := make(chan Packet, 10)
	send := make(chan Packet, 10)
	go listen(receive, readPort)
	go broadcast(send, writePort)
	return receive, send
}

func listen(receive chan Packet, port string) {
	localAddress, _ := net.ResolveUDPAddr("udp", ":"+port)
	connection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	var message Packet
	for {
		inputBytes := make([]byte, 4096)
		length, _, _ := connection.ReadFromUDP(inputBytes)
		buffer := bytes.NewBuffer(inputBytes[:length])
		decoder := gob.NewDecoder(buffer)
		decoder.Decode(&message)
		//Filters out all messages not relevant for the system
		receive <- message

	}
}

func broadcast(send chan Packet, port string) {
	destinationAddress, _ := net.ResolveUDPAddr("udp", broadcast_addr+":"+port)
	fmt.Printf("Destination address: %+v\n", destinationAddress)
	connection, err := net.DialUDP("udp", nil, destinationAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	for {
		message := <-send
		encoder.Encode(message)
		connection.Write(buffer.Bytes())
		buffer.Reset()
	}
}

// func UDPexample() {
// 	ln, err := net.ListenPacket("udp", "localhost:3000")
// 	defer ln.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	new_loop := make(chan bool)
// 	go func() {
// 		a, b := 1, 10
// 		for {
// 			conn, err := net.Dial("udp", "localhost:3000")
// 			defer conn.Close()
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			enc := gob.NewEncoder(conn)
// 			send_msg := &msgs.Msg{A: a, B: b}
// 			err = enc.Encode(send_msg)
// 			a, b = a+1, b+1
// 			new_loop <- true
// 			time.Sleep(1000 * time.Millisecond)
// 		}
// 	}()
// 	for {
// 		buf := make([]byte, 1024)
// 		n, _, err := ln.ReadFrom(buf)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		buffer := bytes.NewBuffer(buf[:n])
// 		dec := gob.NewDecoder(buffer)
// 		test_msg := &msgs.Msg{}
// 		err = dec.Decode(test_msg)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("Received: %+v\n", test_msg)
// 		<-new_loop
// 	}
// }

// func TCPexample() {
// 	srv, err := net.Listen("tcp", "localhost:3000")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer srv.Close()
// 	go client.ClientProxy()
// 	for {
// 		conn, err := srv.Accept()
// 		defer conn.Close()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		go server.ServerHandle(&conn)
// 	}
// }
