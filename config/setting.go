package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Database struct {
	Username string
	Password string
	DBname   string
}

type Setting struct {
	Database Database
	Commands []Command
}

type Command struct {
	Name string
	SQL  string
}

var SETTING Setting

func LoadSetting(file string) {
	setting, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer setting.Close()

	byte, _ := ioutil.ReadAll(setting)

	if err = json.Unmarshal(byte, &SETTING); err != nil {
		panic(err)
	}
}
