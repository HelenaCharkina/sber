package settings

import (
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
)

type Settings struct {
	Port string

	DBHost     string
	DBPort     string
	DBUsername string
	DBName     string
	DBSSLMode  string
}

var Config *Settings

func InitConfig() (err error) {
	Config = new(Settings)

	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	settingsFile, err := ini.Load(filepath.Join(pwd, "conf.ini"))
	if err != nil {
		return
	}
	err = settingsFile.MapTo(&Config)
	if err != nil {
		return
	}

	return
}
