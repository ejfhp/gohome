package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

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
	var err error
	cmd := os.Args[1]
	switch cmd {
	case "help":
		advancedHelp(os.Args[1:])
		break
	case "plant":
		err = showPlant(os.Args[1:])
		break
	case "show":
		err = showHome(os.Args[1:])
		break
	case "do":
		err := executeCommand(os.Args[2:])
		if err != nil {
			fmt.Printf("Cannot complete command executiion: %+v\n", err)
		}
		break
	case "listen":
		err = listen()
		break
	case "remote":
		err = remoteControl()
		break
	default:
		basicHelp()
		break
	}
	if err != nil {
		fmt.Printf("Unfortunately something went wrong: %v\n", err)
	}
}

func executeCommand(command []string) error {
	home, err := openHome()
	if err != nil {
		return errors.Wrapf(err, "cannot open Home")
	}
	fmt.Printf("who is %s\n", command[0])
	who := gohome.NewWho(command[0])
	if who.Desc == "" {
		return errors.Errorf("unknown <who> in command:%s", command[0])
	}
	fmt.Printf("what is %s\n", command[1])
	what, err := who.WhatFromDesc(command[1])
	if err != nil {
		return errors.Errorf("Cannot get <what> from command: %s due to: %v", command[1], err)
	}
	fmt.Printf("where is %s\n", command[2])
	where, err := home.Plant.WhereFromDesc(command[2])
	if err != nil {
		return errors.Wrapf(err, "Cannot get <where> from command: %s due to: %v", command[2], err)
	}
	fmt.Printf("executing command, who:%s what:%s where:%s\n", who.Desc, what.Desc, where.Desc)
	cmd := gohome.NewCommand(who, what, where)
	return home.Do(cmd)
}

func showPlant(command []string) error {
	home, err := openHome()
	if err != nil {
		return errors.Wrapf(err, "cannot open Home")
	}
	fmt.Println("-------------------")
	fmt.Printf("Plant: %s\n\n", home.Plant.Name)
	fmt.Printf("Ambients:\n")
	for a, amb := range home.Plant.Ambients {
		fmt.Printf("     %s: %d\n", a, amb.Num)
		for l, n := range amb.Lights {
			fmt.Printf("          %s: %d\n", l, n)
		}
	}
	return nil
}

func showHome(command []string) error {
	home, err := openHome()
	if err != nil {
		return errors.Wrapf(err, "cannot open Home")
	}
	queryStatus := gohome.SystemMessages["QUERY_ALL"]
	statuses, err := home.Ask(queryStatus)
	if err != nil {
		return errors.Wrapf(err, "cannot get plant status, queryFrame: %s", queryStatus.Kind)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(w, "N\tWHO\tWHERE\tSTATUS")
	for i, m := range statuses {
		if err != nil {
			fmt.Printf("failed to decode message '%v' due to: %v", m, err)
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", i, m.Who.Desc, m.Where.Desc, m.What.Desc)
	}
	w.Flush()
	return nil
}

func listen() error {
	home, err := openHome()
	if err != nil {
		return errors.Wrapf(err, "cannot open Home")
	}
	listen, _, errs := home.Listen()
	ok := true
	for ok {
		select {
		case e, ok := <-errs:
			fmt.Printf(">>>>> error received (ok? %t): %v\n", ok, e)
		case f, ok := <-listen:
			if v, _ := gohome.IsValid(f); v {
				msg := home.Plant.ParseFrame(f)
				fmt.Printf(">>>>> received (ok? %t): '%s' '%s' '%s'  msg: '%v'\n", ok, msg.Who.Desc, msg.What.Desc, msg.Where.Desc, msg.Kind)
			} else {
				fmt.Printf(">>>>> message invalid: '%s'\n", f)
			}
		}
	}
	return nil
}

func openPlantFile() (*os.File, error) {
	gohomePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	gohomeDir := filepath.Dir(gohomePath)
	plantFilePath := filepath.Join(gohomeDir, defaultConf)
	config, err := os.Open(plantFilePath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func openHome() (*gohome.Home, error) {
	config, err := openPlantFile()
	defer config.Close()
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open configuration file: %s", defaultConf)
	}
	plant, err := gohome.NewPlant(config)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot load plant from configuration file: %s", defaultConf)
	}
	home := gohome.NewHome(plant)
	return home, nil
}

func remoteControl() error {
	home, err := openHome()
	if err != nil {
		return errors.Wrapf(err, "cannot open Home")
	}
	pubsub, err := gohome.NewPubSub()
	if err != nil {
		return errors.Wrapf(err, "cannot access Google Pub/Sub")
	}
	incoming, errs := pubsub.Listen(home)
	for true {
		select {
		case inMsg := <-incoming:
			fmt.Printf("Command from remote %s, %v \n", inMsg.Frame(), inMsg)
			home.Do(inMsg)
			break
		case err := <-errs:
			return errors.Wrapf(err, "errors while listening to the Pub/Sub")
		}
	}
	return nil
}

func basicHelp() {
	fmt.Printf("GoHome,\n")
	fmt.Printf("a simple command line tool to control a Bticino MyHome plant.\n")
	fmt.Printf("\n")
	fmt.Printf("Basic commands:\n")
	fmt.Printf("     %s help: extended help\n", os.Args[0])
	fmt.Printf("     %s plant: print current plant from file gohome.json\n", os.Args[0])
	fmt.Printf("     %s show: show status of all home components\n", os.Args[0])
	fmt.Printf("     %s listen: listen to network and show events\n", os.Args[0])
	fmt.Printf("     %s do: listen to network and show events\n", os.Args[0])
}

func advancedHelp(pars []string) {
	fmt.Printf("ADVANCED HELP\n")
	fmt.Printf("      default configuration file is \"gohome.json\"\n\n")
	fmt.Printf("      %s  <who> <what> <where>\n", os.Args[0])
	fmt.Printf("        who:   LIGHT (currently work only on lights)\n")
	fmt.Printf("        what:  <command>\n")
	fmt.Printf("        where: <room>.<light> (in case of single light)\n")
	fmt.Printf("        where: <room>         (in case of ambient)\n")
	fmt.Printf("        where: general        (in case of general)\n")
	fmt.Printf("\n\nFor LIGHT <command> is one of:\n")
	for _, v := range gohome.NewWho("LIGHT").Actions {
		fmt.Printf("      %v\n", v)
	}
}
