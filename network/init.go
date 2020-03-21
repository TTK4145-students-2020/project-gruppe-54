package network

import (
	"./msgs"
)

func InitNetwork() error {
	err := msgs.InitMessages()
	return err
}
