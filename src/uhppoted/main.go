package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uhppote"
	"uhppoted/commands"
	"uhppoted/config"
	"uhppoted/eventlog"
	"uhppoted/rest"
)

const (
	LOGFILESIZE = 1
	IDLE        = time.Duration(60 * time.Second)
)

var VERSION = "v0.04.0"
var retries = 0

func main() {
	flag.Parse()

	cmd, err := commands.Parse()
	if err != nil {
		fmt.Printf("\nError parsing command line: %v\n\n", err)
		os.Exit(1)
	}

	if cmd != nil {
		ctx := commands.Context{}
		if err = cmd.Execute(ctx); err != nil {
			fmt.Printf("\nERROR: %v\n\n", err)
			os.Exit(1)
		}

		return
	}

	// ... default to 'run'

	sysinit()

	config, err := config.LoadConfig(*configuration)
	if err != nil {
		fmt.Printf("\n   WARN: Error loading configuration: %v\n", err)
	}

	if err := os.MkdirAll(*dir, os.ModeDir|os.ModePerm); err != nil {
		log.Fatal(fmt.Sprintf("Error creating working directory '%v'", *dir), err)
	}

	pid := fmt.Sprintf("%d\n", os.Getpid())

	if err := ioutil.WriteFile(*pidFile, []byte(pid), 0644); err != nil {
		log.Fatal("Error creating pid file: %v\n", err)
	}

	defer cleanup(*pidFile)

	// ... use syslog for console logging?

	if *useSyslog {
		logger, err := syslog.New(syslog.LOG_NOTICE, "uhppoted")

		if err != nil {
			log.Fatal("Error opening syslog: ", err)
			return
		}

		log.SetOutput(logger)
	}

	run(&config, *logfile, *logfilesize)
}

func cleanup(pid string) {
	os.Remove(pid)
}

func run(c *config.Config, logfile string, logfilesize int) {
	// ... setup logging

	events := eventlog.Ticker{Filename: logfile, MaxSize: logfilesize}
	logger := log.New(&events, "", log.Ldate|log.Ltime|log.LUTC)

	// ... syscall SIG handlers

	interrupt := make(chan os.Signal, 1)
	rotate := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(rotate, syscall.SIGHUP)

	go func() {
		for {
			<-rotate
			log.Printf("Rotating uhppoted log file '%s'\n", logfile)
			events.Rotate()
		}
	}()

	// ... listen forever

	for {
		err := listen(c, logger, interrupt)

		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		log.Printf("exit\n")
		break
	}
}

func listen(c *config.Config, logger *log.Logger, interrupt chan os.Signal) error {
	// ... listen

	log.Printf("... listening")

	u := uhppote.UHPPOTE{
		BindAddress:      &c.BindAddress,
		BroadcastAddress: &c.BroadcastAddress,
		Devices:          make(map[uint32]*net.UDPAddr),
		Debug:            true,
	}

	for id, d := range c.Devices {
		if d.Address != nil {
			u.Devices[id] = d.Address
		}
	}

	go func() {
		rest.Run(&u, logger)
	}()

	defer rest.Close()

	touched := time.Now()
	closed := make(chan struct{})

	// ... wait until interrupted/closed

	k := time.NewTicker(15 * time.Second)
	tick := time.NewTicker(5 * time.Second)

	defer k.Stop()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			if err := watchdog(touched); err != nil {
				return err
			}

		case <-k.C:
			log.Printf("... keep-alive")
			keepalive()

		case <-interrupt:
			log.Printf("... interrupt")
			return nil

		case <-closed:
			log.Printf("... closed")
			return errors.New("Server error")
		}
	}

	log.Printf("... exit")
	return nil
}
