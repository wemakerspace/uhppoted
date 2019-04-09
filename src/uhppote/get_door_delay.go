package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

type GetDoorDelayRequest struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
	Door         uint8  `uhppote:"offset:8"`
}

type GetDoorDelayResponse struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
	Door         uint8  `uhppote:"offset:8"`
	Unit         uint8  `uhppote:"offset:9"`
	Delay        uint8  `uhppote:"offset:10"`
}

func (u *UHPPOTE) GetDoorDelay(serialNumber uint32, door byte) (*types.DoorDelay, error) {
	request := GetDoorDelayRequest{
		MsgType:      0x82,
		SerialNumber: serialNumber,
		Door:         door,
	}

	reply := GetDoorDelayResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0x82 {
		return nil, errors.New(fmt.Sprintf("GetDoorDelay returned incorrect message type: %02X\n", reply.MsgType))
	}

	if reply.Unit != 0x03 {
		return nil, errors.New(fmt.Sprintf("GetDoorDelay returned incorrect time unit: %02X\n", reply.Unit))
	}

	return &types.DoorDelay{
		SerialNumber: reply.SerialNumber,
		Door:         reply.Door,
		Delay:        reply.Delay,
	}, nil
}
