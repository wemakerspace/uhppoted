package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"uhppote"
	"uhppote-cli/commands"
)

type addr struct {
	address *net.UDPAddr
}

var cli = []commands.Command{
	&commands.VersionCommand{VERSION},
	&commands.GetDevicesCommand{},
	&commands.SetAddressCommand{},
	&commands.GetStatusCommand{},
	&commands.GetTimeCommand{},
	&commands.SetTimeCommand{},
	&commands.GetDoorDelayCommand{},
	&commands.SetDoorDelayCommand{},
	&commands.GetListenerCommand{},
	&commands.SetListenerCommand{},
	&commands.GetCardsCommand{},
	&commands.GetCardCommand{},
	&commands.GrantCommand{},
	&commands.RevokeCommand{},
	&commands.RevokeAllCommand{},
	&commands.GetEventsCommand{},
	&commands.GetEventIndexCommand{},
	&commands.SetEventIndexCommand{},
	&commands.OpenDoorCommand{},
	&commands.ListenCommand{},
}

var VERSION = "v0.00.0"
var debug = false
var local = addr{nil}
var broadcast = addr{nil}

func main() {
	flag.Var(&local, "bind", "Sets the local IP address and port to which to bind (e.g. 192.168.0.100:60001)")
	flag.Var(&broadcast, "broadcast", "Sets the IP address and port for UDP broadcast (e.g. 192.168.0.255:60000)")
	flag.BoolVar(&debug, "debug", false, "Displays vaguely useful information while processing a command")
	flag.Parse()

	u := uhppote.UHPPOTE{
		BindAddress:      local.address,
		BroadcastAddress: broadcast.address,
		Debug:            debug,
	}

	cmd, err := parse()
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}

	if cmd == nil {
		help()
		return
	}

	err = cmd.Execute(&u)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func parse() (commands.Command, error) {
	var cmd commands.Command = nil
	var err error = nil

	if len(os.Args) > 1 {
		switch flag.Arg(0) {
		default:
			for _, c := range cli {
				if c.CLI() == flag.Arg(0) {
					cmd = c
				}
			}
		}
	}

	return cmd, err
}

func (b *addr) String() string {
	return b.address.String()
}

func (b *addr) Set(s string) error {
	address, err := net.ResolveUDPAddr("udp", s)
	if err != nil {
		return err
	}

	b.address = address

	return nil
}

func help() {
	if len(flag.Args()) > 0 && flag.Arg(0) == "help" {
		if len(flag.Args()) > 1 {

			if flag.Arg(1) == "commands" {
				helpCommands()
				return
			}

			for _, c := range cli {
				if c.CLI() == flag.Arg(1) {
					c.Help()
					return
				}
			}

			fmt.Printf("Invalid command: %v. Type 'help commands' to get a list of supported commands\n", flag.Arg(1))
			return
		}
	}

	usage()
}

func usage() {
	fmt.Println()
	fmt.Println("  Usage: uhppote-cli [options] <command>")
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println()
	fmt.Println("    help             Displays this message")
	fmt.Println("                     For help on a specific command use 'uhppote-cli help <command>'")

	for _, c := range cli {
		fmt.Printf("    %-16s %s\n", c.CLI(), c.Description())
	}

	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --bind      Sets the local IP address and port to use")
	fmt.Println("    --broadcast Sets the IP address and port to use for UDP broadcast")
	fmt.Println("    --debug     Displays vaguely useful internal information")
	fmt.Println()
}

func helpCommands() {
	fmt.Println("Supported commands:")
	fmt.Println()

	for _, c := range cli {
		fmt.Printf(" %-16s %s\n", c.CLI(), c.Usage())
	}

	fmt.Println()
}
