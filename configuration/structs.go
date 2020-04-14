package structs

import (
	"../hardware/driver-go/elevio"
)

type Channels struct {
	DelegateOrder     chan elevio.ButtonEvent
	OrderCompleted    chan elevio.ButtonEvent
	TakingOrder       chan elevio.ButtonEvent
	TakeExternalOrder chan elevio.ButtonEvent
	MetaData          <-chan MetaData
}

type MetaData struct {
	NumNodes  int
	NumFloors int
	Id        int
	ElevPort  string
}
