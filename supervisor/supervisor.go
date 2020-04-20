package supervisor

import (
	"time"

	"./watchdog"
)

var timeoutPeriod time.Duration = 10 * time.Second

//WatchOrder ... watches an order
func WatchOrder(isNotDone chan bool) {
	wd := watchdog.NewWatchdog(timeoutPeriod)
	select {
	case <-wd.GetKickChannel():
		isNotDone <- true
		wd.Stop()
		return
	}
}
