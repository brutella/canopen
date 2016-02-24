package canopen

const (
	MessageTypeNMT       uint16 = 0x000
	MessageTypeSync      uint16 = 0x080
	MessageTypeTimestamp uint16 = 0x100
	MessageTypeTPDO1     uint16 = 0x180
	MessageTypeRPDO1     uint16 = 0x200
	MessageTypeTPDO2     uint16 = 0x280
	MessageTypeRPDO2     uint16 = 0x300
	MessageTypeTPDO3     uint16 = 0x380
	MessageTypeRPDO3     uint16 = 0x400
	MessageTypeTPDO4     uint16 = 0x480
	MessageTypeRPDO4     uint16 = 0x500
	// MessageTypeTSDO represents the type of SDO server response messages
	MessageTypeTSDO uint16 = 0x580
	// MessageTypeRSDO represents the type of SDO client request messages
	MessageTypeRSDO      uint16 = 0x600
	MessageTypeHeartbeat uint16 = 0x700
)

// MaxNodeID defines the highest node id
const MaxNodeID uint8 = 0x7F

const (
	// MaskCobID is used to get 11 bits from an uint16 for the COB-ID
	MaskCobID = 0x7FF
	// MaskNodeID is used to extract the 7-bit node id from the COB-ID
	MaskNodeID = 0x7F
	// MaskMessageType is used to extract the 4-bit message type from the COB-ID
	MaskMessageType = 0x780
)
