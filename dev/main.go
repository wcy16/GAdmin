package main

import (
	"flag"
	"gadmin"
	"github.com/go-bindata/go-bindata"
	"path/filepath"
)

func main() {
	var c bool

	flag.BoolVar(&c, "compile", false, "compile assets")

	if c {
		compile("static", "../static/", "../static/", "../static/static.go", true, true)
		compile("template", "../template", "../template/", "../template/template.go", false, false)
	} else {
		gadmin.Serve("../config/setting.json")
	}
}

func compile(pkg, prefix, input, output string, httpFileSystem, recursive bool) {
	cfg := bindata.NewConfig()

	cfg.Package = pkg
	cfg.HttpFileSystem = httpFileSystem
	cfg.Prefix = prefix
	cfg.Output = output
	cfg.Input = make([]bindata.InputConfig, 1)
	cfg.Input[0] = bindata.InputConfig{
		Path:      filepath.Clean(input),
		Recursive: recursive,
	}

	cfg.Dev = false

	err := bindata.Translate(cfg)

	if err != nil {
		panic(err)
	}
}
