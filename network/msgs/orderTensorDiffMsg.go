package msgs

import (
	"errors"

	"../../hardware/driver-go/elevio"
)

type OrderTensorDiff int

const (
	DIFF_ADD    OrderTensorDiff = 0
	DIFF_REMOVE                 = 1
)

type OrderTensorDiffMsg struct {
	Order elevio.ButtonEvent
	Diff  OrderTensorDiff
	Id    int
}

func (msg OrderTensorDiffMsg) port() string {
	return ORDER_TENSOR_DIFF_MSG_PORT
}

func (msg *OrderTensorDiffMsg) setId(Id int) {
	msg.Id = Id
}

func (msg *OrderTensorDiffMsg) GetId() int {
	return msg.Id
}

func (msg OrderTensorDiffMsg) Send() {
	send(&msg)
}

func (msg *OrderTensorDiffMsg) Listen() error {
	m, err := listen(msg)
	if err != nil {
		return err
	}
	if m, ok := m.(*OrderTensorDiffMsg); ok {
		*msg = *m
	} else {
		return errors.New("failed casting to msg type after listen")
	}
	return nil
}
