package internalcontrol

import (
	"../hardware/driver-go/elevio"
)

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
