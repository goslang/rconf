package main

import (
	"fmt"
	"github.com/goslang/rconf"
	"os"
)

var conf struct {
	Host string
	Port string
	Name string
}

func main() {
	interpreter, err := rconf.NewInterpreter(dsl)

	file, err := os.Open("examples/ez.rb")
	if err != nil {
		panic(err)
	}

	err = interpreter(file)
	fmt.Println(err)
	fmt.Println(conf)
}

func dsl(bc rconf.BindContext) {
	bc.Block("database", func(bc rconf.BindContext) {
		bc.BindString("host", &conf.Host)
		bc.BindString("port", &conf.Port)
		bc.BindString("name", &conf.Name)
	})
}
