package internalcontrol

import (
	"../hardware/driver-go/elevio"
)

//OrderStruct ... orderstruct
type OrderStruct struct {
	Floor  int
	button OrderType
}

// OrderType ... type of order
type OrderType int

const (
	HallUp OrderType = 0
	HallDn           = 1
	Cab              = 2
)

// ButtonTypeMap ... maps buttontype to ordertype
var ButtonTypeMap = map[elevio.ButtonType]OrderType{
	elevio.BT_HallUp:   HallUp,
	elevio.BT_HallDown: HallDn,
	elevio.BT_Cab:      Cab,
}

type ElevatorDirection int

const (
	Stop ElevatorDirection = 0
	Down                   = -1
	Up                     = 1
)
// ElevatorDirectionMap ... maps elevatorDirection to MotorDirection
var ElevatorDirectionMap = map[elevio.MotorDirection]ElevatorDirection{
	elevio.MD_Stop: Stop,
	elevio.MD_Down: Down,
	elevio.MD_Up:   Up,
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
			//AddOrder(drvOrder.Floor, drvOrder.Button)

		}
	}
}
