package canopen

import (
	"fmt"
	"github.com/FabianPetersen/can"
	"github.com/jpillora/maplock"
	"time"
)

var Lock = maplock.New()

type TransferAbort struct{}

func (e TransferAbort) Error() string {
	return "Server aborted upload"
}

type UnexpectedSCSResponse struct {
	Expected uint8
	Actual   uint8
}

func (e UnexpectedSCSResponse) Error() string {
	return fmt.Sprintf("unexpected server command specifier %X (expected %X)", e.Actual, e.Expected)
}

type UnexpectedResponseLength struct {
	Expected int
	Actual   int
}

func (e UnexpectedResponseLength) Error() string {
	return fmt.Sprintf("unexpected response length %X (expected %X)", e.Actual, e.Expected)
}

type UnexpectedToggleBit struct {
	Expected bool
	Actual   bool
}

func (e UnexpectedToggleBit) Error() string {
	return fmt.Sprintf("unexpected toggle bit %t (expected %t)", e.Actual, e.Expected)
}

// A Client handles message communication by sending a request
// and waiting for the response.
type Client struct {
	Bus     *can.Bus
	Timeout time.Duration
}

// Do sends a request and waits for a response.
// If the response frame doesn't arrive on time, an error is returned.
func (c *Client) Do(req *Request) (*Response, error) {
	rch := can.Wait(c.Bus, req.ResponseID, c.Timeout)

	if err := c.Bus.Publish(req.Frame.CANFrame()); err != nil {
		return nil, err
	}

	resp := <-rch

	return &Response{CANopenFrame(resp.Frame), req}, resp.Err
}
