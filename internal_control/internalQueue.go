package internalcontrol

import (
	"fmt"

	"../hardware/driver-go/elevio"
)

const numFloors int = 5 // 1-indeksert

const numButtons int = 3

var internalQueue [numButtons][numFloors]int

type Button struct {
	Dir   elevio.MotorDirection
	Floor int
}

func printQueue() {
	for button := 0; button < numButtons; button++ {
		for floor := 0; floor < numFloors; floor++ {
			fmt.Println(internalQueue[button][floor])
		}
	}
}

func initQueue() {
	for button := 0; button < numButtons; button++ {
		for floor := 0; floor < numFloors; floor++ {
			internalQueue[button][floor] = 0
		}
	}
	fmt.Println("Initialised internal Queue")

}

func AddOrder(order elevio.ButtonEvent) {
	println("Adding order")
	internalQueue[order.Button][order.Floor] = 1
}

func DeleteOrder(currentFloor int) {
	for i := 0; i < 3; i++ {
		internalQueue[i][currentFloor] = 0
	}
}

func ordersInFloor(currentFloor int) bool {
	for i := 0; i < 3; i++ {
		if internalQueue[i][currentFloor] == 1 {
			return true
		}
	}
	return false
}

func ordersAbove(currentFloor int) bool {
	for i := currentFloor + 1; i < numFloors-1; i++ {
		if internalQueue[0][i] == 1 || internalQueue[1][i] == 1 || internalQueue[2][i] == 1 {
			return true
		}
	}
	return false
}

func ordersBelow(currentFloor int) bool {
	for i := currentFloor - 1; i > -1; i-- {
		if internalQueue[0][i] == 1 || internalQueue[1][i] == 1 || internalQueue[2][i] == 1 {
			return true
		}
	}
	return false
}
