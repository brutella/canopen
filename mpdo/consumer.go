package mpdo

import (
	"github.com/FabianPetersen/can"
	"github.com/FabianPetersen/canopen"
)

type Consumer struct {
	ObjectIndex canopen.ObjectIndex

	ObserveCobID uint16
	ReceiveCobID uint8
}

func (consumer *Consumer) Do(bus *can.Bus, channel chan [4]byte) error {
	bus.SetPassFilter([]uint32{uint32(consumer.ObserveCobID)})
	bus.SubscribeFunc(func(frame can.Frame) {
		// Check if the objectIndex is a match
		if frame.Data[0] == consumer.ReceiveCobID && frame.Data[1] == consumer.ObjectIndex.Index.B0 && frame.Data[2] == consumer.ObjectIndex.Index.B1 && frame.Data[3] == consumer.ObjectIndex.SubIndex {
			channel <- [4]byte{
				frame.Data[4],
				frame.Data[5],
				frame.Data[6],
				frame.Data[7],
			}
		}
	})
	return bus.ConnectAndPublish()
}
