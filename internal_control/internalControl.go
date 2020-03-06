package internalcontrol

import (
	c "../configuration"
	"../hardware/driver-go/elevio"
)

// InternalControl .. Responsible for internal control of a single elevator
func InternalControl(ch c.Channels) {
	var numFloors int = 4
	println("Connecting to server")
	elevio.Init("localhost:15657", numFloors)

	initQueue()
	FsmInit()
	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvStop := make(chan bool)

	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollStopButton(drvStop)
	go FSM()
	for {
		select {
		case floor := <-drvFloors: //Sensor senses a new floor
			//println("updating floor:", floor)
			FsmUpdateFloor(floor)
		case drvOrder := <-drvButtons: // a new button is pressed on this elevator
			ch.DelegateOrder <- drvOrder //Delegate this order
		case ExtOrder := <-ch.TakeExternalOrder:
			AddOrder(ExtOrder)
			ch.TakingOrder <- true //if nothing fails tell order it can update matrix..
		}

	}
}
