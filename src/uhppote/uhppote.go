package uhppote

import (
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"time"
	codec "uhppote/encoding/UTO311-L0x"
)

type UHPPOTE struct {
	BindAddress net.UDPAddr
	Debug       bool
}

func (u *UHPPOTE) Exec(request, reply interface{}) error {
	p, err := codec.Marshal(request)
	if err != nil {
		return err
	}

	q, err := u.Execute(p)

	err = codec.Unmarshal(q, reply)
	if err != nil {
		return err
	}

	return nil
}

func (u *UHPPOTE) Broadcast(request interface{}) ([][]byte, error) {
	p, err := codec.Marshal(request)
	if err != nil {
		return [][]byte{}, err
	}

	return u.broadcast(p)
}

func (u *UHPPOTE) Execute(cmd []byte) ([]byte, error) {
	reply := make([]byte, 2048)

	if u.Debug {
		regex := regexp.MustCompile("(?m)^(.*)")

		fmt.Printf(" ... command %v bytes\n", len(cmd))
		fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(cmd), " ...         $1"))
	}

	broadcast, err := net.ResolveUDPAddr("udp", "255.255.255.255:60000")

	if err != nil {
		return nil, makeErr("Failed to resolve UDP broadcast address", err)
	}

	connection, err := net.ListenUDP("udp", &u.BindAddress)

	if err != nil {
		return nil, makeErr("Failed to open UDP socket", err)
	}

	defer close(connection)

	N, err := connection.WriteTo(cmd, broadcast)

	if err != nil {
		return nil, makeErr("Failed to write to UDP socket", err)
	}

	if u.Debug {
		fmt.Printf(" ... sent %v bytes\n", N)
	}

	err = connection.SetDeadline(time.Now().Add(5000 * time.Millisecond))

	if err != nil {
		return nil, makeErr("Failed to set UDP timeout", err)
	}

	N, remote, err := connection.ReadFromUDP(reply)

	if err != nil {
		return nil, makeErr("Failed to read from UDP socket", err)
	}

	if u.Debug {
		regex := regexp.MustCompile("(?m)^(.*)")

		fmt.Printf(" ... received %v bytes from %v\n", N, remote)
		fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(reply[:N]), " ...          $1"))
	}

	return reply[:N], nil
}

func (u *UHPPOTE) broadcast(cmd []byte) ([][]byte, error) {
	replies := make([][]byte, 0)

	if u.Debug {
		regex := regexp.MustCompile("(?m)^(.*)")

		fmt.Printf(" ... command %v bytes\n", len(cmd))
		fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(cmd), " ...         $1"))
	}

	broadcast, err := net.ResolveUDPAddr("udp", "255.255.255.255:60000")

	if err != nil {
		return nil, makeErr("Failed to resolve UDP broadcast address", err)
	}

	connection, err := net.ListenUDP("udp", &u.BindAddress)

	if err != nil {
		return nil, makeErr("Failed to open UDP socket", err)
	}

	defer close(connection)

	N, err := connection.WriteTo(cmd, broadcast)

	if err != nil {
		return nil, makeErr("Failed to write to UDP socket", err)
	}

	if u.Debug {
		fmt.Printf(" ... sent %v bytes\n", N)
	}

	go func() {
		for {
			reply := make([]byte, 2048)
			N, remote, err := connection.ReadFromUDP(reply)

			if err != nil {
				break
			} else {
				replies = append(replies, reply[:N])

				if u.Debug {
					regex := regexp.MustCompile("(?m)^(.*)")

					fmt.Printf(" ... received %v bytes from %v\n", N, remote)
					fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(reply[:N]), " ...          $1"))
				}
			}
		}
	}()

	time.Sleep(2500 * time.Millisecond)

	return replies, err
}

func close(connection net.Conn) {
	connection.Close()
}
