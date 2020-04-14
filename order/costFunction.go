package order

import (
	"math/rand"

	"../hardware/driver-go/elevio"
	"../network/msgs"
)

// func calculateCost(order elevio.ButtonEvent) int {
// 	// add something for order.ButtonType == BT_Cab e.g. extremely high as it should
// 	var cost = order.Floor - elevio.GetFloor() // Må lage getFloor, getDirection
// 	if cost == 0 && getDirection() != MotorDirection.MD_Stop {
// 		cost += 4
// 	}
// 	if cost > 0 && getDirection() == MotorDirection.MD_Down {
// 		cost += 3
// 	}
// 	if cost < 0 && getDirection() == MotorDirection.MD_Up {
// 		cost += 3
// 	}
// 	if cost != 0 && getDirection() == MotorDirection.MD_Stop {
// 		cost++
// 	}
// 	return cost
// }

func calculateCost(order elevio.ButtonEvent) uint { return uint(rand.Intn(10)) }

//getCostFuction

func sendCost(order elevio.ButtonEvent, id int) {
	cost := calculateCost(order)
	costMsg := msgs.CostMsg{Cost: cost}
	costMsg.Send()
}
