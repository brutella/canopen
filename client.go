package canopen

import (
	"github.com/brutella/can"
	"time"
)

type Client struct {
	Bus     *can.Bus
	Timeout time.Duration
}

func (c *Client) Do(req *Request) (*Response, error) {
	if err := c.Bus.Publish(req.Frame.CANFrame()); err != nil {
		return nil, err
	}

	resp := <-can.Wait(c.Bus, req.ResponseID, c.Timeout)

	return &Response{CANopenFrame(resp.Frame), req}, resp.Err
}
