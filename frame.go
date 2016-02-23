package canopen

import (
	"fmt"
	"github.com/brutella/can"
)

// A Frame represents a CANopen frame.
type Frame struct {
	// CobID is the 11-bit communication object identifier – CANopen only uses 11-bit identifiers.
	// Bits 0-6 represent the 7-bit node ID. Bits 7-11 represent the 4-bit message type.
	CobID uint16
	// Rtr represents the Remote Transmit Request flag.
	Rtr bool
	// Data contains 8 bytes
	Data []uint8
}

// CANopenFrame returns a CANopen frame from a CAN frame.
func CANopenFrame(frm can.Frame) Frame {
	canopenFrame := Frame{}

	canopenFrame.CobID = uint16(frm.ID & can.MaskIDSff)
	canopenFrame.Rtr = (frm.ID & can.MaskRtr) == can.MaskRtr
	canopenFrame.Data = frm.Data[:]

	return canopenFrame
}

// NewFrame returns a frame with an id and data bytes.
func NewFrame(id uint16, data []uint8) Frame {
	return Frame{
		CobID: id & MaskCobID, // only use first 11 bits
		Data:  data,
	}
}

func (frm Frame) String() string {
	t := frm.MessageType()
	id := frm.NodeID()
	switch t {
	case MessageTypeTimestamp:
		if time, _ := frm.Timestamp(); time != nil {
			return fmt.Sprintf("Timestamp: timestamp %s", time.String())
		}
	case MessageTypeHeartbeat:
		return fmt.Sprintf("Heartbeat: node #%d", id)
	}

	return fmt.Sprintf("Message type %X, node id %d, data: % X", t, id, frm.Data)
}

// MessageType returns the message type.
func (frm Frame) MessageType() uint16 {
	return frm.CobID & MaskMessageType
}

// NodeID returns the node id.
func (frm Frame) NodeID() uint8 {
	return uint8(frm.CobID & MaskNodeID)
}

// CANFrame returns a CAN frame representing the CANopen frame.
func (frm Frame) CANFrame() can.Frame {
	var data [8]uint8
	n := len(frm.Data)
	copy(data[:n], frm.Data[:n])

	// Convert CANOpen COB-ID to CAN id including RTR flag
	id := uint32(frm.CobID)
	if frm.Rtr == true {
		id = id | can.MaskRtr
	}

	return can.Frame{
		ID:     id,
		Length: uint8(len(frm.Data)),
		Data:   data,
	}
}

// Marshal returns the byte encoding of frm.
func Marshal(frm Frame) (b []byte, err error) {
	canFrm := frm.CANFrame()

	return can.Marshal(canFrm)
}

// Unmarshal parses the bytes b and stores the result in the value
// pointed to by frm.
func Unmarshal(b []byte, frm *Frame) error {
	canFrm := can.Frame{}
	if err := can.Unmarshal(b, &canFrm); err != nil {
		return err
	}

	*frm = CANopenFrame(canFrm)

	return nil
}
