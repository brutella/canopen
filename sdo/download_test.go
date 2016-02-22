package sdo

import (
	"bytes"
	"github.com/brutella/can"
	"github.com/brutella/canopen"
	"log"
	"testing"
)

type downloadReadWriteCloser struct {
	buf bytes.Buffer
}

func (rw *downloadReadWriteCloser) Read(b []byte) (n int, err error) {
	return rw.buf.Read(b)
}

func (rw *downloadReadWriteCloser) Write(b []byte) (n int, err error) {
	var frm can.Frame
	if err := can.Unmarshal(b, &frm); err != nil {
		return 0, err
	}

	switch frm.Data[0] & CommandSpecifierMask {
	case ClientIntiateDownload:
		frm = downloadInitiateFrame
	case ClientSegmentDownload:
		frm = downloadSegmentFrame
	default:
		log.Fatal("Unknown command")
		break
	}
	if b, err := can.Marshal(frm); err == nil {
		return rw.buf.Write(b)
	} else {
		return 0, err
	}

	return len(b), nil
}

func (rw *downloadReadWriteCloser) Close() error { return nil }

var downloadInitiateFrame = can.Frame{
	ID:     0x580,
	Length: 0x8,
	Flags:  0x0,
	Res0:   0x0,
	Res1:   0x0,
	Data: [can.MaxFrameDataLength]uint8{
		ServerInitiateDownload,
		0xBB, 0xAA,
		0xCC,
		0x0, 0x0, 0x0, 0x0},
}

var downloadSegmentFrame = can.Frame{
	ID:     0x580,
	Length: 0x8,
	Flags:  0x0,
	Res0:   0x0,
	Res1:   0x0,
	Data: [can.MaxFrameDataLength]uint8{
		ServerSegmentDownload,
		0xBB, 0xAA,
		0xCC,
		0x0, 0x0, 0x0, 0x0},
}

func TestDownload(t *testing.T) {
	rwc := can.NewReadWriteCloser(&downloadReadWriteCloser{})
	bus := can.NewBus(rwc)

	go bus.ConnectAndPublish()
	defer bus.Disconnect()

	object := canopen.NewObjectIndex(0xAABB, 0xCC)
	download := Download{
		ObjectIndex:   object,
		RequestCobID:  0x600,
		ResponseCobID: uint16(downloadSegmentFrame.ID),
		// 0x2 + WRITE (String) + 0x91 (Datatype)
		Data: []byte{0x2, 0x57, 0x52, 0x49, 0x54, 0x45, 0x91},
	}

	err := download.Do(bus)
	if err != nil {
		t.Fatal(err)
	}
}
