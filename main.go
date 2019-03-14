package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/savardiego/gohome/own"
)

const defaultConf = "gohome.plant"

func main() {
	//command line must be WHO WHAT WHERE
	if len(os.Args) == 1 {
		basicHelp()
		return
	}
	if os.Args[1] == "help" {
		advancedHelp(os.Args[2:])
		return
	}
	plant, err := getPlant(defaultConf)
	cmd := os.Args[1:]
	if len(cmd) < 3 {
		basicHelp()
		return
	}
	whoArg := cmd[0]
	whatArg := cmd[1]
	whereArg := cmd[2]
	who := own.NewWho(whoArg)
	what := who.NewWhat(whatArg)
	where := plant.AddressOfLight() // unificare.. se esiste solo ambient è ambient, altrimenti è luce

}

func basicHelp() {
	fmt.Printf("GoHome,\n")
	fmt.Printf("a simple command line tool to control a Bticino MyHome plant.\n")
	fmt.Printf("\n")
	fmt.Printf("For extented help:\n")
	fmt.Printf("     %s help\n", os.Args[0])
}

func advancedHelp(pars []string) {
	fmt.Printf("----- advanced help\n")
	fmt.Printf("    gohome  <who> <what> <where>\n")
}

func getPlant(file string) (*own.Plant, error) {
	config, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open configuration file: %s", file)
	}
	defer config.Close()
	plant := own.NewPlant(config)
	return plant, nil
}
