package msgs

import (
	"errors"
)

type ackMsg struct {
	Msg interface{}
	Id  int
}

func (msg ackMsg) port() string {
	return ACK_MSG_PORT
}

func (msg ackMsg) Send() {
	send(&msg)
}

func (msg *ackMsg) setId(Id int) {
	msg.Id = Id
}

func (msg *ackMsg) GetId() int {
	return msg.Id
}

func (msg *ackMsg) Listen() error {
	m, err := listen(msg)
	if err != nil {
		return err
	}
	if m, ok := m.(*ackMsg); ok {
		*msg = *m
	} else {
		return errors.New("failed casting to msg type after listen")
	}
	return nil
}

func (msg *ackMsg) Ack() {}
