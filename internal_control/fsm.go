package internalcontrol

import (
	"time"

	"../hardware/driver-go/elevio"
)

var state int

var Floor int

const (
	IDLE      = 0
	DRIVE     = 1
	DOOR_OPEN = 2
)

func FsmInit() {

	state = 0
	Floor = elevio.GetFloor()
	if Floor == -1 {
		Floor = 0
	}
}

func FsmUpdateFloor(newFloor int) {
	Floor = newFloor
}

func FSM(doorsOpen chan<- int) {
	for {
		switch state {
		case IDLE:
			if ordersAbove(Floor) {
				//println("order above,going up, current Floor: ", Floor)
				elevio.SetMotorDirection(elevio.MD_Up)
				state = DRIVE
			}
			if ordersBelow(Floor) {
				//println("order below, going down, current Floor: ", Floor)
				elevio.SetMotorDirection(elevio.MD_Down)
				state = DRIVE
			}
			if ordersInFloor(Floor) {
				//println("order below, going down, current Floor: ", Floor)
				state = DOOR_OPEN
			}
		case DRIVE:
			if ordersInFloor(Floor) { // this is the problem : the floor is being kept constant at e.g. 2 while its moving
				elevio.SetMotorDirection(elevio.MD_Stop)
				state = DOOR_OPEN
			}
		case DOOR_OPEN:
			elevio.SetDoorOpenLamp(true)
			elevio.SetMotorDirection(elevio.MD_Stop)
			println("DOOR OPEN")
			DeleteOrder(Floor)
			doorsOpen <- Floor
			timer1 := time.NewTimer(2 * time.Second)
			<-timer1.C
			elevio.SetDoorOpenLamp(false)
			println("DOOR CLOSE")
			state = IDLE
		}
	}
}
