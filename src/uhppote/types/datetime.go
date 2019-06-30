package types

import (
	"errors"
	"fmt"
	"time"
	"uhppote/encoding/bcd"
)

type DateTime time.Time

func (d DateTime) String() string {
	return time.Time(d).Format("2006-01-02 15:04:05")
}

func (d DateTime) MarshalUT0311L0x() ([]byte, error) {
	encoded, err := bcd.Encode(time.Time(d).Format("20060102150405"))

	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("Error encoding datetime %v to BCD: [%v]", d, err))
	}

	if encoded == nil {
		return []byte{}, errors.New(fmt.Sprintf("Unknown error encoding datetime %v to BCD", d))
	}

	return *encoded, nil
}

func (d *DateTime) UnmarshalUT0311L0x(bytes []byte) (interface{}, error) {
	decoded, err := bcd.Decode(bytes[0:7])
	if err != nil {
		return nil, err
	}

	datetime, err := time.ParseInLocation("20060102150405", decoded, time.Local)
	if err != nil {
		return nil, err
	}

	v := DateTime(datetime)

	return &v, nil
}
