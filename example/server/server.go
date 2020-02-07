package server

import (
    "fmt"
    "../msgs"
    "net"
    "log"
    "time"
    "encoding/gob"
)

func ServerHandle(conn *net.Conn) {
	dec := gob.NewDecoder(*conn)
	msg := &msgs.Msg{}
	err := dec.Decode(msg)
	if err != nil {
		log.Fatal(err)
	}
	curTime := time.Now()
	fmt.Printf("%s: %+v\n", curTime.Format("2006-02-05 15:04:05.000000"), msg)
}