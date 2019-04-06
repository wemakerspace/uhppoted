package commands

import (
	"fmt"
	"uhppote"
)

type ListDevicesCommand struct {
}

func NewListDevicesCommand() (*ListDevicesCommand, error) {
	return &ListDevicesCommand{}, nil
}

func (c *ListDevicesCommand) Execute(u *uhppote.UHPPOTE) error {
	devices, err := u.FindDevices()

	if err == nil {
		for _, device := range devices {
			fmt.Printf("%s\n", device.String())
		}
	}

	return err
}

func (c *ListDevicesCommand) CLI() string {
	return "get-devices"
}

func (c *ListDevicesCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] list-devices [command options]")
	fmt.Println()
	fmt.Println(" Searches the local network for UHPPOTE access control boards reponding to a poll")
	fmt.Println(" on the default UDP port 60000. Returns a list of boards one per line in the format:")
	fmt.Println()
	fmt.Println(" <serial number> <IP address> <subnet mask> <gateway> <MAC address> <hexadecimal version> <firmware date>")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays vaguely useful internal information")
	fmt.Println()
}