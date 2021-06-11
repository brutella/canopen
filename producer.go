package canopen

import (
	"github.com/FabianPetersen/can"
	"time"
)

// Produce repeatedly sends a frame on a bus after a timeout.
// The sending can be stopped by using the returned channel.
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

// ProduceHeartbeat repeatedly sends a CANopen heartbeat frame.
func ProduceHeartbeat(nodeID uint8, state uint8, bus *can.Bus, timeout time.Duration) chan<- struct{} {
	frame := NewHeartbeatFrame(nodeID, state)
	return Produce(frame, bus, timeout)
}
