package internal_control

import (
	"time"

	"../hardware/driver-go/elevio"
	"../supervisor/watchdog"
)

var state int

var Floor int

const (
	IDLE      = 0
	DRIVE     = 1
	DOOR_OPEN = 2
)

func FsmInit() {
	state = IDLE
	Floor = 0
}

func FsmUpdateFloor(newFloor int) {
	Floor = newFloor
	println("Floor: %d%v", Floor)
}

func FSM() {
	//println("Floor: ", Floor)
	switch state {
	case IDLE:
		if ordersAbove(Floor) {
			println("order above, current Floor: ", Floor)
			elevio.SetMotorDirection(elevio.MD_Up)
			state = DRIVE
		}
		if ordersBelow(Floor) {
			println("order below, current Floor: ", Floor)
			elevio.SetMotorDirection(elevio.MD_Down)
			state = DRIVE
		}
	case DRIVE:

		if ordersInFloor(Floor) {
			elevio.SetMotorDirection(elevio.MD_Stop)
			state = DOOR_OPEN
		}
	case DOOR_OPEN:

		elevio.SetDoorOpenLamp(true)
		wd := watchdog.NewWatchdog(time.Second * 3)
		if <-wd.getKickChannel() {
			elevio.SetDoorOpenLamp(false)
			wd.Stop()
		}
		orderDone <- true
		if ordersAbove(Floor) {
			elevio.SetMotorDirection(elevio.MD_Up)
			direction = elevio.MD_Up
			state = DRIVE
		} else if ordersBelow(Floor) {
			elevio.SetMotorDirection(elevio.MD_Down)
			direction = elevio.MD_Down
			state = DRIVE
		} else {
			state = IDLE
		}
	}
}
