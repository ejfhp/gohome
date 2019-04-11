package main

import "testing"

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
