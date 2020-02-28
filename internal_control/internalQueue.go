package internal_control

import "../hardware/driver-go/elevio"


var numFloors int = 4	
var internalQueue [3][numFloors]int

type Button struct {
	Dir   elevio.MotorDirection
	Floor int
}


func ChooseDirection() Dir{
	switch elevio.MotorDirection {
	case elevio.MD_Stop:
		if ordersAbove(elevator) {
			return elevio.MD_Up
		} else if ordersBelow(elevator) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}
	case elevio.MD_Up:
		if ordersAbove(elevator) {
			return elevio.MD_Up
		} else if ordersBelow(elevator) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Down:
		if ordersBelow(elevator) {
			return elevio.MD_Down
		} else if ordersAbove(elevator) {
			return elevio.MD_Up
		} else {
			return elevio.MD_Stop
		}
	}
	return elevio.MD_Stop
}

func AddOrder(Floor int, typeOfOrder int) {
	internalQueue[typeOfOrder][Floor] = 1
}

func DeleteOrder(currentFloor int) {
	for i := 0; i < 3; i++ {
		internalQueue[i][currentFloor] = 0
	}
}

func ordersInFloor(currentFloor int ) int{
	for i := 0; i < 3; i++ {
		if internalQueue[i][currentFloor] == 1 {
			return 1
		}
	}
	return 0
}


func ordersAbove(currentFloor int) int {
	for i := currentFloor + 1; i < numFloors-1; i++ {
		if internalQueue[0][i] == 1 || internalQueue[1][i] == 1 || internalQueue[2][i] == 1 {
			return 1
		}
	}
	return 0
}

func ordersBelow(currentFloor int) int {
	for i := currentFloor - 1; i > -1; i-- {
		if internalQueue[0][i] == 1 || internalQueue[1][i] == 1 || internalQueue[2][i] == 1 {
			return 1
		}
	}
	return 0
}

