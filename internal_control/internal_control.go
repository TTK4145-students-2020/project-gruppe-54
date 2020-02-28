package internal_control

import "../hardware/driver-go/elevio"
import "fmt"

func InternalControl(){

    var numFloors int = 4

    elevio.Init("localhost:15657", numFloors)
    
    var direction elevio.MotorDirection = elevio.MD_Up
    elevio.SetMotorDirection(direction)
    
    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
	drv_stop    := make(chan bool)    
	orderDone   := make(chan bool)
    
    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollStopButton(drv_stop)
	

	a:= <- drv_floors
	fmt.Println(a)
    
    for {
		select{
		case order := <- drv_buttons:
			AddOrder(order.Floor, order.Button)
			/*switch state{
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
					}*/
		case a := <- orderDone:
			DeleteOrder(drv_floors)

		default:
			fmt.Println("TAKABOOM")
			//main loop for internalqueue
		}
	}
}
		

