package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/savardiego/gohome"
)

const defaultConf = "gohome.json"

var ErrWrongCommand = errors.New("Wrong command, cannot execute.")
var ErrConfigNotFound = errors.New("Config file not found or not readable.")

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
	err = gohome.LoadPlant(config)
	if err != nil {
		log.Printf("cannot open configuration file: %s due to: %v", defaultConf, err)
		return ErrConfigNotFound
	}
	config.Close()
	who := gohome.NewWho(command[0])
	if who == "" {
		log.Printf("Wrong <who>: %s", command[0])
		return ErrWrongCommand
	}
	what := who.NewWhat(command[1])
	if what == "" {
		log.Printf("Wrong <what>: %s", command[0])
		return ErrWrongCommand
	}
	where, err := gohome.NewWhere(command[2])
	if err != nil {
		log.Printf("Wrong <where>: not found: %v", err)
		return ErrWrongCommand
	}
	log.Printf("executing command, who:%s what:%s where:%s\n", who, what, where)
  command := gohome.NewCommand(who, what, where)
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
