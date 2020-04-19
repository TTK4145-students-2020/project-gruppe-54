package internalcontrol

import (
	"fmt"
	"time"

	"../hardware/driver-go/elevio"
)

var state int
var dir elevio.MotorDirection
var Floor int

const (
	IDLE      = 0
	DRIVE     = 1
	DOOR_OPEN = 2
)

func FsmInit() {

	state = IDLE
	// Needs to start in a well-defined state
	for Floor = elevio.GetFloor(); Floor < 0; Floor = elevio.GetFloor() {
		elevio.SetMotorDirection(elevio.MD_Up)
		time.Sleep(1 * time.Second)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	fmt.Println("FSM initialized!")
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
				dir = elevio.MD_Up
				elevio.SetMotorDirection(dir)
				state = DRIVE
			}
			if ordersBelow(Floor) {
				//println("order below, going down, current Floor: ", Floor)
				dir = elevio.MD_Down
				elevio.SetMotorDirection(dir)
				state = DRIVE
			}
			if ordersInFloor(Floor) {
				//println("order below, going down, current Floor: ", Floor)
				state = DOOR_OPEN
			}
		case DRIVE:
			if ordersInFloor(Floor) { // this is the problem : the floor is being kept constant at e.g. 2 while its moving
				dir = elevio.MD_Stop
				elevio.SetMotorDirection(dir)
				state = DOOR_OPEN
			}
		case DOOR_OPEN:
			elevio.SetDoorOpenLamp(true)
			dir = elevio.MD_Stop
			elevio.SetMotorDirection(dir)
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

func GetDirection() elevio.MotorDirection {
	return dir
}
