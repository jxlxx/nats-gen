package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"

	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/jxlxx/nats-gen"
	"github.com/jxlxx/nats-gen/spec"
)

var (
	config      string
	help        bool
	showVersion bool
)

func init() {
	flag.StringVarP(&config, "config", "c", "", "path to config file")
	flag.BoolVarP(&help, "help", "h", false, "print help message and exit")
	flag.BoolVarP(&showVersion, "version", "v", false, "print version and exit")
}

func main() {
	flag.Parse()
	if help {
		printHelp()
		return
	}
	if showVersion {
		fmt.Println(getVersion())
		return
	}

	var s spec.Spec
	read(config, &s)
	g := natsgen.New()
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

func goInstallVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	return info.Main.Version
}

func getVersion() string {
	if natsgen.Version != "" {
		return natsgen.Version
	}
	return goInstallVersion()
}

func printHelp() {
	fmt.Println(("TODO"))
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
