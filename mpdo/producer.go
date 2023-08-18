package mpdo

import (
	"github.com/FabianPetersen/can"
	"github.com/FabianPetersen/canopen"
	"strconv"
)

// Producer represents a MPDO send message
type Producer struct {
	ObjectIndex canopen.ObjectIndex

	Data         [4]byte
	RequestCobID uint16
	ReceiveCobID uint8
}

func (producer Producer) Do(bus *can.Bus) error {
	// Do not allow multiple messages for the same device
	key := strconv.Itoa(int(producer.ReceiveCobID))
	canopen.Lock.Lock(key)
	defer canopen.Lock.Unlock(key)

	return bus.Publish(can.Frame{
		ID: uint32(producer.RequestCobID),
		Data: [8]byte{
			producer.ReceiveCobID,
			producer.ObjectIndex.Index.B0, producer.ObjectIndex.Index.B1,
			producer.ObjectIndex.SubIndex,
			producer.Data[0],
			producer.Data[1],
			producer.Data[2],
			producer.Data[3],
		},
	})
}
