package canopen

type NMTCommand byte

const (
	GoToOperational        NMTCommand = 0x1
	GoToStopped            NMTCommand = 0x2
	GoToPreOperation       NMTCommand = 0x80
	GoToResetNode          NMTCommand = 0x81
	GoToResetCommunication NMTCommand = 0x82
)

type NMTState byte

const (
	BootUp         NMTState = 0x0
	Stopped        NMTState = 0x04
	Operational    NMTState = 0x05
	PreOperational NMTState = 0x7f
)

func NewHeartbeatFrame(id NodeID, state NMTState) Frame {
	return Frame{
		CobID: uint16(Heartbeat) + uint16(id),
		Data:  []byte{byte(state)},
	}
}
