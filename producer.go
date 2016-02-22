package canopen

import (
	"github.com/brutella/can"
	"time"
)

func Produce(frame Frame, bus *can.Bus, timeout time.Duration) chan<- struct{} {
	stop := make(chan struct{})
	canFrame := frame.CANFrame()
	go func() {
		for {
			if err := bus.Publish(canFrame); err != nil {
				return
			}

			select {
			case <-stop:
				return
			case <-time.After(timeout):
				break
			}
		}
	}()
	return stop
}

func ProduceHeartbeat(id NodeID, state NMTState, bus *can.Bus, timeout time.Duration) chan<- struct{} {
	frame := NewHeartbeatFrame(id, state)
	return Produce(frame, bus, timeout)
}
