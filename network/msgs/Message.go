package msgs

const (
	BROADCAST_ADDR = "255.255.255.255"
	UDP_TIMEOUT    = 50 // ms
)

type messager interface {
	Send()
	Listen()
	port() string
}
