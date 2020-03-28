package msgs

import (
	"errors"
)

type CostMsg struct {
	Cost uint
	id   int
}

func (msg CostMsg) port() string {
	return COST_MSG_PORT
}

func (msg CostMsg) Send() {
	send(&msg)
}

func (msg *CostMsg) setId(Id int) {
	msg.id = Id
}

func (msg CostMsg) GetId() int {
	return msg.id
}

func (msg *CostMsg) Listen() error {
	m, err := listen(msg)
	if err != nil {
		return err
	}
	if m, ok := m.(*CostMsg); ok {
		*msg = *m
	} else {
		return errors.New("failed casting to msg type after listen")
	}
	return nil
}
