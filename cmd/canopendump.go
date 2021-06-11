// This program logs CANopen frames to the console.
package main

import (
	"flag"
	"fmt"
	"github.com/FabianPetersen/can"
	"github.com/FabianPetersen/canopen"
	"log"
	"os"
	"os/signal"
)

var i = flag.String("if", "", "network interface name")

func main() {
	flag.Parse()
	if len(*i) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	bus, err := can.NewBusForInterfaceWithName(*i)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("+------+--------------+-------------------------+")
	log.Println("| Node | Message Type | Bytes                   |")
	log.Println("+------+--------------+-------------------------+")
	bus.SubscribeFunc(logCANFrame)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	go func() {
		select {
		case <-c:
			bus.Disconnect()
			os.Exit(1)
		}
	}()

	bus.ConnectAndPublish()
}

// logCANFrame logs a frame with the same format as candump from can-utils.
func logCANFrame(frm can.Frame) {
	canopenFrm := canopen.CANopenFrame(frm)
	var msgType string
	var data string
	switch canopenFrm.MessageType() {
	case canopen.MessageTypeNMT:
		msgType = "Nmt"
	case canopen.MessageTypeHeartbeat:
		msgType = "Heartbeat"
	case canopen.MessageTypeSync:
		msgType = "Sync"
	case canopen.MessageTypeTimestamp:
		msgType = "Timestamp"
		if time, _ := canopenFrm.Timestamp(); time != nil {
			data = fmt.Sprintf("Timestamp %s", time.String())
		}
	case canopen.MessageTypeTSDO:
		msgType = "SDO Response"
	case canopen.MessageTypeRSDO:
		msgType = "SDO Request"
	default:
		msgType = "Unkown"
	}

	log.Printf("| %-4d | %-12s | % -23X | %s", canopenFrm.NodeID(), msgType, frm.Data, data)
}
