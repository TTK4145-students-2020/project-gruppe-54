package order

import (
	"../hardware/driver-go/elevio"
	ic "../internal_control"
	"../network/msgs"
)

func calculateCost(order elevio.ButtonEvent) uint {
	var cost = abs(order.Floor - elevio.GetFloor())
	if cost == 0 && ic.GetDirection() != elevio.MD_Stop {
		cost += 4
	}
	if cost > 0 && (ic.GetDirection() == elevio.MD_Down || ic.GetDirection() == elevio.MD_Up) {
		cost += 3
	}
	if cost != 0 && ic.GetDirection() == elevio.MD_Stop {
		cost++
	}
	return cost
}

func abs(x int) uint {
	if x < 0 {
		return uint(-x)
	} else {
		return uint(x)
	}
}

func sendCost(order elevio.ButtonEvent, id int) {
	cost := calculateCost(order)
	costMsg := msgs.CostMsg{Cost: cost}
	costMsg.Send()
}
