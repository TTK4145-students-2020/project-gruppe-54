package network

import (
	c "../configuration"
	"./msgs"
)

func InitNetwork(metaDataChan <-chan c.MetaData) error {
	err := msgs.InitMessages(metaDataChan)
	return err
}
