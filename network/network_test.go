package network

import (
	"bytes"
	"encoding/gob"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestMsg struct {
	A int
}

func TestInitSender(t *testing.T) {
	addr := "localhost"
	port := "3000"
	msg_ch := InitSender(addr)
	test_msg := TestMsg{1}
	ln, err := net.ListenPacket("udp", ":"+port)
	if err != nil {
		t.Errorf("Could not initialise server at %s\n", addr)
	}
	defer ln.Close()
	rec_ch := make(chan TestMsg)
	go func() {
		rec_msg := TestMsg{}
		inBytes := make([]byte, 4096)
		n, _, err := ln.ReadFrom(inBytes)
		if err != nil {
			t.Error("Could not read from UDP")
		}
		buf := bytes.NewBuffer(inBytes[:n])
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&rec_msg)
		if err != nil {
			t.Error("Could not decode received msg")
		}
		rec_ch <- rec_msg
	}()
	msg_ch <- test_msg
	rec_msg := <-rec_ch
	if !cmp.Equal(rec_msg, test_msg) {
		t.Errorf("rec_msg not equal to test_msg\nrec_msg: %+v\ntest_msg: %+v\n", rec_msg, test_msg)
	}

}
