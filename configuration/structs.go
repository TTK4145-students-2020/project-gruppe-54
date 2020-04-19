package structs

import (
	"../hardware/driver-go/elevio"
)

type Channels struct {
	DelegateOrder      chan elevio.ButtonEvent
	OrderCompleted     chan elevio.ButtonEvent
	CompletedOrder     chan elevio.ButtonEvent
	TakeExternalOrder  chan elevio.ButtonEvent
	MetaData           <-chan MetaData
	UpdateOrderTensor  chan []Node
	CurrentOrderTensor chan []Node
}

type MetaData struct {
	NumNodes  int
	NumFloors int
	Id        int
	ElevPort  string
}

type FloorOrders struct {
	Inside,
	OutsideUp,
	OutsideDown bool
}

type Node struct {
	Floor []FloorOrders
}

type UpdateOrderTensorChan chan []Node
type CurrentOrderTensorChan chan []Node
