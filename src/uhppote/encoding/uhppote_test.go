package uhppote

import (
	"encoding/hex"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"testing"
	"time"
	"uhppote/types"
)

var findDevicesRequest = []byte{
	0x17, 0x94, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

func TestMarshalFindDevicesRequest(t *testing.T) {
	request := struct {
		MsgType byte `uhppote:"offset:1"`
	}{
		0x94,
	}

	m, err := Marshal(request)

	if err != nil {
		t.Errorf("Marshal(%s) returned unexpected error: %v", "FindDevicesRequest", err)
		return
	}

	if !reflect.DeepEqual(m, findDevicesRequest) {
		t.Errorf("Invalid byte array for uhppote.Marshal(%s):\nExpected:\n%s\nReturned:\n%s", "FindDevicesRequest", print(findDevicesRequest), print(m))
		return
	}
}

func TestUnmarshal(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x08, 0x92, 0x20, 0x18, 0x08, 0x16,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := struct {
		Code         byte             `uhppote:"offset:1"`
		SerialNumber uint32           `uhppote:"offset:4"`
		IpAddress    net.IP           `uhppote:"offset:8"`
		SubnetMask   net.IP           `uhppote:"offset:12"`
		Gateway      net.IP           `uhppote:"offset:16"`
		MacAddress   net.HardwareAddr `uhppote:"offset:20"`
		Version      types.Version    `uhppote:"offset:26"`
		Date         types.Date       `uhppote:"offset:28"`
	}{}

	err := Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if reply.Code != 0x94 {
		t.Errorf("Expected command code 0x%02X, got: 0x%02X\n", 0x94, reply.Code)
	}

	if reply.SerialNumber != 423187757 {
		t.Errorf("Expected serial number %v, got: %v\n", 423187757, reply.SerialNumber)
	}

	if !reflect.DeepEqual(reply.IpAddress, net.IPv4(192, 168, 0, 0)) {
		t.Errorf("Expected IP address '%v', got: '%v'\n", net.IPv4(192, 168, 0, 0), reply.IpAddress)
	}

	if !reflect.DeepEqual(reply.SubnetMask, net.IPv4(255, 255, 255, 0)) {
		t.Errorf("Expected subnet mask '%v', got: '%v'\n", net.IPv4(255, 255, 255, 0), reply.SubnetMask)
	}

	if !reflect.DeepEqual(reply.Gateway, net.IPv4(0, 0, 0, 0)) {
		t.Errorf("Expected subnet mask '%v', got: '%v'\n", net.IPv4(0, 0, 0, 0), reply.Gateway)
	}

	MAC, _ := net.ParseMAC("00:66:19:39:55:2d")
	if !reflect.DeepEqual(reply.MacAddress, MAC) {
		t.Errorf("Expected MAC address '%v', got: '%v'\n", MAC, reply.MacAddress)
	}

	if reply.Version != 0x0892 {
		t.Errorf("Expected version '0x%04X', got: '0x%04X'\n", 0x0892, reply.Version)
	}

	date, _ := time.ParseInLocation("20060102", "20180816", time.Local)
	if reply.Date.Date != date {
		t.Errorf("Expected date '%v', got: '%v'\n", date, reply.Date)
	}
}

func print(m []byte) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), "$1"))
}