package canopen

const (
	// MessageTypeNMT represents the type of network management messages
	MessageTypeNMT uint16 = 0x000
	// MessageTypeSync represents the type of synchronous messages
	MessageTypeSync uint16 = 0x080
	// MessageTypeEmergency represents the type of emergency messages
	MessageTypeEmergency uint16 = 0x080
	// MessageTypeTimestamp represents the type of timestamp messages
	MessageTypeTimestamp uint16 = 0x100
	// MessageTypeTPDO1 represents the type of TPDO1 messages
	MessageTypeTPDO1 uint16 = 0x180
	// MessageTypeRPDO1 represents the type of RPDO1 messages
	MessageTypeRPDO1 uint16 = 0x200
	// MessageTypeTPDO2 represents the type of TPDO2 messages
	MessageTypeTPDO2 uint16 = 0x280
	// MessageTypeRPDO2 represents the type of RPDO2 messages
	MessageTypeRPDO2 uint16 = 0x300
	// MessageTypeTPDO3 represents the type of TPDO3 messages
	MessageTypeTPDO3 uint16 = 0x380
	// MessageTypeRPDO3 represents the type of RPDO3 messages
	MessageTypeRPDO3 uint16 = 0x400
	// MessageTypeTPDO4 represents the type of TPDO4 messages
	MessageTypeTPDO4 uint16 = 0x480
	// MessageTypeRPDO4 represents the type of RPDO4 messages
	MessageTypeRPDO4 uint16 = 0x500
	// MessageTypeTSDO represents the type of SDO server response messages
	MessageTypeTSDO uint16 = 0x580
	// MessageTypeRSDO represents the type of SDO client request messages
	MessageTypeRSDO uint16 = 0x600
	// MessageTypeHeartbeat represents the type of heartbeat messages
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
