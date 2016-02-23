package canopen

// NMTCommand is the type of a NMT command.
type NMTCommand byte

const (
	GoToOperational        NMTCommand = 0x1
	GoToStopped            NMTCommand = 0x2
	GoToPreOperation       NMTCommand = 0x80
	GoToResetNode          NMTCommand = 0x81
	GoToResetCommunication NMTCommand = 0x82
)

// NMTState is the type of a NMT state.
type NMTState byte

const (
	BootUp         NMTState = 0x0
	Stopped        NMTState = 0x04
	Operational    NMTState = 0x05
	PreOperational NMTState = 0x7f
)
