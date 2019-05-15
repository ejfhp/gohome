package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/savardiego/gohome"
)

const defaultConf = "gohome.json"

func main() {
	//command line must be WHO WHAT WHERE
	if len(os.Args) < 2 {
		basicHelp()
		return
	}
	cmd := os.Args[1]
	switch cmd {
	case "help":
		advancedHelp(os.Args)
		break
	case "light":
		executeCommand(os.Args[2:])
		break
	default:
		basicHelp()
		break
	}
}

func executeCommand(command []string) error {
	config, err := os.Open(defaultConf)
	if err != nil {
		return errors.Wrapf(err, "cannot open configuration file: %s", defaultConf)
	}
	plant, err := gohome.LoadPlant(config)
	if err != nil {
		return errors.Wrapf(err, "cannot load plant from configuration file: %s", defaultConf)
	}
	config.Close()
	who := gohome.NewWho(command[0])
	if who == "" {
		return errors.Errorf("unknown <who> in command:%s", command[0])
	}
	what := who.NewWhat(command[1])
	if what == "" {
		errors.Errorf("unknown <what> in command: %s", command[1])
	}
	where, err := plant.NewWhere(command[2])
	if err != nil {
		errors.Wrapf(err, "wrong <where> in command: %s", command[2])
	}
	log.Printf("executing command, who:%s what:%s where:%s\n", who, what, where)
	cmd := gohome.NewCommand(who, what, where)
	home := gohome.NewHome(plant)
	return home.Do(cmd)
}

func listen() error {
	config, err := os.Open(defaultConf)
	if err != nil {
		return errors.Wrapf(err, "cannot open configuration file: %s", defaultConf)
	}
	plant, err := gohome.LoadPlant(config)
	if err != nil {
		return errors.Wrapf(err, "cannot load plant from configuration file: %s", defaultConf)
	}
	config.Close()
	home := gohome.NewHome(plant)
	home.Listen()
	return nil
}

func basicHelp() {
	fmt.Printf("GoHome,\n")
	fmt.Printf("a simple command line tool to control a Bticino MyHome plant.\n")
	fmt.Printf("\n")
	fmt.Printf("For extented help:\n")
	fmt.Printf("     %s help\n", os.Args[0])
}

func advancedHelp(pars []string) {
	fmt.Printf("ADVANCED HELP\n")
	fmt.Printf("      default configuration file is \"gohome.json\"\n\n")
	fmt.Printf("      %s  <who> <what> <where>\n", os.Args[0])
	fmt.Printf("        who= light\n")
	fmt.Printf("        what= <command>\n")
	fmt.Printf("        where= <room>.<light> (in case of single light)\n")
	fmt.Printf("        where= <room>         (in case of ambient)\n")
	fmt.Printf("        where= general        (in case of general)\n")
	fmt.Printf("\n\nCOMMANDS\n   LIGHT\n")
	for k := range gohome.WhoWhat[gohome.Light] {
		fmt.Printf("      %s\n", k)
	}
}
