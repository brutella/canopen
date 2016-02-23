package sdo

var CommandSpecifierMask byte = 0xE0

const (
	TransferExpedited     = 0x2
	TransferSizeIndicated = 0x1
	TransferSizeMask      = 0xC
	TransferSegmentToggle = 0x10 // 0001 0000
	TransferAbort         = 0x80
)
