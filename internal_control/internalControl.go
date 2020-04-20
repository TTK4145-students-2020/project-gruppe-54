package internalcontrol

import (
	c "../configuration"
	"../hardware/driver-go/elevio"
)

// InternalControl .. Responsible for internal control of a single elevator
func InternalControl(ch c.Channels) {
	var numFloors int = 4
	println("Connecting to server")
	port := (<-ch.MetaData).ElevPort
	elevio.Init("localhost:"+port, numFloors)

	initQueue()
	FsmInit()
	drvButtons := make(chan elevio.ButtonEvent)
	drvFloors := make(chan int)
	drvStop := make(chan bool)
	doorsOpen := make(chan int)

	go elevio.PollButtons(drvButtons)
	go elevio.PollFloorSensor(drvFloors)
	go elevio.PollStopButton(drvStop)
	go FSM(doorsOpen)
	for {
		select {
		case floor := <-drvFloors: //Sensor senses a new floor
			FsmUpdateFloor(floor)
		case drvOrder := <-drvButtons: // a new button is pressed on this elevator
			ch.DelegateOrder <- drvOrder //Delegate this order
		case ExtOrder := <-ch.TakeExternalOrder:
			AddOrder(ExtOrder)
		case floor := <-doorsOpen:
			orderOutsideUpCompleted := elevio.ButtonEvent{
				Floor:  floor,
				Button: elevio.BT_HallUp,
			}
			orderOutsideDownCompleted := elevio.ButtonEvent{
				Floor:  floor,
				Button: elevio.BT_HallDown,
			}
			orderInsideCompleted := elevio.ButtonEvent{
				Floor:  floor,
				Button: elevio.BT_Cab,
			}
			ch.CompletedOrder <- orderOutsideUpCompleted
			ch.CompletedOrder <- orderOutsideDownCompleted
			ch.CompletedOrder <- orderInsideCompleted
		}

	}
}
