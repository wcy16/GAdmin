package config

import (
	"encoding/json"
	"fmt"
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
	Commands []*Command
}

type Command struct {
	Name        string
	Description string
	SQL         string
	Query       bool
	Params      int
	Comments    []string
	Input       []interface{} `json:"-"`
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
	} else {
		for _, c := range SETTING.Commands {
			compileCommand(c)
		}
	}
}

var cmdInput = `
<div class="input-group" style="margin: 5px;">
<span class="input-group-prepend"><a href="#" class="btn btn-default disabled">%s</a></span>
<input type="text" class="form-control" id="execmd-%d" required></input>
</div>
`

func compileCommand(c *Command) {
	l := len(c.Comments)
	for i := 0; i < c.Params; i++ {
		if i >= l {
			c.Comments = append(c.Comments, "")
		}
		input := fmt.Sprintf(cmdInput, c.Comments[i], i)
		c.Input = append(c.Input, input)
	}
}
