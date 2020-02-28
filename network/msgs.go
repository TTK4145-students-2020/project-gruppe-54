package network

type RecCh interface {
}

type CostMsg struct {
	channel chan CostMsg
	Cost    float64
	ID      int
}
