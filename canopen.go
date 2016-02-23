package canopen

// NodeID is the type of a CANopen frame id.
type NodeID uint8

// FunctionCode is the type of the CANopen frame type
type FunctionCode uint16

const (
	NmtService FunctionCode = 0x000
	Sync       FunctionCode = 0x080
	Emergency  FunctionCode = 0x080
	Timestamp  FunctionCode = 0x100
    TPDO1      FunctionCode = 0x180
    RPDO1      FunctionCode = 0x200
    TPDO2      FunctionCode = 0x280
    RPDO2      FunctionCode = 0x300
    TPDO3      FunctionCode = 0x380
    RPDO3      FunctionCode = 0x400
    TPDO4      FunctionCode = 0x480
    RPDO4      FunctionCode = 0x500
    TSDO       FunctionCode = 0x580
    RSDO       FunctionCode = 0x600
	Heartbeat  FunctionCode = 0x700
)

const MaxNodeID uint8 = 0x7F

const (
	CobMaskID        = 0x7FF // 11-bit cob identifier
	NodeMaskID       = 0x7F  // 7-bit identifier
	FunctionCodeMask = 0x780 // 4-bit function code
)
