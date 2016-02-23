package canopen

// A Request represents a CANopen request published on a CAN bus and received by another CANopen node.
type Request struct {
	// The Frame of the request
	Frame Frame

	// The ResponseID of the response frame
	ResponseID uint32
}

// NewRequest returns a request containing the frame to be sent 
// and the expected response frame id.
func NewRequest(frm Frame, respID uint32) *Request {
	return &Request{
		Frame:      frm,
		ResponseID: respID,
	}
}
