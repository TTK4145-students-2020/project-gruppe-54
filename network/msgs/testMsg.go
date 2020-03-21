package msgs

import "errors"

type TestMsg struct {
	A int
}

func (msg TestMsg) port() string {
	return TEST_MSG_PORT
}

func (msg TestMsg) Send() {
	send(&msg)
}

func (msg *TestMsg) Listen() error {
	m, err := listen(msg)
	if err != nil {
		return err
	}
	if m, ok := m.(*TestMsg); ok {
		*msg = *m
	} else {
		return errors.New("failed casting to msg type after listen")
	}
	return nil
}
