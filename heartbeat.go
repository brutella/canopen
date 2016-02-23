package canopen

// NewHeartbeatFrame returns a new CANopen heartbeat frame containing the node id and state. 
func NewHeartbeatFrame(id NodeID, state NMTState) Frame {
	return Frame{
		CobID: uint16(Heartbeat) + uint16(id),
		Data:  []byte{byte(state)},
	}
}