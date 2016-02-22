package canopen

import (
	"github.com/brutella/can"
)

func SendFrame(frame Frame, wr can.Writer) error {
	return wr.WriteFrame(frame.CANFrame())
}

func GetFrame(r can.Reader) (frame Frame, err error) {
	canFrame := can.Frame{}

	err = r.ReadFrame(&canFrame)
	if err != nil {
		return
	}

	return CANopenFrame(canFrame), nil
}
