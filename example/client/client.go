package client

import (
    "../msgs"
    "fmt"
    "net"
    "log"
    "time"
    "encoding/gob"
)

func Test() {
    fmt.Println("Testing!")
}

func ClientProxy() {
	a, b := 1, 2
	for {
		conn, err := net.Dial("tcp", "localhost:3000")
		if err != nil {
			log.Fatal(err)
		}
		enc := gob.NewEncoder(conn)
		msg := &msgs.Msg{a, b}
		fmt.Printf("Sent: %+v\n", msg)
		err = enc.Encode(msg)
		if err != nil {
			log.Fatal(err)
		}
		conn.Close()
		a, b = b, a
		time.Sleep(500 * time.Millisecond)
	}
}