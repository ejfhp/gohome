package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/savardiego/gohome/own"
)

func main() {
	//command line must be WHO WHAT WHERE
	if len(os.Args) == 1 {
		basicHelp()
	}
}

func basicHelp() {
	fmt.Printf("GoHome,\n")
	fmt.Printf("a simple command line tool to control a Bticino MyHome plant.\n")
	fmt.Printf("\n")
	fmt.Printf("For extented help:\n")
	fmt.Printf("     %s help\n", os.Args[0])
}

func getPlant(file string) (*own.Plant, error) {
	config, err := os.Open("gohome.conf")
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open configuration file %s", file)
	}
	defer config.Close()
	plant := own.NewPlant(config)
	return plant, nil
}
