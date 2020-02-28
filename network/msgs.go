package network

type MsgCh interface {
	Ch() chan interface{}
}

type CostMsg struct {
	Cost float64
	ID   int
}

type CostMsgCh struct {
	channel chan interface{}
}

func (self CostMsgCh) Ch() chan interface{} {
	return (self).channel
}

type TestMsg struct {
	A int
}

type networkMsg interface {
	a()
}

func (msg TestMsg) a() {}

type TestMsgCh struct {
	channel chan interface{}
}

func (self TestMsgCh) Ch() chan interface{} {
	return (self).channel
}

const (
	BROADCAST_ADDR = "255.255.255.255"
	COST_MSG_PORT  = "3000"
	TEST_MSG_PORT  = "15000"
)
