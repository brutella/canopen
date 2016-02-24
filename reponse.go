package canopen

// A Response represents the response which resulted from a request.
type Response struct {
	// The response frame
	Frame Frame

	// The Request which triggers the response
	Request *Request
}
