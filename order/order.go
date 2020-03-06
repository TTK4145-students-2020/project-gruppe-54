package order

import(
	"fmt"
	"../supervisor"
)

//OrderStruct ... orderstruct
type OrderStruct struct {
	Floor int
	
}


//state matrix



func listenForOrderCompleted(){
	
}

//func receiveOrder ...
	//make chan order
	//start watchdog

	//order complete

func sendOrder(elev int){
	go WatchOrder() // does everything by itself, no need to check on it.. resends order if it doesnt receive finishedoRder
}

func delegateOrder(order OrderStruct){
	chosenElev = lowestCost()
	sendOrder(chosenElev)
	WatchOrder()
}

func delegateOrders() {
	newOrders := make(chan elevio.ButtonEvent)
	for {
		select {
			case newOrder := <-newOrders:
				delegateOrder(newOrder)				
		}
	}
}


func receiveOrders(){

}

func mainOrderthing(){
	go delegateOrders()
	go receiveOrders()
}