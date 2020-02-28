package internal_control

import (
	"../hardware/driver-go/elevio"
)

var state int

var floor int

const (
	IDLE      = 0
	DRIVE     = 1
	DOOR_OPEN = 2
)

func FsmInit() {
	state = IDLE
}

func FsmUpdateFloor(newFloor int) {
	floor = newFloor
	println("floor: %d%v", floor)
}

func FSM() {
	//elevio.SetMotorDirection(elevio.MD_Up)

	switch state {
	case IDLE:
		if ordersAbove(floor) {
			println("order above, current floor: ", floor)
			elevio.SetMotorDirection(elevio.MD_Up)
			//direction = elevio.MD_Up
			state = DRIVE
		}
		if ordersBelow(floor) {
			println("order below, current floor: ", floor)
			elevio.SetMotorDirection(elevio.MD_Down)
			//direction = elevio.MD_Down
			state = DRIVE
		}
	case DRIVE:

		if ordersInFloor(floor) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			//direction = elevio.MD_Stop
			state = DOOR_OPEN
		}
	case DOOR_OPEN:

		/*elevio.SetDoorOpenLamp(true)
		wd := watchdog.NewWatchdog(time.Second * 3)
		if <-wd.getKickChannel() {
			elevio.SetDoorOpenLamp(false)
			wd.Stop()
		}
		orderDone <- true
		if ordersAbove(floor) {
			elevio.SetMotorDirection(elevio.MD_Up)
			direction = elevio.MD_Up
			state = DRIVE
		} else if ordersBelow(floor) {
			elevio.SetMotorDirection(elevio.MD_Down)
			direction = elevio.MD_Down
			state = DRIVE
		} else {
			state = IDLE
		}*/
	}
}
