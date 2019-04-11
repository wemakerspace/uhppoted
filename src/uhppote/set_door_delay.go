package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

type SetDoorDelayRequest struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
	Door         uint8  `uhppote:"offset:8"`
	Unit         uint8  `uhppote:"offset:9"`
	Delay        uint8  `uhppote:"offset:10"`
}

type SetDoorDelayResponse struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
	Door         uint8  `uhppote:"offset:8"`
	Unit         uint8  `uhppote:"offset:9"`
	Delay        uint8  `uhppote:"offset:10"`
}

func (u *UHPPOTE) SetDoorDelay(serialNumber uint32, door uint8, delay uint8) (*types.DoorDelay, error) {
	request := SetDoorDelayRequest{
		MsgType:      0x80,
		SerialNumber: serialNumber,
		Door:         door,
		Unit:         0x03,
		Delay:        delay,
	}

	reply := SetDoorDelayResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0x80 {
		return nil, errors.New(fmt.Sprintf("SetDoorDelay returned incorrect message type: %02X\n", reply.MsgType))
	}

	if reply.Unit != 0x03 {
		return nil, errors.New(fmt.Sprintf("SetDoorDelay returned incorrect time unit: %02X\n", reply.Unit))
	}

	return &types.DoorDelay{
		SerialNumber: reply.SerialNumber,
		Door:         reply.Door,
		Delay:        reply.Delay,
	}, nil
}
