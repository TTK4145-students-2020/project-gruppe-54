package network

import (
	"fmt"
	"testing"
	"time"

	"./msgs"
	"github.com/google/go-cmp/cmp"
)

func TestSendAndListen(t *testing.T) {
	InitNetwork()
	a := 100
	sender := msgs.TestMsg{A: a}
	receiver := msgs.TestMsg{}
	received := make(chan bool)
	go func() {
		err := receiver.Listen()
		if err != nil {
			t.Error("Listen failed with", err)
		}
		received <- true
	}()
	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Millisecond)
		sender.Send()
	}
	<-received
	if !cmp.Equal(sender, receiver) {
		t.Errorf("rec_msg not equal to test_msg\nrec_msg: %+v\ntest_msg: %+v\n", sender, receiver)
	} else {
		fmt.Printf("Success!\nSent: %+v\nReceived %+v\n", sender, receiver)
	}
}
