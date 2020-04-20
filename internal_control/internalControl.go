package internalcontrol

import (
	"fmt"

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
	i := 0
	for {
		i++
		if i%100000 == 0 {
			fmt.Println("InternalControl")
		}
		select {
		case floor := <-drvFloors: //Sensor senses a new floor
			fmt.Println("floor")
			//println("updating floor:", floor)
			FsmUpdateFloor(floor)
		case drvOrder := <-drvButtons: // a new button is pressed on this elevator
			fmt.Println("drvButtons")
			ch.DelegateOrder <- drvOrder //Delegate this order
		case ExtOrder := <-ch.TakeExternalOrder:
			fmt.Println("TakeExternalOrder")
			AddOrder(ExtOrder)
		case floor := <-doorsOpen:
			fmt.Println("doorsOpen")
			order_OutsideUp_Completed := elevio.ButtonEvent{
				Floor:  floor,
				Button: elevio.BT_HallUp,
			}
			order_OutsideDown_Completed := elevio.ButtonEvent{
				Floor:  floor,
				Button: elevio.BT_HallDown,
			}
			order_Inside_Completed := elevio.ButtonEvent{
				Floor:  floor,
				Button: elevio.BT_Cab,
			}
			ch.CompletedOrder <- order_OutsideUp_Completed
			ch.CompletedOrder <- order_OutsideDown_Completed
			ch.CompletedOrder <- order_Inside_Completed
		}

	}
}
