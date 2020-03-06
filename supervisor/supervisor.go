package supervisor

import (
	"time"

	"./watchdog"
)

var timeoutPeriod int = 5

//WatchOrder ... watches an order
func WatchOrder(isNotDone chan bool) {
	wd := watchdog.NewWatchdog(time.Second * time.Duration(timeoutPeriod))
	select { // whichever one of order complete or wd finished first
	case <-wd.GetKickChannel():
		isNotDone <- true
		wd.Stop()
		return
	}
}

func main() {
	//stop := make(chan bool)
	//go NewOrder(a)
	//for {

	//}
	//stop <- true
	//<-stop
}
