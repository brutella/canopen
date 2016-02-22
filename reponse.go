package canopen

// Response represents the response from a request.
type Response struct {
	// The Frame of the response
	Frame Frame

	// The Request that was sent to obtain this Response
	Request *Request
}
