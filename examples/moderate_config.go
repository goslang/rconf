package main

import (
	"fmt"
	"github.com/goslang/rconf"
	"os"
)

var conf struct {
	Database struct {
		Host       string
		Port       int
		Name       string
		SearchPath string
	}

	Server struct {
		Host string
		Port int
	}

	Jobs struct {
		Concurrency int
		ServerURI   string
	}
}

func main() {
	interpreter, err := rconf.NewInterpreter(dsl)

	file, err := os.Open("examples/moderate.rb")
	if err != nil {
		panic(err)
	}

	err = interpreter(file)
	fmt.Println(err)
	fmt.Println(conf)
}

func dsl(bc rconf.BindContext) {
	bc.BlockWithArg("database", func(bc rconf.BindContext) {
		backend := bc.StringArg(0)

		bc.BindString("host", &conf.Database.Host)
		bc.BindInt("port", &conf.Database.Port)
		bc.BindString("name", &conf.Database.Name)

		if backend == "postgres" {
			bc.BindString("search_path", &conf.Database.SearchPath)
		}
	})

	bc.Block("server", func(bc rconf.BindContext) {
		bc.BindString("host", &conf.Server.Host)
		bc.BindInt("port", &conf.Server.Port)
	})

	bc.Block("jobs", func(bc rconf.BindContext) {
		bc.BindString("server_uri", &conf.Jobs.ServerURI)
		bc.BindInt("concurrency", &conf.Jobs.Concurrency)
	})
}
