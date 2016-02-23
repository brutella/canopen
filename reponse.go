package canopen

// A Response represents the response from a request.
type Response struct {
	// The frame of the response
	Frame Frame

	// The Request which triggers the response.
	Request *Request
}
