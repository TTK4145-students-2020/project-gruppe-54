package order

type FloorOrders struct {
	Inside,
	OutsideUp,
	OutsideDown bool
}

type Node struct {
	Floor []FloorOrders
}

func initOrderTensor(numNodes int, numFloors int) (<-chan []Node, chan<- []Node) {
	orderTensor := make([]Node, numNodes)
	for i := 0; i < 4; i++ {
		orderTensor[i].Floor = make([]FloorOrders, numFloors)
		for floor := 0; floor < 3; floor++ {
			orderTensor[i].Floor[floor] = FloorOrders{
				Inside:      false,
				OutsideUp:   false,
				OutsideDown: false,
			}
		}
	}
	updateOrderTensor := make(<-chan []Node, 1)
	currentOrderTensor := make(chan<- []Node, 1)
	go orderTensorServer(orderTensor, updateOrderTensor, currentOrderTensor)
	return updateOrderTensor, currentOrderTensor
}

func orderTensorServer(orderTensor []Node, updateOrderTensor <-chan []Node, currentOrderTensor chan<- []Node) {
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
