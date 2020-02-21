package driver


import "./elevio"
import "internalQueue"
import "fmt"

func main(){

    numFloors := 4

    elevio.Init("localhost:15657", numFloors)
    
    var direction elevio.MotorDirection = elevio.MD_Up
    elevio.SetMotorDirection(direction)
    
    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_stop    := make(chan bool)    
    
    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollStopButton(drv_stop)
	
	var state int
	const (
		IDLE      = 0
		DRIVE    = 1
		DOOR_OPEN = 2
	)

	doorTimedOut := time.NewTimer(3 * time.Second)
	engineErrorTimer := time.NewTimer(3 * time.Second)
	doorTimedOut.Stop()
	engineErrorTimer.Stop()
	orderCleared := false
	ch.Elevator <- elevator

    
    for {
		switch State{
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
		}
        
    }    
}
