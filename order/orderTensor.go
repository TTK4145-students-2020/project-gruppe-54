package order

type FloorOrders struct {
	Inside,
	OutsideUp,
	OutsideDown bool
}

type Node struct {
	Floor []FloorOrders
}

type UpdateOrderTensorChan <-chan []Node
type CurrentOrderTensorChan chan<- []Node

func InitOrderTensor(numNodes int, numFloors int) (UpdateOrderTensorChan, CurrentOrderTensorChan) {
	orderTensor := make([]Node, numNodes)
	for i := 0; i < numNodes; i++ {
		orderTensor[i].Floor = make([]FloorOrders, numFloors)
		for floor := 0; floor < numFloors; floor++ {
			orderTensor[i].Floor[floor] = FloorOrders{
				Inside:      false,
				OutsideUp:   false,
				OutsideDown: false,
			}
		}
	}
	updateOrderTensor := make(UpdateOrderTensorChan, 1)
	currentOrderTensor := make(CurrentOrderTensorChan, 1)
	go orderTensorServer(orderTensor, updateOrderTensor, currentOrderTensor)
	return updateOrderTensor, currentOrderTensor
}

func orderTensorServer(orderTensor []Node, updateOrderTensor UpdateOrderTensorChan, currentOrderTensor CurrentOrderTensorChan) {
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
