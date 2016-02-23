package canopen

import (
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	fr := NewFrame(0xFFFF, []uint8{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})

	b, err := Marshal(fr)
	if err != nil {
		t.Fatal(err)
	}

	if is, want := b, []byte{0xFF, 0x7, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}; reflect.DeepEqual(is, want) == false {
		t.Fatalf("\n% X !=\n% X\n", is, want)
	}
}

func TestUnmarshal(t *testing.T) {
	b := []byte{
		0x01, 0x7, 0x00, 0x00,
		0x01,
		0x00,
		0x00,
		0x00,
		0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	frm := Frame{}

	Unmarshal(b, &frm)

	if is, want := frm.CobID, uint16(0x701); is != want {
		t.Fatalf("is=%X, want=%X", is, want)
	}

	if is, want := frm.MessageType(), MessageTypeHeartbeat; is != want {
		t.Fatalf("is=%X, want=%X", is, want)
	}
}
