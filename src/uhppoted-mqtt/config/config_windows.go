package config

import (
	"golang.org/x/sys/windows"
	"path/filepath"
)

var secrets string = filepath.Join(workdir(), "mqtt.hotp.secrets")
var counters string = filepath.Join(workdir(), "mqtt.hotp.counters")

func workdir() string {
	programData, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return `C:\uhppoted`
	}

	return filepath.Join(programData, "uhppoted")
}