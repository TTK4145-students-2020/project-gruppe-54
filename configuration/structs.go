package structs

import (
	"../hardware/driver-go/elevio"
)

type Channels struct {
	DelegateOrder     chan elevio.ButtonEvent
	OrderCompleted    chan bool
	TakingOrder       chan bool
	TakeExternalOrder chan elevio.ButtonEvent
}
