package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/jxlxx/nats-gen/micro"
	"github.com/jxlxx/nats-gen/micro/spec"
)

var (
	config string
)

func init() {
	flag.StringVarP(&config, "config", "c", "", "path to config file")
}

func main() {
	flag.Parse()
	var s spec.Spec
	read(config, &s)
	g := micro.New()
	files, err := g.Run(s)
	if err != nil {
		log.Fatalln(err)
	}
	if err := g.Write(); err != nil {
		log.Fatalln(err)
	}
	for _, f := range files {
		cmd := exec.Command("go", "fmt", f)
		if err := cmd.Run(); err != nil {
			log.Fatalln(err)
		}
	}
}

func read(filename string, conf interface{}) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("err getting working dir: ", err)
	}
	fullPath := fmt.Sprintf("%s/%s", wd, filename)
	f, err := os.ReadFile(fullPath)
	if err != nil {
		log.Fatalln("err reading yaml: ", err)
	}
	err = yaml.Unmarshal(f, conf)
	if err != nil {
		log.Fatalln("err unmarshal: ", err)
	}
}
