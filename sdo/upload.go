package sdo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/FabianPetersen/can"
	"github.com/FabianPetersen/canopen"
	"log"
	"time"
)

const (
	ClientIntiateUpload = 0x40 // 0100 0000
	ClientSegmentUpload = 0x60 // 0110 0000

	ServerInitiateUpload = 0x40 // 0100 0000
	ServerSegmentUpload  = 0x00 // 0000 0000
)

// Upload represents a SDO upload process to read data from a CANopen
// device – upload because the receiving node uploads data to another node.
type Upload struct {
	ObjectIndex canopen.ObjectIndex

	RequestCobID  uint16
	ResponseCobID uint16
}

func (upload Upload) Do(bus *can.Bus) ([]byte, error) {
	c := &canopen.Client{bus, time.Second * 2}
	// Initiate
	frame := canopen.Frame{
		CobID: upload.RequestCobID,
		Data: []byte{
			byte(ClientIntiateUpload),
			upload.ObjectIndex.Index.B0, upload.ObjectIndex.Index.B1,
			upload.ObjectIndex.SubIndex,
			0x0, 0x0, 0x0, 0x0,
		},
	}

	req := canopen.NewRequest(frame, uint32(upload.ResponseCobID))
	resp, err := c.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	frame = resp.Frame
	b0 := frame.Data[0] // == 0100 nnes
	scs := b0 & TransferMaskCommandSpecifier
	switch scs {
	case ServerInitiateUpload:
		break
	case TransferAbort:
		return nil, errors.New("Server aborted upload")
	default:
		log.Fatalf("Unexpected server command specifier %X", scs)
	}

	if isExpedited(frame) {
		// number of segment bytes with no data
		n := 0
		if isSizeIndicated(frame) {
			n = sizeValue(frame)
		}
		return frame.Data[4 : 8-n], nil
	}

	// Read segment data length
	var n uint32
	b := bytes.NewBuffer(frame.Data[4:8])
	if err := binary.Read(b, binary.LittleEndian, &n); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	t := 0
	for {
		// Upload segment
		cmd := byte(ClientSegmentUpload)
		if t%2 == 1 {
			cmd |= TransferSegmentToggle
		}

		t += 1

		frame = canopen.Frame{
			CobID: upload.RequestCobID,
			Data: []byte{
				cmd,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			},
		}

		req = canopen.NewRequest(frame, uint32(upload.ResponseCobID))
		resp, err = c.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		// Segment response
		frame := resp.Frame
		b0 = frame.Data[0]

		n := int((b0 & 0xE) >> 1)
		buf.Write(resp.Frame.Data[1 : 8-n])
		if isLast(frame) {
			break
		}
	}

	return buf.Bytes(), nil
}

func isExpedited(frame canopen.Frame) bool {
	return frame.Data[0]&TransferExpedited == TransferExpedited
}

func isSizeIndicated(frame canopen.Frame) bool {
	return frame.Data[0]&TransferSizeIndicated == TransferSizeIndicated
}

func isLast(frame canopen.Frame) bool {
	return frame.Data[0]&0x1 == 0x1
}

func sizeValue(frame canopen.Frame) int {
	return int(frame.Data[0] & TransferMaskSize >> 2)
}
