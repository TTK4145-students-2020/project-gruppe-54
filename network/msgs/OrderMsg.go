package msgs

import (
	"errors"

	"../../hardware/driver-go/elevio"
)

type OrderMsg struct {
	Order elevio.ButtonEvent
	Id    int
}

func (msg OrderMsg) port() string {
	return ORDER_MSG_PORT
}

func (msg OrderMsg) Send() {
	send(&msg)
}

func (msg *OrderMsg) Listen() error {
	m, err := listen(msg)
	if err != nil {
		return err
	}
	if m, ok := m.(*OrderMsg); ok {
		*msg = *m
	} else {
		return errors.New("failed casting to msg type after listen")
	}
	return nil
}
