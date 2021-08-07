package main

import (
	"fmt"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type MyConfig struct {
	Version int
	Name    string
	Tags    []string
}

var doc = `
version = 2
name = "go-toml"
tags = ["go", "toml"]
`

func main() {
	r := strings.NewReader(doc)
	d := toml.NewDecoder(r)
	d.SetStrict(true)

	var cfg MyConfig
	err := d.Decode(&cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println("version:", cfg.Version)
	fmt.Println("name:", cfg.Name)
	fmt.Println("tags:", cfg.Tags)
}
