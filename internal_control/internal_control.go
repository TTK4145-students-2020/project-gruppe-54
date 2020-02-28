package internal_control

import (
	"../hardware/driver-go/elevio"
)

// InternalControl .. Responsible for internal control of a single elevator
func InternalControl() {

	var numFloors int = 4
	println("Connecting to server")
	elevio.Init("localhost:15657", numFloors)

	initQueue()
	printQueue()
	FsmInit()
	//var direction elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(direction)
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_stop := make(chan bool)
	//orderDone := make(chan bool)
	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollStopButton(drv_stop)
	println("HERE")

	for {
		go FSM()
		select {
		case floor := <-drv_floors:
			println("updating floor")
			FsmUpdateFloor(floor)

		case order := <-drv_buttons:
			println("new order")
			AddOrder(order.Floor, order.Button)
			printQueue()
		/*case a := <-orderDone:
		if a {
			DeleteOrder(floor)
		}*/
		default:
			//fmt.Println("TAKABOOM")
		}
	}
}
