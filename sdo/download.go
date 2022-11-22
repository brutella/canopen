package sdo

import (
	"github.com/FabianPetersen/can"
	"github.com/FabianPetersen/canopen"
	"strconv"

	"bytes"
	"encoding/binary"
	"time"
)

const (
	DownloadInitiateRequest  = 0x20 // 0010 0000
	DownloadInitiateResponse = 0x60 // 0110 0000

	DownloadSegmentRequest  = 0x00 // 0000 0000
	DownloadSegmentResponse = 0x20 // 0010 0000
)

// Download represents a SDO download process to write data to a CANopen
// device – download because the receiving node downloads data.
type Download struct {
	ObjectIndex canopen.ObjectIndex

	Data          []byte
	RequestCobID  uint16
	ResponseCobID uint16
}

func (download Download) Do(bus *can.Bus) error {
	// Do not allow multiple messages for the same device
	key := strconv.Itoa(int(download.RequestCobID))
	canopen.Lock.Lock(key)
	defer canopen.Lock.Unlock(key)

	if err := download.doInit(bus); err != nil {
		return err
	}

	return download.doSegments(bus)
}

func (download Download) doInit(bus *can.Bus) error {
	frame, err := download.initFrame()
	if err != nil {
		return err
	}

	req := canopen.NewRequest(frame, uint32(download.ResponseCobID))
	c := &canopen.Client{Bus: bus, Timeout: time.Second * 2}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	frame = resp.Frame
	switch scs := frame.Data[0] >> 5; scs {
	case 3: // response
		// Check if this is the correct response for the requested message
		if frame.Data[1] != download.ObjectIndex.Index.B0 || frame.Data[2] != download.ObjectIndex.Index.B1 || frame.Data[3] != download.ObjectIndex.SubIndex {
			return canopen.TransferAbort{
				AbortCode: getAbortCodeBytes(frame),
			}
		}
		return nil
	case 4: // abort
		return canopen.TransferAbort{
			AbortCode: getAbortCodeBytes(frame),
		}
	default:
		return canopen.UnexpectedSCSResponse{
			Expected:  3,
			Actual:    scs,
			AbortCode: getAbortCodeBytes(frame),
		}
	}
}

// initFrame returns the initial frame of the download.
// If the download data is less than 4 bytes, the init frame data contains all download data.
// If the download data is more than 4 bytes, the init frame data contains the overall length of the download data.
func (download Download) initFrame() (frame canopen.Frame, err error) {
	fdata := make([]byte, 4)

	// css = 1 (download init request)
	fdata[0] = setBit(fdata[0], 5)
	fdata[1] = download.ObjectIndex.Index.B0
	fdata[2] = download.ObjectIndex.Index.B1
	fdata[3] = download.ObjectIndex.SubIndex

	n := len(download.Data)
	if n <= 4 { // does download data fit into one frame?
		// e = 1 (expedited)
		fdata[0] = setBit(fdata[0], 1)
		// s = 1
		fdata[0] = setBit(fdata[0], 0)

		// n = number of unused bytes in frame.Data
		emptyBytes := 4 - n
		if emptyBytes == 2 || emptyBytes == 3 {
			fdata[0] = setBit(fdata[0], 3)
		}
		if emptyBytes == 1 || emptyBytes == 3 {
			fdata[0] = setBit(fdata[0], 2)
		}

		// copy all download data into frame data
		fdata = append(fdata, download.Data...)
	} else {
		// e = 0
		// n = 0 (frame.Data contains the overall )
		// s = 1
		fdata[0] = setBit(fdata[0], 0)

		var buf bytes.Buffer
		if err = binary.Write(&buf, binary.LittleEndian, uint32(n)); err != nil {
			return
		}

		// copy overall length of download data into frame data
		fdata = append(fdata, buf.Bytes()...)
	}

	// CiA301 Standard expects all (8) bytes to be sent
	for len(fdata) < 8 {
		fdata = append(fdata, 0x0)
	}

	frame.CobID = download.RequestCobID
	frame.Data = fdata

	return
}

func (download Download) doSegments(bus *can.Bus) error {
	frames := download.segmentFrames()

	c := &canopen.Client{Bus: bus, Timeout: time.Second * 2}
	for _, frame := range frames {
		req := canopen.NewRequest(frame, uint32(download.ResponseCobID))
		resp, err := c.Do(req)
		if err != nil {
			return err
		}

		switch scs := resp.Frame.Data[0] >> 5; scs {
		case 1:
			break
		case 4:
			return canopen.TransferAbort{
				AbortCode: getAbortCodeBytes(resp.Frame),
			}
		default:
			return canopen.UnexpectedSCSResponse{
				Expected:  1,
				Actual:    scs,
				AbortCode: getAbortCodeBytes(resp.Frame),
			}
		}

		// check toggle bit
		if hasBit(frame.Data[0], 4) != hasBit(resp.Frame.Data[0], 4) {
			return canopen.UnexpectedToggleBit{
				Expected:  hasBit(frame.Data[0], 4),
				Actual:    hasBit(resp.Frame.Data[0], 4),
				AbortCode: getAbortCodeBytes(resp.Frame),
			}
		}
	}

	return nil
}

func (download Download) segmentFrames() (frames []canopen.Frame) {
	if len(download.Data) <= 4 {
		return
	}

	junks := splitN(download.Data, 7)
	for i, junk := range junks {
		fdata := make([]byte, 1)

		if len(junk) < 7 {
			fdata[0] |= uint8(7-len(junk)) << 1
		}

		if i%2 == 1 {
			// toggle bit 5
			fdata[0] = setBit(fdata[0], 4)
		}

		if i == len(junks)-1 {
			// c = 1 (no more segments to download)
			fdata[0] = setBit(fdata[0], 0)
		}

		fdata = append(fdata, junk...)

		// CiA301 Standard expects all (8) bytes to be sent
		for len(fdata) < 8 {
			fdata = append(fdata, 0x0)
		}

		frame := canopen.Frame{
			CobID: download.RequestCobID,
			Data:  fdata,
		}

		frames = append(frames, frame)
	}

	return
}
