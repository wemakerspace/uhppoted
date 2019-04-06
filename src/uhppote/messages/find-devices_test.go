package messages

import (
	"net"
	"reflect"
	"testing"
	"time"
	"uhppote/encoding"
	"uhppote/types"
)

func TestParseFindResponse(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x08, 0x92, 0x20, 0x18, 0x08, 0x16,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := struct {
		MsgType      byte             `uhppote:"offset:1"`
		SerialNumber uint32           `uhppote:"offset:4"`
		IpAddress    net.IP           `uhppote:"offset:8"`
		SubnetMask   net.IP           `uhppote:"offset:12"`
		Gateway      net.IP           `uhppote:"offset:16"`
		MacAddress   net.HardwareAddr `uhppote:"offset:20"`
		Version      types.Version    `uhppote:"offset:26"`
		Date         types.Date       `uhppote:"offset:28"`
	}{}

	err := uhppote.Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Fins returned error from valid message: %v\n", err)
	}

	if reply.MsgType != 0x94 {
		t.Errorf("Fins returned incorrect 'message type' from valid message: %02x\n", reply.MsgType)
	}

	if reply.SerialNumber != 423187757 {
		t.Errorf("Fins returned incorrect 'serial number' from valid message: %v\n", reply.SerialNumber)
	}

	if !reflect.DeepEqual(reply.IpAddress, net.IPv4(192, 168, 0, 0)) {
		t.Errorf("Fins returned incorrect 'IP address' from valid message: %v\n", reply.IpAddress)
	}

	if !reflect.DeepEqual(reply.SubnetMask, net.IPv4(255, 255, 255, 0)) {
		t.Errorf("Fins returned incorrect 'subnet mask' from valid message: %v\n", reply.SubnetMask)
	}

	if !reflect.DeepEqual(reply.Gateway, net.IPv4(0, 0, 0, 0)) {
		t.Errorf("Fins returned incorrect 'gateway' from valid message: %v\n", reply.Gateway)
	}

	MAC, _ := net.ParseMAC("00:66:19:39:55:2d")
	if !reflect.DeepEqual(reply.MacAddress, MAC) {
		t.Errorf("Fins returned incorrect 'MAC address' from valid message: %v\n", reply.MacAddress)
	}

	if reply.Version != 0x0892 {
		t.Errorf("Fins returned incorrect 'version' from valid message: %v\n", reply.Version)
	}

	date, _ := time.ParseInLocation("20060102", "20180816", time.Local)
	if reply.Date.Date != date {
		t.Errorf("Fins returned incorrect 'date' from valid message: %v\n", reply.Date)
	}
}