package order

import (
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
	orderComplete := false

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
				orderComplete = true
				return
			}
		case <-isNotDone:
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
	costMsg := msgs.CostMsg{}
	go func() {
		for {
			select {
			case <-timer.C:
				stopPolling <- true
				close(stopPolling)
				return
			default:
				err := costMsg.Listen()
				if err != nil {
					continue
				}
				newCost <- costMsg
			}
		}
	}()
L:
	for {
		select {
		case <-stopPolling:
			break L
		case newCostMsg := <-newCost:
			if id := newCostMsg.Id; id >= 0 && id < numNodes {
				costs[id] = newCostMsg.Cost
			}
		}
	}
	return costs
}

//ControlOrders ... Delegates new orders it receives on channel newOrders
func ControlOrders(ch c.Channels) {
	newOrders := make(chan msgs.OrderMsg, 1000)
	go handleNewOrder(newOrders, ch)
	go listenForNewOrders(newOrders, ch)
	go checkForAcceptedOrders(newOrders, ch)
	for {
		select {
		case newOrder := <-ch.DelegateOrder:
			orderMsg := msgs.OrderMsg{Order: newOrder}
			orderMsg.Id = (<-ch.MetaData).Id
			delegateOrder(orderMsg, ch)
			newOrders <- orderMsg
		case orderCompleted := <-ch.CompletedOrder: // the external order has been taken
			orderTensorDiffMsg := msgs.OrderTensorDiffMsg{
				Order: orderCompleted,
				Diff:  msgs.DIFF_REMOVE,
				Id:    (<-ch.MetaData).Id}
			for i := 0; i < 5; i++ {
				orderTensorDiffMsg.Send()
				time.Sleep(1 * time.Millisecond)
			}
		}
	}
}

func listenForNewOrders(newOrders chan msgs.OrderMsg, ch c.Channels) {
	newOrder := msgs.OrderMsg{}
	for {
		err := newOrder.Listen()
		if err != nil {
			continue
		} else {
			newOrders <- newOrder
		}
	}
}

func handleNewOrder(newOrders chan msgs.OrderMsg, ch c.Channels) {
	metaData := <-ch.MetaData
	for {
		newOrder := <-newOrders
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
			min := uint(math.Inf(0))
			minID := metaData.Id
			for id, cost := range costs {
				if cost < min {
					min = cost
					minID = id
				}
			}
			if minID == metaData.Id || min == costs[metaData.Id] {
				acceptOrder(newOrder.Order, ch)
			}
		}
	}
}

func checkForAcceptedOrders(newOrders chan msgs.OrderMsg, ch c.Channels) {
	acceptedOrder := msgs.OrderTensorDiffMsg{}
	i := 0
	for {
		i++
		err := acceptedOrder.Listen()
		if err != nil {
			continue
		}
		if acceptedOrder.Diff == msgs.DIFF_ADD {
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
					orderMsg := msgs.OrderMsg{
						Order: order,
						Id:    id,
					}
					delegateOrder(orderMsg, ch)
					newOrders <- orderMsg //redelegate
				case <-isDone:
					// Propagate the signal
					isDone <- true
				}
			}()
		}
	}
}

func acceptOrder(order elevio.ButtonEvent, ch c.Channels) {
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
	ch.TakeExternalOrder <- order
}

func equalOrders(order1, order2 elevio.ButtonEvent) bool {
	return order1.Floor == order2.Floor && order1.Button == order2.Button
}
