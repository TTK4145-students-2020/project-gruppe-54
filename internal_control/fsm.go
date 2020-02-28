package internal_control

import "../hardware/driver-go/elevio"

var state int
const (
	IDLE      = 0
	DRIVE     = 1
	DOOR_OPEN = 2
)

func fsmInit() {
	state = IDLE
}


/*
func FsmButtonHandler(order elevio.ButtonEvent){
	switch state{
		case IDLE:
			if ordersAbove(drv_floors){
				SetMotorDirection(elevio.MD_Up)
				direction = elevio.MD_Up
				state = DRIVE
			}
			if ordersBelow(drv_floors){
				SetMotorDirection(elevio.MD_Down)
				direction = elevio.MD_Down
				state = DRIVE
			}
		case DRIVE:
			if ordersInFloor(drv_floors){
				SetMotorDirection(elevio.MD_Stop)
				direction = elevio.MD_Stop
				state = DOOR_OPEN
			}
		case DOOR_OPEN:
			SetDoorOpenLamp(true)
			wd := watchdog.NewWatchdog(time.Second*3)
			if <-wd.getKickChannel(){
				SetDoorOpenLamp(false)
				wd.Stop()
			}
			orderDone <- true
			if ordersAbove(drv_floors){
				SetMotorDirection(elevio.MD_Up)
				direction = elevio.MD_Up
				state = DRIVE
			}else if ordersBelow(drv_floors){
				SetMotorDirection(elevio.MD_Down)
				direction = elevio.MD_Down
				state = DRIVE
			}else{
				state = IDLE
			}
		}
	}

*/