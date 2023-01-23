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

	if err, _ := download.doInitFrame(bus, false); err != nil {
		return err
	}

	return download.doSegments(bus)
}

func (download Download) DoBlock(bus *can.Bus) error {
	// Do not allow multiple messages for the same device
	key := strconv.Itoa(int(download.RequestCobID))
	canopen.Lock.Lock(key)
	defer canopen.Lock.Unlock(key)

	var err error
	var segmentsPerBlock int
	if err, segmentsPerBlock = download.doInitFrame(bus, true); err != nil {
		return err
	}

	return download.doBlock(bus, segmentsPerBlock)
}

func (download Download) doInitFrame(bus *can.Bus, isBlockTransfer bool) (error, int) {
	segmentsPerBlock := 0
	frame, err := download.initFrame(isBlockTransfer)
	if err != nil {
		return err, segmentsPerBlock
	}

	req := canopen.NewRequest(frame, uint32(download.ResponseCobID))
	c := &canopen.Client{Bus: bus, Timeout: time.Second * 2}
	resp, err := c.Do(req)
	if err != nil {
		return err, segmentsPerBlock
	}

	frame = resp.Frame
	if isBlockTransfer {
		segmentsPerBlock = int(frame.Data[4])
	}

	scs := frame.Data[0] >> 5
	if !isBlockTransfer && scs == 3 || isBlockTransfer && scs == 5 { // Success
		// Check if this is the correct response for the requested message
		if frame.Data[1] != download.ObjectIndex.Index.B0 || frame.Data[2] != download.ObjectIndex.Index.B1 || frame.Data[3] != download.ObjectIndex.SubIndex {
			return canopen.TransferAbort{
				AbortCode: getAbortCodeBytes(frame),
			}, segmentsPerBlock
		}

	} else if scs == 4 { // Abort
		return canopen.TransferAbort{
			AbortCode: getAbortCodeBytes(frame),
		}, segmentsPerBlock

	} else {
		return canopen.UnexpectedSCSResponse{
			Expected:  3,
			Actual:    scs,
			AbortCode: getAbortCodeBytes(frame),
		}, segmentsPerBlock
	}

	return nil, segmentsPerBlock
}

// initFrame returns the initial frame of the download.
// If the download data is less than 4 bytes, the init frame data contains all download data.
// If the download data is more than 4 bytes, the init frame data contains the overall length of the download data.
func (download Download) initFrame(isBlockTransfer bool) (frame canopen.Frame, err error) {
	fdata := make([]byte, 4)

	if isBlockTransfer {
		// css = 6 (download init block request)
		fdata[0] = setBit(fdata[0], 6)
		fdata[0] = setBit(fdata[0], 7)
	} else {
		// css = 1 (download init request)
		fdata[0] = setBit(fdata[0], 5)
	}

	fdata[1] = download.ObjectIndex.Index.B0
	fdata[2] = download.ObjectIndex.Index.B1
	fdata[3] = download.ObjectIndex.SubIndex

	n := len(download.Data)
	if n <= 4 && !isBlockTransfer { // does download data fit into one frame?
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
		if isBlockTransfer {
			// Always indicate size in block transfer
			fdata[0] = setBit(fdata[0], 1)
		} else {
			// e = 0
			// n = 0 (frame.Data contains the overall )
			// s = 1
			fdata[0] = setBit(fdata[0], 0)
		}

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

func (download Download) doBlock(bus *can.Bus, segmentsPerBlock int) error {
	frames := download.segmentFrames(true)

	c := &canopen.Client{Bus: bus, Timeout: time.Second * 2}
	segmentIndex := 0
	index := 0
	for segmentIndex < len(frames) {
		// Don't wait for the confirmation frame
		var err error = nil
		for ; err == nil && index+1 < segmentsPerBlock && (segmentIndex+index+1) < len(frames); index++ {
			frames[segmentIndex+index].Data[0] = getFirstByte(index, false, 7, true)
			err = bus.PublishMinDuration(frames[segmentIndex+index].CANFrame(), 2*time.Microsecond)
		}

		// Wait for the confirmation frame
		frames[segmentIndex+index].Data[0] = getFirstByte(index, segmentIndex+index+1 == len(frames), 7, true)
		req := canopen.NewRequest(frames[segmentIndex+index], uint32(download.ResponseCobID))
		resp, err := c.DoMinDuration(req, 1*time.Millisecond)
		if err != nil {
			return err
		}

		// Mask out the correct bits
		scs := resp.Frame.Data[0] >> 5
		ss := resp.Frame.Data[0] & 0x3

		if scs == 5 && ss == 2 {
			segmentsPerBlock = int(resp.Frame.Data[2])
			ackSegment := int(resp.Frame.Data[1])

			// Everything is fine, move along please
			if ackSegment == segmentsPerBlock {
				segmentIndex += ackSegment
				index = 0
				// Retry from the last acked segment
			} else {
				index = ackSegment
			}
		} else {
			return canopen.UnexpectedSCSResponse{
				Expected:  5,
				Actual:    scs,
				AbortCode: getAbortCodeBytes(resp.Frame),
			}
		}
	}

	// Send the end block
	err := download.doBlockEnd(c)
	if err != nil {
		return err
	}

	return nil
}

func (download Download) doBlockEnd(c *canopen.Client) error {
	fdata := make([]byte, 8)

	// n (Set the length of data in the last frame in the last segment)
	fdata[0] |= uint8(7-(len(download.Data)%7)) << 2

	// css = 6 (download init block request)
	fdata[0] = setBit(fdata[0], 6)
	fdata[0] = setBit(fdata[0], 7)

	// cs = 1 (indicate download end)
	fdata[0] = setBit(fdata[0], 0)

	req := canopen.NewRequest(canopen.NewFrame(download.RequestCobID, fdata), uint32(download.ResponseCobID))
	resp, err := c.DoMinDuration(req, 0)

	if err != nil {
		return err
	}

	scs := resp.Frame.Data[0] >> 5
	ss := hasBit(resp.Frame.Data[0], 0)

	if scs == 5 && ss {
		return nil
	} else {
		return canopen.UnexpectedSCSResponse{
			Expected:  5,
			Actual:    scs,
			AbortCode: getAbortCodeBytes(resp.Frame),
		}
	}
}

func (download Download) doSegments(bus *can.Bus) error {
	frames := download.segmentFrames(false)

	c := &canopen.Client{Bus: bus, Timeout: time.Second * 2}
	for _, frame := range frames {
		req := canopen.NewRequest(frame, uint32(download.ResponseCobID))
		resp, err := c.DoMinDuration(req, 2*time.Millisecond)
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

func (download Download) segmentFrames(isBlockTransfer bool) (frames []canopen.Frame) {
	if !isBlockTransfer && len(download.Data) <= 4 {
		return
	}

	junks := splitN(download.Data, 7)
	for i, junk := range junks {
		fdata := append([]byte{getFirstByte(i, i == len(junks)-1, len(junk), isBlockTransfer)}, junk...)

		// CiA301 Standard expects all (8) bytes to be sent
		for len(fdata) < 8 {
			fdata = append(fdata, 0x0)
		}

		frames = append(frames, canopen.Frame{
			CobID: download.RequestCobID,
			Data:  fdata,
		})
	}

	return
}

func getFirstByte(i int, isLast bool, junkLength int, isBlockTransfer bool) byte {
	firstByte := byte(0)
	if !isBlockTransfer {
		if junkLength < 7 {
			firstByte |= uint8(7-junkLength) << 1
		}

		if i%2 == 1 {
			// toggle bit 5
			firstByte = setBit(firstByte, 4)
		}
	} else {
		// Set the segment number
		firstByte = byte(1 + i)
	}

	if isLast {
		// c = 1 (no more segments to download)
		if isBlockTransfer {
			firstByte = setBit(firstByte, 7)
		} else {
			firstByte = setBit(firstByte, 0)
		}
	}

	return firstByte
}
