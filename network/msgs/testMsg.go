package msgs

type TestMsg struct {
	A  int
	Id int
}

func (msg TestMsg) port() string {
	return TEST_MSG_PORT
}

func (msg *TestMsg) setId(Id int) {
	msg.Id = Id
}

func (msg *TestMsg) GetId() int {
	return msg.Id
}

func (msg *TestMsg) Send() {
	sendTest(*msg)
}

func (msg *TestMsg) Listen() error {
	m, err := listenTest(*msg)
	if err != nil {
		return err
	}
	*msg = m
	return nil
}
