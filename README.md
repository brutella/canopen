# canopen

canopen lets you send and received data in a [CANopen](https://en.wikipedia.org/wiki/CANopen) network.

The library contains basic functionality and doesn't aim to be a complete implementation of the CANopen protocol. You can find a complete CANopen implementation in the [CANopenNode](https://github.com/CANopenNode) project.

## CANopen

CANopen is a protocol to communicate on a CAN bus. Data is exchanged between nodes using CANopen frames, which uses a subset of the bytes in a CAN frames. This project extends the [can](https://github.com/brutella/can) library to interact with a CANopen nodes. Setup your hard- and software as described [there](https://github.com/brutella/can/blob/master/README.md).

You can find a very good documentation about CANopen [here](http://www.a-m-c.com/download/sw/dw300_3-0-3/CAN_Manual300_3-0-3.pdf).

## Usage

##### Setup the CAN bus interface

```go
bus, _ := can.NewBusForInterfaceWithName("can0")
bus.ConnectAndPublish()
```

##### CANopen request/response communication

Parts of CANopen protocol are based on request/response communication. The library makes it easy to send a request and wait for the reponse.

```go
// Frame to be sent
payload := []byte{...}
frame := canopen.NewFrame(canopen.MessageTypeRSDO, payload)

// Expected id of response frame
respID := canopen.MessageTypeTSDO

req := NewRequest(frame, respID)

// Create client which sends request and waits for response
client := &Client{bus, time.Second * 1}
resp, _ := client.Do(req)
```

# Contact

Matthias Hochgatterer

Github: [https://github.com/brutella](https://github.com/brutella/)

Twitter: [https://twitter.com/brutella](https://twitter.com/brutella)

# License

*canopen* is available under the MIT license. See the LICENSE file for more info.