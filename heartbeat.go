package canopen

// NewHeartbeatFrame returns a new CANopen heartbeat frame containing the node id and state.
func NewHeartbeatFrame(id uint8, state uint8) Frame {
	return Frame{
		CobID: uint16(MessageTypeHeartbeat) + uint16(id),
		Data:  []byte{byte(state)},
	}
}
