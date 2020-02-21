package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestInitSender(t *testing.T) {
	test_msg := TestMsg{}
	msg_ch := InitSender(test_msg)
	localAddress, _ := net.ResolveUDPAddr("udp", ":"+TEST_MSG_PORT)
	connection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	var rec_msg interface{}
	inputBytes := make([]byte, 4096)
	go func() {
		a := 19
		for {
			test_msg = TestMsg{a}
			msg_ch.Ch() <- test_msg
			time.Sleep(500 * time.Millisecond)
			a = a + 1
		}
	}()
	length, _, _ := connection.ReadFromUDP(inputBytes)
	fmt.Println("Received packet!")
	buffer := bytes.NewBuffer(inputBytes[:length])
	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(&rec_msg)
	if err != nil {
		log.Fatalf("decoding: %s", err)
	}
	if !cmp.Equal(rec_msg, test_msg) {
		t.Errorf("rec_msg not equal to test_msg\nrec_msg: %+v\ntest_msg: %+v\n", rec_msg, test_msg)
	} else {
		fmt.Printf("Success!\nSent: %+v\nReceived %+v\n", test_msg, rec_msg)
	}
}
