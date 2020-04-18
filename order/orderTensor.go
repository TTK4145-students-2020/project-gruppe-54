package order

import (
	conf "../configuration"
	"../hardware/driver-go/elevio"
	"../network/msgs"
)

func InitOrderTensor(numNodes int, numFloors int) (chan []conf.Node, chan []conf.Node) {
	orderTensor := make([]conf.Node, numNodes)
	for i := 0; i < numNodes; i++ {
		orderTensor[i].Floor = make([]conf.FloorOrders, numFloors)
		for floor := 0; floor < numFloors; floor++ {
			orderTensor[i].Floor[floor] = conf.FloorOrders{
				Inside:      false,
				OutsideUp:   false,
				OutsideDown: false,
			}
		}
	}
	updateOrderTensor := make(conf.UpdateOrderTensorChan, 1)
	currentOrderTensor := make(conf.CurrentOrderTensorChan, 1)
	go orderTensorServer(orderTensor, updateOrderTensor, currentOrderTensor)
	return updateOrderTensor, currentOrderTensor
}

func UpdateOrderTensor(updateOrderTensorCh chan<- []conf.Node, currentOrderTensorCh <-chan []conf.Node, orderDiff msgs.OrderTensorDiffMsg) {
	newTensor := <-currentOrderTensorCh
	id := orderDiff.Id
	order := orderDiff.Order

	var diff bool
	switch orderDiff.Diff {
	case msgs.DIFF_ADD:
		diff = true
	case msgs.DIFF_REMOVE:
		diff = false
	}

	switch order.Button {
	case elevio.BT_HallUp:
		newTensor[id].Floor[order.Floor].OutsideUp = diff
	case elevio.BT_HallDown:
		newTensor[id].Floor[order.Floor].OutsideDown = diff
	case elevio.BT_Cab:
		newTensor[id].Floor[order.Floor].Inside = diff
	}
	updateOrderTensorCh <- newTensor
}

func orderTensorServer(orderTensor []conf.Node, updateOrderTensor conf.UpdateOrderTensorChan, currentOrderTensor conf.CurrentOrderTensorChan) {
	for {
		select {
		case <-updateOrderTensor:
			orderTensor = <-updateOrderTensor
		default:
			currentOrderTensor <- orderTensor
		}
	}
}

// func printOrderMatrix(orderMatrix [][]int) {
// 	counter := 1
// 	//orderTypes := [3]string{"Down: ", "Up: ", "Cab: "}
// 	for _, x := range orderMatrix {
// 		fmt.Print("Heis nr. ", counter, ": [ ")
// 		counter++
// 		//counter2 := 0
// 		for _, y := range x {
// 			//fmt.Print(orderTypes[counter2])adf
// 			fmt.Print(y, " ")
// 		}
// 		fmt.Print("]\n")
// 	}
// }
