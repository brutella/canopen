package canopen

import (
	"github.com/FabianPetersen/can"
	"github.com/jpillora/maplock"
	"strconv"
	"time"
)

var m = maplock.New()

// A Client handles message communication by sending a request
// and waiting for the response.
type Client struct {
	Bus     *can.Bus
	Timeout time.Duration
}

// Do sends a request and waits for a response.
// If the response frame doesn't arrive on time, an error is returned.
func (c *Client) Do(req *Request) (*Response, error) {
	// Do not allow multiple messages for the same device
	key := strconv.Itoa(int(req.ResponseID))
	m.Lock(key)
	defer m.Unlock(key)

	rch := can.Wait(c.Bus, req.ResponseID, c.Timeout)

	if err := c.Bus.Publish(req.Frame.CANFrame()); err != nil {
		return nil, err
	}

	resp := <-rch

	return &Response{CANopenFrame(resp.Frame), req}, resp.Err
}
