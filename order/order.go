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

func listenForOrderCompleted(order elevio.ButtonEvent, isDone chan bool, isNotDone chan bool) {
	orderComplete := false //Odin sett in noe her om at den ser på matrisen og hvis den får inn at det er utført blir det true

	// Listen for diff in order tensor
	orderTensorDiffMsg := msgs.OrderTensorDiffMsg{}
	receivedDiff := make(chan msgs.OrderTensorDiffMsg)
	killPolling := make(chan bool)

	defer func() {
		if orderComplete {
			// Propagate the signal
			isDone <- true
		} else {
			// Propagate the signal
			isNotDone <- true
		}
		killPolling <- true
	}()

	// Start polling for diff to tensor
	go func() {
		for {
			select {
			case <-killPolling:
				return
			default:
				err := orderTensorDiffMsg.Listen()
				if err != nil {
					continue
				}
				receivedDiff <- orderTensorDiffMsg
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
		case <-isNotDone:
			fmt.Println("Order not completed :(")
			orderComplete = false
			return
		}
	}
}

func delegateOrder(orderMsg msgs.OrderMsg, ch c.Channels) {
	for i := 0; i < 5; i++ {
		orderMsg.Send()
		time.Sleep(1 * time.Millisecond)
	}
}

func collectCosts(numNodes int) []uint {
	costs := make([]uint, numNodes)
	for i := 0; i < numNodes; i++ {
		costs[i] = uint(math.Inf(0)) // Initialize all costs to infinite
	}
	stopPolling := make(chan bool)
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
	costMsg := msgs.CostMsg{}
L:
	for {
		select {
		case <-stopPolling:
			break L
		default:
			err := costMsg.Listen()
			if err != nil {
				continue
			}
			newCost <- costMsg
		}
	}
	return costs
}

//ControlOrders ... Delegates new orders it receives on channel newOrders
func ControlOrders(ch c.Channels) {
	go checkForNewOrders(ch) //Continously check for new orders given to this elev
	go checkForAcceptedOrders(ch)
	for {
		select {
		case newOrder := <-ch.DelegateOrder:
			orderMsg := msgs.OrderMsg{Order: newOrder}
			orderMsg.Id = (<-ch.MetaData).Id
			delegateOrder(orderMsg, ch)
		case orderCompleted := <-ch.CompletedOrder: // the external order has been taken
			orderTensorDiffMsg := msgs.OrderTensorDiffMsg{
				Order: orderCompleted,
				Diff:  msgs.DIFF_REMOVE,
				Id:    (<-ch.MetaData).Id}
			UpdateOrderTensor(ch.UpdateOrderTensor, ch.CurrentOrderTensor, orderTensorDiffMsg)
			for i := 0; i < 5; i++ {
				orderTensorDiffMsg.Send()
				time.Sleep(1 * time.Millisecond)
			}
		}
	}
}

func checkForNewOrders(ch c.Channels) {
	newOrder := msgs.OrderMsg{}
	metaData := <-ch.MetaData
	for {
		err := newOrder.Listen()
		if err != nil {
			continue
		}
		// If the order comes from inside, only the current elevator can complete it
		if newOrder.Order.Button == elevio.BT_Cab && newOrder.Id == (<-ch.MetaData).Id {
			acceptOrder(newOrder.Order, ch)
		} else {
			go func() {
				cost := calculateCost(newOrder.Order)
				costMsg := msgs.CostMsg{Cost: cost}
				for i := 0; i < 5; i++ {
					costMsg.Send()
					time.Sleep(1 * time.Millisecond)
				}
			}()
			costs := collectCosts(metaData.NumNodes)
			// Find minimum
			min := uint(math.Inf(0))
			minId := metaData.Id
			for id, cost := range costs {
				if cost < min {
					min = cost
					minId = id
				}
			}
			orderDiff := msgs.OrderTensorDiffMsg{
				Order: newOrder.Order,
				Diff:  msgs.DIFF_ADD,
				Id:    minId,
			}
			UpdateOrderTensor(ch.UpdateOrderTensor, ch.CurrentOrderTensor, orderDiff)
			// Needs to ensure that order is taken, if min is itself
			if minId == metaData.Id || min == costs[metaData.Id] {
				acceptOrder(newOrder.Order, ch)
			}
		}
	}
}

func checkForAcceptedOrders(ch c.Channels) {
	acceptedOrder := msgs.OrderTensorDiffMsg{}
	for {
		err := acceptedOrder.Listen()
		if err != nil {
			continue
		}
		if acceptedOrder.Diff == msgs.DIFF_ADD {
			fmt.Println("New accepted order")
			go func() {
				order := acceptedOrder.Order
				id := acceptedOrder.Id

				isNotDone := make(chan bool, 1)
				isDone := make(chan bool, 1)

				go sv.WatchOrder(isNotDone)
				go listenForOrderCompleted(order, isDone, isNotDone)
				select {
				case <-isNotDone:
					// Propagate the signal
					isNotDone <- true
					fmt.Printf("Order %+v is not done\n", order)
					// isDone <- false
					orderMsg := msgs.OrderMsg{
						Order: order,
						Id:    id,
					}
					delegateOrder(orderMsg, ch) //redelegate
					fmt.Println("delegated again")
				case <-isDone:
					// Propagate the signal
					isDone <- true
					fmt.Printf("Order %+v is done\n", order)
					// Nothing to do
				}
			}()
		}
	}
}

func acceptOrder(order elevio.ButtonEvent, ch c.Channels) {
	ch.TakeExternalOrder <- order
	orderTensorDiffMsg := msgs.OrderTensorDiffMsg{
		Order: order,
		Diff:  msgs.DIFF_ADD,
		Id:    (<-ch.MetaData).Id,
	}
	UpdateOrderTensor(ch.UpdateOrderTensor, ch.CurrentOrderTensor, orderTensorDiffMsg)
	for i := 0; i < 5; i++ {
		orderTensorDiffMsg.Send()
		time.Sleep(1 * time.Millisecond)
	}
}

func equalOrders(order1, order2 elevio.ButtonEvent) bool {
	return order1.Floor == order2.Floor && order1.Button == order2.Button
}
