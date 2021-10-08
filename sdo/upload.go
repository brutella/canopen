package sdo

import (
	"github.com/brutella/can"
	"github.com/brutella/canopen"

	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	UploadInitiateRequest  = 0x40 // 0100 0000
	UploadInitiateResponse = 0x40 // 0100 0000

	UploadSegmentRequest  = 0x60 // 0110 0000
	UploadSegmentResponse = 0x00 // 0000 0000
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
			byte(UploadInitiateRequest),
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
	switch scs := frame.Data[0] >> 5; scs {
	case 2:
		break
	case TransferAbort:
		return nil, errors.New("Server aborted upload")
	default:
		log.Fatalf("Unexpected server command specifier %X", scs)
	}

	if hasBit(frame.Data[0], 1) { // e = 1?
		// number of segment bytes with no data
		var n uint8
		if hasBit(frame.Data[0], 0) { // s = 1?
			n = (frame.Data[0] >> 2) & 0x3
		}
		return frame.Data[4 : 8-n], nil
	}

	// Read segment data length
	var n uint32
	b := bytes.NewBuffer(frame.Data[4:8])
	if err := binary.Read(b, binary.LittleEndian, &n); err != nil {
		return nil, err
	}

	var i int
	var buf bytes.Buffer
	for {
		data := make([]byte, 8)

		// ccs = 3
		data[0] |= 3 << 5

		if i%2 == 1 {
			// t = 1
			data[0] = setBit(data[0], 4)
		}

		i += 1

		frame = canopen.Frame{
			CobID: upload.RequestCobID,
			Data:  data,
		}

		req = canopen.NewRequest(frame, uint32(upload.ResponseCobID))
		resp, err = c.Do(req)
		if err != nil {
			return nil, err
		}

		if hasBit(frame.Data[0], 4) != hasBit(resp.Frame.Data[0], 4) {
			return nil, fmt.Errorf("unexpected toggle bit %t", hasBit(resp.Frame.Data[0], 4))
		}

		n := (resp.Frame.Data[0] >> 1) & 0x7
		buf.Write(resp.Frame.Data[1 : 8-n])

		if hasBit(resp.Frame.Data[0], 0) { // c = 1?
			break
		}
	}

	return buf.Bytes(), nil
}
