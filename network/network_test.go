package network

import (
	"fmt"
	"testing"
	"time"

	c "../configuration"
	"./msgs"
)

func TestSendAndListen(t *testing.T) {

	testMetaDataSender := c.MetaData{NumNodes: 2, NumFloors: 0, Id: 12}
	testMetaDataReceiver := c.MetaData{NumNodes: 2, NumFloors: 0, Id: 1}

	testMetaDataSenderChan := make(chan c.MetaData, 1)
	testMetaDataReceiverChan := make(chan c.MetaData, 1)
	go func() {
		for {
			testMetaDataSenderChan <- testMetaDataSender
			testMetaDataReceiverChan <- testMetaDataReceiver
		}
	}()
	msgs.InitTestMessage(testMetaDataSenderChan, testMetaDataReceiverChan)

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
	if sender.A != receiver.A {
		t.Errorf("rec_msg not equal to test_msg\nrec_msg: A: %d, id: %d\ntest_msg: %+v\n", sender.A, sender.GetId(), receiver)
	} else {
		fmt.Printf("Success!\nSent: A: %d, id: %d\nReceived: A: %d, id: %d\n", sender.A, sender.GetId(), receiver.A, receiver.GetId())
	}
}
