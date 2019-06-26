package main

import (
	"fmt"
	"os"
	"testing"
)

func TestExecuteCommand1(t *testing.T) {
	commandline := []string{"light", "TURN_ON", "kitchen.main"}
	if err := executeCommand(commandline); err != nil {
		t.Errorf("Command failed due to: %v", err)
	}
}
func TestExecuteCommand2(t *testing.T) {
	commandline := []string{"light", "ON_30_SEC", "kitchen.main"}
	if err := executeCommand(commandline); err != nil {
		t.Errorf("Command failed due to: %v", err)
	}
}

func TestCommandMain(t *testing.T) {
	fmt.Printf("Start main..\n")
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	os.Args = []string{"gohome", "light", "TURN_OFF", "cucina.tavolo"}
	main()
	fmt.Printf("Runned\n")
}

func TestCommandShow(t *testing.T) {
	fmt.Printf("Start main..\n")
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	os.Args = []string{"gohome", "show"}
	main()
	fmt.Printf("Runned\n")
}
