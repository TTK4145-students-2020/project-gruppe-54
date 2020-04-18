package order

import (
	"fmt"
	"math"
	"time"

	c "../configuration"
	"../hardware/driver-go/elevio"
	"../network/msgs"
	sv "../supervisor"
)

var numNodes int = 3
var numFloors int = 4

func listenForOrderCompleted(order elevio.ButtonEvent, isDone chan bool, giveUp chan bool) {
	orderComplete := false //Odin sett in noe her om at den ser på matrisen og hvis den får inn at det er utført blir det true

	// Listen for diff in order tensor
	orderTensorDiffMsg := msgs.OrderTensorDiffMsg{}
	receivedDiff := make(chan msgs.OrderTensorDiffMsg)
	killPolling := make(chan bool)

	defer func() {
		killPolling <- true
		if orderComplete {
			isDone <- true
		}
	}()

	// Start polling for diff to tensor
	go func() {
		for {
			err := orderTensorDiffMsg.Listen()
			if err != nil {
				continue
			}
			receivedDiff <- orderTensorDiffMsg
			if <-killPolling {
				return
			}
		}
	}()
	for {
		select {
		case receivedDiffMsg := <-receivedDiff:
			if receivedDiffMsg.Order.Floor == order.Floor && receivedDiffMsg.Diff == msgs.DIFF_REMOVE {
				fmt.Printf("Order %+v completed!\n", receivedDiffMsg)
				orderComplete = true
				return
			}
		case <-giveUp:
			fmt.Println("Order not completed :(")
			orderComplete = false
			return
		}
	}
}

func delegateOrder(order elevio.ButtonEvent, ch c.Channels) {
	orderMsg := msgs.OrderMsg{Order: order}
	for i := 0; i < 5; i++ {
		orderMsg.Send()
		time.Sleep(1 * time.Millisecond)
	}

	isNotDone := make(chan bool, 1)
	isDone := make(chan bool, 1)

	go sv.WatchOrder(isNotDone)
	go listenForOrderCompleted(order, isDone, isNotDone)
	select {
	case <-isNotDone:
		fmt.Printf("Order %+v is not done\n", order)
		// FIXME: Dirty fix
		close(isDone)
		delegateOrder(order, ch) //redelegate
	case <-isDone:
		close(isNotDone)
		// Nothing to do
	}
}

func collectCosts(numNodes int) []uint {
	costs := make([]uint, numNodes)
	for i := 0; i < numNodes; i++ {
		costs[i] = uint(math.Inf(0)) // Initialize all costs to infinite
	}
	stopPolling := make(chan bool)
	complete := make(chan bool)
	// Stop polling after 50 ms
	timer := time.NewTimer(50 * time.Millisecond)
	newCost := make(chan msgs.CostMsg)
	go func() {
		for {
			select {
			case newCostMsg := <-newCost:
				if id := newCostMsg.GetId(); id >= 0 && id < numNodes {
					costs[id] = newCostMsg.Cost
				}
			case <-timer.C:
				stopPolling <- true
				return
			}
		}
	}()
	go func() {
		costMsg := msgs.CostMsg{}
		for {
			select {
			case <-stopPolling:
				// fmt.Println("Stopped polling")
				complete <- true
				return
			default:
				// fmt.Println("Listening for costs")
				err := costMsg.Listen()
				if err != nil {
					// fmt.Printf("Error: %s\n", err)
					continue
				}
				// fmt.Printf("Received from %d\n", costMsg.GetId())
				newCost <- costMsg
			}
		}
	}()
	<-complete
	return costs
}

//ControlOrders ... Delegates new orders it receives on channel newOrders
func ControlOrders(ch c.Channels) {
	go checkForExternalOrders(ch) //Continously check for new orders given to this elev
	for {
		select {
		case newOrder := <-ch.DelegateOrder:
			// case <-ch.DelegateOrder:
			go delegateOrder(newOrder, ch)
		case orderTaken := <-ch.TakingOrder: // the external order has been taken
			orderTensorDiffMsg := msgs.OrderTensorDiffMsg{
				Order: orderTaken,
				Diff:  msgs.DIFF_REMOVE,
				Id:    (<-ch.MetaData).Id}
			UpdateOrderTensor(ch.UpdateOrderTensor, ch.CurrentOrderTensor, orderTensorDiffMsg)
			// fmt.Printf("PC %d completed order %+v\n", (<-ch.MetaData).Id, orderTaken)
			for i := 0; i < 5; i++ {
				orderTensorDiffMsg.Send()
				time.Sleep(1 * time.Millisecond)
			}
		}
	}
}

func checkForExternalOrders(ch c.Channels) {
	newOrder := msgs.OrderMsg{Order: elevio.ButtonEvent{Button: elevio.BT_HallUp, Floor: 0}, Id: 0}
	metaData := <-ch.MetaData
	for {
		err := newOrder.Listen()
		if err != nil {
			continue
		}
		// fmt.Printf("Received order: %+v\n", newOrder)
		// go func() {
		// 	cost := calculateCost(newOrder.Order)
		// 	costMsg := msgs.CostMsg{Cost: cost}
		// 	for i := 0; i < 5; i++ {
		// 		costMsg.Send()
		// 		time.Sleep(1 * time.Millisecond)
		// 	}
		// 	fmt.Printf("Sending cost: %+v\n", costMsg)
		// }()
		// costs := collectCosts(metaData.NumNodes)
		// Find minimum
		var costs [10]uint
		min := uint(math.Inf(0))
		minId := metaData.Id
		// fmt.Println("Printing costs...")
		// for id, cost := range costs {
		// 	fmt.Printf("Id: %d, cost: %d\n", id, cost)
		// 	if cost < min {
		// 		min = cost
		// 		minId = id
		// 	}
		// }
		// fmt.Printf("MinId: %d, MinCost: %d\n", minId, min)
		// Needs to ensure that order is taken, if min is itself
		orderDiff := msgs.OrderTensorDiffMsg{
			Order: newOrder.Order,
			Diff:  msgs.DIFF_ADD,
			Id:    minId,
		}
		UpdateOrderTensor(ch.UpdateOrderTensor, ch.CurrentOrderTensor, orderDiff)
		if minId == metaData.Id || min == costs[metaData.Id] {
			ch.TakeExternalOrder <- newOrder.Order
			orderTensorDiffMsg := msgs.OrderTensorDiffMsg{
				Order: newOrder.Order,
				Diff:  msgs.DIFF_ADD,
				Id:    metaData.Id,
			}
			UpdateOrderTensor(ch.UpdateOrderTensor, ch.CurrentOrderTensor, orderTensorDiffMsg)
			// fmt.Printf("PC %d takes order %+v\n", (<-ch.MetaData).Id, newOrder.Order)
			for i := 0; i < 5; i++ {
				orderTensorDiffMsg.Send()
				time.Sleep(1 * time.Millisecond)
			}
		}
	}
}

func equalOrders(order1, order2 elevio.ButtonEvent) bool {
	return order1.Floor == order2.Floor && order1.Button == order2.Button
}
