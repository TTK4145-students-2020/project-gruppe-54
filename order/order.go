package order

//OrderStruct ... orderstruct
type OrderStruct struct {
	Floor int
	
}


//state matrix



//func orderComplete

//func receiveOrder ...
	//make chan order
	//start watchdog

	//order complete


func delegateOrders() {
	orders := make(chan elevio.ButtonEvent)
	for {
		select {
			case newOrder := <-orders:
				chooseElevator()
				
		}
	}
}
func chooseElevator() int{
	GetCostFuction

}