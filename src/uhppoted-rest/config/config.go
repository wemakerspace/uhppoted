package config

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"uhppote/encoding/conf"
)

type DeviceMap map[uint32]*Device

type Device struct {
	Address *net.UDPAddr
	Door    []string
}

type REST struct {
	HttpEnabled        bool
	HttpPort           uint16
	HttpsEnabled       bool
	HttpsPort          uint16
	TLSKeyFile         string
	TLSCertificateFile string
	CACertificateFile  string
	CORSEnabled        bool
}

type OpenApi struct {
	Enabled   bool
	Directory string
}

type Config struct {
	BindAddress      *net.UDPAddr `conf:"bind.address"`
	BroadcastAddress *net.UDPAddr `conf:"broadcast.address"`
	Devices          DeviceMap    `conf:"/^UT0311-L0x\\.([0-9]+)\\.(.*)/"`
	REST             `conf:"/^rest\\.(.*)/"`
	OpenApi          `conf:"/^openapi\\.(.*)/"`
}

func NewConfig() *Config {
	bind, broadcast := DefaultIpAddresses()

	c := Config{
		BindAddress:      &bind,
		BroadcastAddress: &broadcast,
		Devices:          make(map[uint32]*Device),

		REST: REST{
			HttpEnabled:        false,
			HttpPort:           8080,
			HttpsEnabled:       true,
			HttpsPort:          8443,
			TLSKeyFile:         "uhppoted.key",
			TLSCertificateFile: "uhppoted.cert",
			CACertificateFile:  "ca.cert",
			CORSEnabled:        false,
		},

		OpenApi: OpenApi{
			Enabled:   false,
			Directory: "./openapi",
		},
	}

	return &c
}

func (c *Config) Load(path string) error {
	if path == "" {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	c.OpenApi.Directory = filepath.Join(filepath.Dir(path), "rest", "openapi")

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return conf.Unmarshal(bytes, c)
}

// Ref. https://stackoverflow.com/questions/23529663/how-to-get-all-addresses-and-masks-from-local-interfaces-in-go
func DefaultIpAddresses() (net.UDPAddr, net.UDPAddr) {
	bind := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 0,
		Zone: "",
	}

	broadcast := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 60000,
		Zone: "",
	}

	copy(bind.IP, net.IPv4zero)
	copy(broadcast.IP, net.IPv4bcast)

	if ifaces, err := net.Interfaces(); err == nil {
	loop:
		for _, i := range ifaces {
			if addrs, err := i.Addrs(); err == nil {
				for _, a := range addrs {
					switch v := a.(type) {
					case *net.IPNet:
						if v.IP.To4() != nil && i.Flags&net.FlagLoopback == 0 {
							copy(bind.IP, v.IP.To4())
							if i.Flags&net.FlagBroadcast != 0 {
								addr := v.IP.To4()
								mask := v.Mask
								binary.BigEndian.PutUint32(broadcast.IP, binary.BigEndian.Uint32(addr)|^binary.BigEndian.Uint32(mask))
							}
							break loop
						}
					}
				}
			}
		}
	}

	return bind, broadcast
}

func (f *DeviceMap) UnmarshalConf(tag string, values map[string]string) (interface{}, error) {
	re := regexp.MustCompile(`^/(.*?)/$`)
	match := re.FindStringSubmatch(tag)
	if len(match) < 2 {
		return f, fmt.Errorf("Invalid 'conf' regular expression tag: %s", tag)
	}

	re, err := regexp.Compile(match[1])
	if err != nil {
		return f, err
	}

	for key, value := range values {
		match := re.FindStringSubmatch(key)
		if len(match) > 1 {
			id, err := strconv.ParseUint(match[1], 10, 32)
			if err != nil {
				return f, fmt.Errorf("Invalid 'testMap' key %s: %v", key, err)
			}

			d, ok := (*f)[uint32(id)]
			if !ok || d == nil {
				d = &Device{
					Door: make([]string, 4),
				}

				(*f)[uint32(id)] = d
			}

			switch match[2] {
			case "address":
				address, err := net.ResolveUDPAddr("udp", value)
				if err != nil {
					return f, fmt.Errorf("Device %v, invalid address '%s': %v", id, value, err)
				} else {
					d.Address = &net.UDPAddr{
						IP:   make(net.IP, net.IPv4len),
						Port: address.Port,
						Zone: "",
					}

					copy(d.Address.IP, address.IP.To4())
				}

			case "door.1":
				d.Door[0] = value

			case "door.2":
				d.Door[1] = value

			case "door.3":
				d.Door[2] = value

			case "door.4":
				d.Door[3] = value
			}
		}
	}

	return f, nil
}

func (f *REST) UnmarshalConf(tag string, values map[string]string) (interface{}, error) {
	re := regexp.MustCompile(`^/(.*?)/$`)
	match := re.FindStringSubmatch(tag)
	if len(match) < 2 {
		return f, fmt.Errorf("Invalid 'conf' regular expression tag: %s", tag)
	}

	re, err := regexp.Compile(match[1])
	if err != nil {
		return f, err
	}

	for key, value := range values {
		match := re.FindStringSubmatch(key)
		if len(match) > 1 {
			switch match[1] {
			case "CORS.enabled":
				if value == "true" {
					f.CORSEnabled = true
				} else if value == "false" {
					f.CORSEnabled = false
				} else {
					return f, fmt.Errorf("Invalid rest.CORS.enabled value: %s:", value)
				}

			case "http.enabled":
				if value == "true" {
					f.HttpEnabled = true
				} else if value == "false" {
					f.HttpEnabled = false
				} else {
					return f, fmt.Errorf("Invalid rest.http.enabled value: %s:", value)
				}

			case "http.port":
				i, err := strconv.ParseUint(value, 10, 16)
				if err != nil {
					return f, err
				}
				f.HttpPort = uint16(i)

			case "https.enabled":
				if value == "true" {
					f.HttpsEnabled = true
				} else if value == "false" {
					f.HttpsEnabled = false
				} else {
					return f, fmt.Errorf("Invalid rest.https.enabled value: %s:", value)
				}

			case "https.port":
				i, err := strconv.ParseUint(value, 10, 16)
				if err != nil {
					return f, err
				}
				f.HttpsPort = uint16(i)

			case "tls.key":
				f.TLSKeyFile = value

			case "tls.certificate":
				f.TLSCertificateFile = value

			case "tls.ca":
				f.CACertificateFile = value
			}
		}
	}

	return f, nil
}

func (f *OpenApi) UnmarshalConf(tag string, values map[string]string) (interface{}, error) {
	re := regexp.MustCompile(`^/(.*?)/$`)
	match := re.FindStringSubmatch(tag)
	if len(match) < 2 {
		return f, fmt.Errorf("Invalid 'conf' regular expression tag: %s", tag)
	}

	re, err := regexp.Compile(match[1])
	if err != nil {
		return f, err
	}

	for key, value := range values {
		match := re.FindStringSubmatch(key)
		if len(match) > 1 {
			switch match[1] {
			case "enabled":
				if value == "true" {
					f.Enabled = true
				} else if value == "false" {
					f.Enabled = false
				} else {
					return f, fmt.Errorf("Invalid rest.openapi.enabled value: %s:", value)
				}

			case "directory":
				f.Directory = value

			}
		}
	}

	return f, nil
}