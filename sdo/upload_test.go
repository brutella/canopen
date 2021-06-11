package sdo

import (
	"bytes"
	"github.com/FabianPetersen/can"
	"github.com/FabianPetersen/canopen"
	"reflect"
	"testing"
)

type uploadReadWriteCloser struct {
	buf bytes.Buffer
}

func (rw *uploadReadWriteCloser) Read(b []byte) (n int, err error) {
	return rw.buf.Read(b)
}

func (rw *uploadReadWriteCloser) Write(b []byte) (n int, err error) {
	if b, err := can.Marshal(uploadFrame); err == nil {
		return rw.buf.Write(b)
	}

	return len(b), nil
}

func (rw *uploadReadWriteCloser) Close() error { return nil }

var uploadFrame = can.Frame{
	ID:     0x580,
	Length: 0x8,
	Flags:  0x0,
	Res0:   0x0,
	Res1:   0x0,
	Data: [can.MaxFrameDataLength]uint8{
		0x42, // 0100 0010 (= expedited)
		0xBB, 0xAA,
		0xCC,
		0x1, 0x2, 0x3, 0x4},
}

func HandleUpload(b []byte) (n int, err error) {
	var buf bytes.Buffer
	if b, err = can.Marshal(uploadFrame); err == nil {
		n, err = buf.Write(b)
	}

	return
}

func TestExpeditedUpload(t *testing.T) {
	rwc := can.NewReadWriteCloser(&uploadReadWriteCloser{})
	bus := can.NewBus(rwc)

	go bus.ConnectAndPublish()
	defer bus.Disconnect()

	object := canopen.NewObjectIndex(0xAABB, 0xCC)
	// Read values for object
	upload := Upload{
		ObjectIndex:   object,
		RequestCobID:  0x600,
		ResponseCobID: uint16(uploadFrame.ID),
	}

	b, err := upload.Do(bus)
	if err != nil {
		t.Fatal(err)
	}

	if is, want := b, []byte{0x1, 0x2, 0x3, 0x4}; reflect.DeepEqual(is, want) != true {
		t.Fatalf("is=%v want=%v", is, want)
	}
}
