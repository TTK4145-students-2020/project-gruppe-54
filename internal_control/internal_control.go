package internalcontrol

import (
	"../hardware/driver-go/elevio"
)
type Order struct {
	button ButtonType
	floor  int
}

const (
	HallUp ButtonType = 0
	HallDn           = 1
	Cab              = 2
)

// ButtonTypeToOrderTypeMap solves the trouble of having two enums
var ButtonTypeToOrderTypeMap = map[elevio.ButtonType]OrderType{
	elevio.BT_HallUp:   HallUp,
	elevio.BT_HallDown: HallDn,
	elevio.BT_Cab:      Cab,
}


type MotorDirection int

const (
	Stop MotorDirection = 0
	Down                = -1
	Up 					= 1
)

var MotorDirectionMap = map[elevio.ButtonType]OrderType{
	elevio.MotorDirection.MD_stop:   HallUp,
	elevio.BT_HallDown: HallDn,
	elevio.BT_Cab:      Cab,
}

// InternalControl .. Responsible for internal control of a single elevator
func InternalControl() {

	var numFloors int = 4
	println("Connecting to server")
	elevio.Init("localhost:15657", numFloors)

	initQueue()
	FsmInit()

	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvStop := make(chan bool)

	newOrders := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollStopButton(drvStop)

	go FSM()
	for {
		select {
		case floor := <-drvFloors:
			//println("updating floor:", floor)
			FsmUpdateFloor(floor)
		case drvOrder := <-drvButtons:
			//println("new order")
			newOrders <- drvOrder
			AddOrder(drvOrder.Floor, drvOrder.Button)

		}
	}
}
