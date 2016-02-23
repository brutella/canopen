package canopen

import (
	"fmt"
	"github.com/brutella/can"
)

// A Frame represents a CANopen frame.
type Frame struct {
	// CobID is the 11-bit communication object identifier – CANopen only uses 11-bit identifiers.
	// Bits 0-6 represent the 7-bit node ID. Bits 7-11 represent the 4-bit function code.
	CobID uint16
    // Rtr represents the Remote Transmit Request flag.
	Rtr   bool
    // Data contains 8 bytes
	Data  []uint8 
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
		CobID: id & CobMaskID, // only use first 11 bits
		Data:  data,
	}
}

func (frm Frame) String() string {
	/*
		IDNmtService uint16 = 0x000
		IDSync       uint16 = 0x080
		IDEmergency  uint16 = 0x080
		Timestamp  uint16 = 0x100
		IDTpdo1      uint16 = 0x180
		IDRpdo1      uint16 = 0x200
		IDTpdo2      uint16 = 0x280
		IDRpdo2      uint16 = 0x300
		IDTpdo3      uint16 = 0x380
		IDRpdo4      uint16 = 0x400
		IDTpdo5      uint16 = 0x480
		IDRpdo5      uint16 = 0x500
		IDTsdo       uint16 = 0x580
		IDRsdo       uint16 = 0x600
		IDHeartbeat  uint16 = 0x700
	*/
	fnc := frm.FunctionCode()
	nid := frm.NodeID()
	switch fnc {
	case Timestamp:
		if t, _ := frm.Timestamp(); t != nil {
			return fmt.Sprintf("Timestamp: timestamp %s", t.String())
		}
	case Heartbeat:
		return fmt.Sprintf("Heartbeat: node #%d", nid)
	}

	return fmt.Sprintf("Function code %X, node id %d, data: % X", fnc, nid, frm.Data)
}

// FunctionCode returns the function code of the frame.
func (frm Frame) FunctionCode() FunctionCode {
	return FunctionCode(frm.CobID & FunctionCodeMask)
}

// NodeID returns the node id of the frame.
func (frm Frame) NodeID() NodeID {
	return NodeID(frm.CobID & NodeMaskID)
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
