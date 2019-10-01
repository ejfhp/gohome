package main

import (
	"fmt"
	"os"
	"testing"
)

func TestExecuteCommand1(t *testing.T) {
	commandline := []string{"LIGHT", "TURN_OFF", "kitchen.main"}
	if err := executeCommand(commandline); err != nil {
		t.Errorf("Command failed due to: %v", err)
	}
}
func TestExecuteCommand2(t *testing.T) {
	commandline := []string{"LIGHT", "ON_30_SEC", "kitchen.main"}
	if err := executeCommand(commandline); err != nil {
		t.Errorf("Command failed due to: %v", err)
	}
}

func TestMainDo(t *testing.T) {
	fmt.Printf("Start main..\n")
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	os.Args = []string{"gohome", "do", "LIGHT", "TURN_OFF", "kitchen.main"}
	main()
	fmt.Printf("Runned\n")
}

func TestMainShow(t *testing.T) {
	fmt.Printf("Start main..\n")
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	os.Args = []string{"gohome", "show"}
	main()
	fmt.Printf("Runned\n")
}

func TestMainPlant(t *testing.T) {
	fmt.Printf("Start main..\n")
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	os.Args = []string{"gohome", "plant"}
	main()
	fmt.Printf("Runned\n")
}

func TestMainHelp(t *testing.T) {
	fmt.Printf("Start main..\n")
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	os.Args = []string{"gohome", "help"}
	main()
	fmt.Printf("Runned\n")
}

func TestMainListen(t *testing.T) {
	fmt.Printf("Start main..\n")
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	os.Args = []string{"gohome", "listen"}
	main()
	fmt.Printf("Runned\n")
}
func TestMainListenTelegram(t *testing.T) {
	fmt.Printf("Start main..\n")
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	os.Args = []string{"gohome", "listen"}
	main()
	fmt.Printf("Runned\n")
}

//gcloud pubsub topics publish calling_home --message='{"who":"LIGHT","what:"TURN_ON", "where":"cucina.main","kind":"COMMAND"}'
func TestRemoteControl(t *testing.T) {
	fmt.Printf("Start main..\n")
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	os.Args = []string{"gohome", "remote"}
	main()
	fmt.Printf("Runned\n")
}

l