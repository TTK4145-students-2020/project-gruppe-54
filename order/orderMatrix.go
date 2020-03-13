package order

import (
	"fmt"
)

func initOrderMatrix(numNodes int, numFloors int) (<-chan [][]int,  chan<- [][]int) {
	rows := numFloors * 3
	coloumns := numNodes
	orderMatrix := make([][]int, coloumns)
	for i := 0; i < coloumns; i++ {
		orderMatrix[i] = make([]int, rows)
		for j := 0; j < coloumns; j++ {
			orderMatrix[i][j] = 0
		}

	}
	return orderMatrix
}

func orderMatrixServer([][]int orderMatrix, updateOrderMatrix <-chan [][]int, currentOrderMatrix  chan<- [][]int) {
	
}

func printOrderMatrix(orderMatrix [][]int) {
	counter := 1
	//orderTypes := [3]string{"Down: ", "Up: ", "Cab: "}
	for _, x := range orderMatrix {
		fmt.Print("Heis nr. ", counter, ": [ ")
		counter++
		//counter2 := 0
		for _, y := range x {
			//fmt.Print(orderTypes[counter2])adf
			fmt.Print(y, " ")
		}
		fmt.Print("]\n")
	}
}
