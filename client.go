package canopen

import (
	"github.com/brutella/can"
	"time"
)

// A Client handles message communication by sending a request
// and waiting for the response.
type Client struct {
	Bus     *can.Bus
	Timeout time.Duration
}

// Do sends a request and waits for a response.
// If the response frame doesn't arrive on time, an error is returned.
func (c *Client) Do(req *Request) (*Response, error) {
	respChan := make(chan can.WaitResponse)
	go func() {
		response := <-can.Wait(c.Bus, req.ResponseID, c.Timeout)
		respChan <- response
	}()
	time.Sleep(time.Millisecond * 15)
	if err := c.Bus.Publish(req.Frame.CANFrame()); err != nil {
		return nil, err
	}
	resp := <-respChan
	return &Response{CANopenFrame(resp.Frame), req}, resp.Err
}
