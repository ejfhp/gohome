package main

import (
	"fmt"
	"os"
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
	cmd := os.Args[1]
	switch cmd {
	case "help":
		advancedHelp(os.Args[1:])
		break
	case "plant":
		showPlant(os.Args[1:])
		break
	case "show":
		showHome(os.Args[1:])
		break
	case "light":
		err := executeCommand(os.Args[1:])
		if err != nil {
			fmt.Printf("Cannot complete command executiion: %+v\n", err)
		}
		break
	case "listen":
		listen()
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
	fmt.Printf("who is %s\n", command[0])
	who := gohome.NewWho(command[0])
	if who == "" {
		return errors.Errorf("unknown <who> in command:%s", command[0])
	}
	fmt.Printf("what is %s\n", command[1])
	what := who.NewWhat(command[1])
	if what == "" {
		return errors.Errorf("unknown <what> in command: %s", command[1])
	}
	fmt.Printf("where is %s\n", command[2])
	where, err := plant.NewWhere(command[2])
	if err != nil {
		return errors.Wrapf(err, "wrong <where> in command: %s", command[2])
	}
	fmt.Printf("executing command, who:%s what:%s where:%s\n", who, what, where)
	cmd := gohome.NewCommand(who, what, where)
	home := gohome.NewHome(plant)
	return home.Do(cmd)
}

func showPlant(command []string) error {
	config, err := os.Open(defaultConf)
	if err != nil {
		return errors.Wrapf(err, "cannot open configuration file: %s", defaultConf)
	}
	plant, err := gohome.LoadPlant(config)
	if err != nil {
		return errors.Wrapf(err, "cannot load plant from configuration file: %s", defaultConf)
	}
	config.Close()
	fmt.Println("-------------------")
	fmt.Printf("Plant: %s\n\n", plant.Name)
	fmt.Printf("Ambients:\n")
	for a, amb := range plant.Ambients {
		fmt.Printf("     %s: %d\n", a, amb.Num)
		for l, n := range amb.Lights {
			fmt.Printf("          %s: %d\n", l, n)
		}
	}
	return nil
}

func showHome(command []string) error {
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
	const queryStatus = "*#1*0##"
	// const queryStatus = "*#5##"
	statuses, err := home.Ask(queryStatus)
	if err != nil {
		return errors.Wrapf(err, "cannot get plant status, queryStatus: %s", queryStatus)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(w, "N\tWHO\tWHERE\tSTATUS")
	for i, m := range statuses {
		who, what, where, err := plant.Parse(m)
		if err != nil {
			fmt.Printf("failed to decode message '%s' due to: %v", m, err)
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", i, who, where, what)
	}
	w.Flush()
	return nil
}

func listen() error {
	// config, err := os.Open("gohome.json")
	// if err != nil {
	// 	t.Errorf("cannot open json file")
	// }
	// defer config.Close()
	// plant, err := gohome.LoadPlant(config)
	// if err != nil {
	// 	t.Errorf("cannot load plant from config file")
	// }
	// if plant.ServerAddress() != "192.168.0.35:20000" {
	// 	t.Errorf("Import plant configuration has wrong address: '%s', len:%d", plant.ServerAddress(), len(plant.ServerAddress()))
	// }
	// plant := makeTestPlant(t)
	// h := gohome.NewHome(plant)
	// if h == nil {
	// 	t.Logf("New Home contruction failed.")
	// 	t.Fail()
	// }
	// listen, stop, errs := h.Listen()
	// config, err := os.Open(defaultConf)
	// if err != nil {
	// 	return errors.Wrapf(err, "cannot open configuration file: %s", defaultConf)
	// }
	// plant, err := gohome.LoadPlant(config)
	// if err != nil {
	// 	return errors.Wrapf(err, "cannot load plant from configuration file: %s", defaultConf)
	// }
	// config.Close()
	// home := gohome.NewHome(plant)
	// home.Listen()
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
