package UTO311_L0x

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

type testType struct {
	bytes []byte
}

func (t testType) MarshalUT0311L0x() ([]byte, error) {
	return []byte{0x20, 0x18, 0x08, 0x16}, nil
}

func (t *testType) UnmarshalUT0311L0x(bytes []byte) (interface{}, error) {
	b := make([]byte, 4)

	for i := 0; i < 4; i++ {
		b[i] = bytes[i] + byte(i)
	}

	v := testType{b}

	return &v, nil
}

func TestMarshalInterface(t *testing.T) {
	expected := []byte{
		0x17, 0x5f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x20, 0x18, 0x08, 0x16, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request := struct {
		MsgType   types.MsgType `uhppote:"value:0x5f"`
		Interface testType      `uhppote:"offset:33"`
	}{
		Interface: testType{},
	}

	m, err := Marshal(request)

	if err != nil {
		t.Errorf("Marshal returned unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Marshal returned invalid message - \nExpected:\n%s\nReturned:\n%s", print(expected), print(m))
		return
	}
}

func TestUnmarshalInterface(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x20, 0x18, 0x08, 0x16, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := struct {
		MsgType   types.MsgType `uhppote:"value:0x94"`
		Interface testType      `uhppote:"offset:33"`
	}{}

	err := Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if reply.MsgType != 0x94 {
		t.Errorf("Expected 'byte':0x%02X, got: 0x%02X\n", 0x94, reply.MsgType)
	}

	if !reflect.DeepEqual(reply.Interface.bytes, []byte{0x20, 0x19, 0x0a, 0x19}) {
		t.Errorf("Expected interface value '%v', got: '%v'\n", []byte{0x20, 0x19, 0x0a, 0x19}, reply.Interface)
	}
}

func TestMarshal(t *testing.T) {
	expected := []byte{
		0x17, 0x5f, 0x7d, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xd2, 0x04, 0x01, 0x00, 0xc0, 0xa8, 0x01, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d,
		0x2d, 0x55, 0x39, 0x19, 0x08, 0x92, 0x20, 0x18, 0x08, 0x16, 0x20, 0x19, 0x09, 0x17, 0x00, 0x00,
		0x00, 0x00, 0x20, 0x19, 0x04, 0x16, 0x12, 0x34, 0x56, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	mac, _ := net.ParseMAC("00:66:19:39:55:2d")
	d20180816, _ := time.ParseInLocation("2006-01-02", "2018-08-16", time.Local)
	d20190917, _ := time.ParseInLocation("2006-01-02", "2019-09-17", time.Local)
	datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-04-16 12:34:56", time.Local)

	dx20190917 := types.Date(d20190917)

	request := struct {
		MsgType      types.MsgType      `uhppote:"value:0x5f"`
		Byte         byte               `uhppote:"offset:2"`
		Uint32       uint32             `uhppote:"offset:4"`
		Uint16       uint16             `uhppote:"offset:8"`
		True         bool               `uhppote:"offset:10"`
		False        bool               `uhppote:"offset:11"`
		Address      net.IP             `uhppote:"offset:12"`
		MacAddress   types.MacAddress   `uhppote:"offset:20"`
		MAC          net.HardwareAddr   `uhppote:"offset:26"`
		SerialNumber types.SerialNumber `uhppote:"offset:32"`
		Version      types.Version      `uhppote:"offset:36"`
		Date         types.Date         `uhppote:"offset:38"`
		DatePtr      *types.Date        `uhppote:"offset:42"`
		NilDatePtr   *types.Date        `uhppote:"offset:46"`
		DateTime     types.DateTime     `uhppote:"offset:50"`
	}{
		Byte:         0x7d,
		Uint32:       423187757,
		Uint16:       1234,
		True:         true,
		False:        false,
		Address:      net.IPv4(192, 168, 1, 2),
		MacAddress:   types.MacAddress(mac),
		MAC:          mac,
		SerialNumber: 423187757,
		Version:      0x0892,
		Date:         types.Date(d20180816),
		DatePtr:      &dx20190917,
		NilDatePtr:   nil,
		DateTime:     types.DateTime(datetime),
	}

	m, err := Marshal(request)

	if err != nil {
		t.Errorf("Marshal returned unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Marshal returned invalid message - \nExpected:\n%s\nReturned:\n%s", print(expected), print(m))
		return
	}
}

func TestMarshalWithoutMsgType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Marshal(...) to fail with a panic('Missing MsgType')")
		}
	}()

	request := struct {
		MsgType byte   `uhppote:"offset:1"`
		Uint32  uint32 `uhppote:"offset:4"`
	}{
		MsgType: 0x5f,
		Uint32:  423187757,
	}

	Marshal(request)
}

func TestMarshalWithDecimalMsgType(t *testing.T) {
	expected := []byte{
		0x17, 0x5f, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request := struct {
		MsgType types.MsgType `uhppote:"value:95"`
		Uint32  uint32        `uhppote:"offset:4"`
	}{
		Uint32: 423187757,
	}

	m, err := Marshal(request)

	if err != nil {
		t.Errorf("Marshal returned unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Marshal returned invalid message - \nExpected:\n%s\nReturned:\n%s", print(expected), print(m))
		return
	}
}

func TestMarshalWithHexadecimalMsgType(t *testing.T) {
	expected := []byte{
		0x17, 0x5f, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request := struct {
		MsgType types.MsgType `uhppote:"value:0x5f"`
		Uint32  uint32        `uhppote:"offset:4"`
	}{
		Uint32: 423187757,
	}

	m, err := Marshal(request)

	if err != nil {
		t.Errorf("Marshal returned unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Marshal returned invalid message - \nExpected:\n%s\nReturned:\n%s", print(expected), print(m))
		return
	}
}

func TestUnmarshal(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x6e, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xd2, 0x04, 0x00, 0x00, 0xc0, 0xa8, 0x00, 0x00,
		0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x00, 0x66,
		0x19, 0x39, 0x55, 0x2d, 0x2d, 0x55, 0x39, 0x19, 0x08, 0x92, 0x20, 0x18, 0x08, 0x16, 0x20, 0x19,
		0x09, 0x17, 0x00, 0x00, 0x00, 0x00, 0x20, 0x18, 0x12, 0x31, 0x12, 0x23, 0x34, 0x01, 0x00, 0x00,
	}

	reply := struct {
		MsgType      types.MsgType      `uhppote:"value:0x94"`
		Byte         byte               `uhppote:"offset:2"`
		Uint32       uint32             `uhppote:"offset:4"`
		Uint16       uint16             `uhppote:"offset:8"`
		Address      net.IP             `uhppote:"offset:12"`
		SubnetMask   net.IP             `uhppote:"offset:16"`
		Gateway      net.IP             `uhppote:"offset:20"`
		MacAddress   types.MacAddress   `uhppote:"offset:24"`
		MAC          net.HardwareAddr   `uhppote:"offset:30"`
		SerialNumber types.SerialNumber `uhppote:"offset:36"`
		Version      types.Version      `uhppote:"offset:40"`
		Date         types.Date         `uhppote:"offset:42"`
		DatePtr      *types.Date        `uhppote:"offset:46"`
		NilDatePtr   *types.Date        `uhppote:"offset:50"`
		DateTime     types.DateTime     `uhppote:"offset:54"`
		True         bool               `uhppote:"offset:61"`
		False        bool               `uhppote:"offset:62"`
	}{}

	err := Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if reply.MsgType != 0x94 {
		t.Errorf("Expected 'byte':0x%02X, got: 0x%02X\n", 0x94, reply.MsgType)
	}

	if reply.Byte != 0x6e {
		t.Errorf("Expected 'byte':%02x, got: %02x\n", 0x6e, reply.Byte)
	}

	if reply.Uint32 != 423187757 {
		t.Errorf("Expected 'uint32':%v, got: %v\n", 423187757, reply.Uint32)
	}

	if reply.Uint16 != 1234 {
		t.Errorf("Expected 'uint16':%v, got: %v\n", 1234, reply.Uint16)
	}

	if !reflect.DeepEqual(reply.Address, net.IPv4(192, 168, 0, 0)) {
		t.Errorf("Expected IP address '%v', got: '%v'\n", net.IPv4(192, 168, 0, 0), reply.Address)
	}

	if !reflect.DeepEqual(reply.SubnetMask, net.IPv4(255, 255, 255, 0)) {
		t.Errorf("Expected subnet mask '%v', got: '%v'\n", net.IPv4(255, 255, 255, 0), reply.SubnetMask)
	}

	if !reflect.DeepEqual(reply.Gateway, net.IPv4(0, 0, 0, 0)) {
		t.Errorf("Expected subnet mask '%v', got: '%v'\n", net.IPv4(0, 0, 0, 0), reply.Gateway)
	}

	MAC, _ := net.ParseMAC("00:66:19:39:55:2d")
	if !reflect.DeepEqual(reply.MacAddress, types.MacAddress(MAC)) {
		t.Errorf("Expected MAC address '%v', got: '%v'\n", MAC, reply.MacAddress)
	}

	if !reflect.DeepEqual(reply.MAC, MAC) {
		t.Errorf("Expected native MAC address '%v', got: '%v'\n", MAC, reply.MAC)
	}

	if reply.Version != 0x0892 {
		t.Errorf("Expected version '0x%04X', got: '0x%04X'\n", 0x0892, reply.Version)
	}

	d20180816, _ := time.ParseInLocation("2006-01-02", "2018-08-16", time.Local)
	if reply.Date != types.Date(d20180816) {
		t.Errorf("Expected date '%v', got: '%v'\n", d20180816, reply.Date)
	}

	d20190917, _ := time.ParseInLocation("2006-01-02", "2019-09-17", time.Local)
	if reply.DatePtr == nil || *reply.DatePtr != types.Date(d20190917) {
		t.Errorf("Expected date '%v', got: '%v'\n", d20190917, reply.DatePtr)
	}

	if reply.NilDatePtr != nil {
		t.Errorf("Expected nil date '%v', got: '%v'\n", nil, reply.NilDatePtr)
	}

	datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2018-12-31 12:23:34", time.Local)
	if reply.DateTime != types.DateTime(datetime) {
		t.Errorf("Expected date '%v', got: '%v'\n", datetime, reply.DateTime)
	}

	if reply.True != true {
		t.Errorf("Expected door 1 '%v', got: '%v\n", true, reply.True)
	}

	if reply.False != false {
		t.Errorf("Expected door 2 '%v', got: '%v\n", false, reply.False)
	}
}

func TestUnmarshalWithDecimalMsgType(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := struct {
		MsgType types.MsgType `uhppote:"value:148"`
		Uint32  uint32        `uhppote:"offset:4"`
	}{}

	err := Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if reply.MsgType != 0x94 {
		t.Errorf("Expected 'byte':0x%02X, got: 0x%02X\n", 0x94, reply.MsgType)
	}

	if reply.Uint32 != 423187757 {
		t.Errorf("Expected 'uint32':%v, got: %v\n", 423187757, reply.Uint32)
	}
}

func TestUnmarshalWithHexadecimalMsgType(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := struct {
		MsgType types.MsgType `uhppote:"value:0x94"`
		Uint32  uint32        `uhppote:"offset:4"`
	}{}

	err := Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if reply.MsgType != 0x94 {
		t.Errorf("Expected 'byte':0x%02X, got: 0x%02X\n", 0x94, reply.MsgType)
	}

	if reply.Uint32 != 423187757 {
		t.Errorf("Expected 'uint32':%v, got: %v\n", 423187757, reply.Uint32)
	}
}

func TestUnmarshalWithoutMsgType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Unmarshal(...) to fail with a panic('Missing MsgType')")
		}
	}()

	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x2d, 0x55, 0x39, 0x19, 0x08, 0x92,
		0x20, 0x18, 0x08, 0x16, 0x20, 0x18, 0x12, 0x31, 0x12, 0x23, 0x34, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := struct {
		MsgType byte   `uhppote:"offset:1"`
		Uint32  uint32 `uhppote:"offset:4"`
	}{}

	Unmarshal(message, &reply)
}

func TestUnmarshalWithInvalidMsgType(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x2d, 0x55, 0x39, 0x19, 0x08, 0x92,
		0x20, 0x18, 0x08, 0x16, 0x20, 0x18, 0x12, 0x31, 0x12, 0x23, 0x34, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := struct {
		MsgType types.MsgType `uhppote:"offset:1, value:0x92"`
		Uint32  uint32        `uhppote:"offset:4"`
	}{}

	err := Unmarshal(message, &reply)

	if err == nil {
		t.Errorf("Expected error: '%v'", " Invalid value in message - expected 92, received 0x94")
		return
	}
}

func print(m []byte) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), "$1"))
}
