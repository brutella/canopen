package canopen

import (
	"reflect"
	"testing"
)

func TestHeartbeatFrame(t *testing.T) {
	frm := NewHeartbeatFrame(0x1, Operational)
	b, err := Marshal(frm)

	if err != nil {
		t.Fatal(err)
	}

	if is, want := b, []byte{0x01, 0x07, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}; reflect.DeepEqual(is, want) == false {
		t.Fatalf("\n% X !=\n% X\n", is, want)
	}
}

func TestHeartbeat(t *testing.T) {
	frm := NewHeartbeatFrame(0x1, Operational)

	if is, want := frm.FunctionCode(), Heartbeat; is != want {
		t.Fatalf("is=%X, want=%X", is, want)
	}
}
