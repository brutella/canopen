package canopen

import (
	"github.com/FabianPetersen/can"
	"testing"
	"time"
)

var testFrame = can.Frame{
	ID:     0xAF,
	Length: 0x8,
	Flags:  0x0,
	Res0:   0x0,
	Res1:   0x0,
	Data:   [can.MaxFrameDataLength]uint8{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8},
}

func TestClientRequestResponse(t *testing.T) {
	bus := can.NewBus(can.NewEchoReadWriteCloser())

	go bus.ConnectAndPublish()
	defer bus.Disconnect()

	req := NewRequest(CANopenFrame(testFrame), testFrame.ID)
	c := &Client{bus, time.Second * 1}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if is, want := resp.Frame.CANFrame().ID, testFrame.ID; is != want {
		t.Fatalf("is=%v want=%v", is, want)
	}

	if is, want := resp.Request, req; is != want {
		t.Fatalf("is=%v want=%v", is, want)
	}
}
