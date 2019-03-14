package own_test

import (
	"testing"

	"github.com/savardiego/gohome/own"
)

func TestNewWho(t *testing.T) {
	expected := "1"
	who := own.NewWho("lightning")
	if who.Text() != expected {
		t.Errorf("Wrong WHO")
	}
}

func TestNewWhat(t *testing.T) {
	expected := "0"
	who := own.NewWho("lightning")
	what := who.NewWhat("TURN_OFF")
	if what.Text() != expected {
		t.Errorf("Wrong WHAT")
	}
}

func TestNewCommand(t *testing.T) {
	plant := getTestPlant()
	turnOn := own.What("1")
	command := own.NewCommand(own.Who("lightning"), turnOn, plant.AddressOfLight("kitchen", "table"))
	expected := own.Command("*1*1*11##")
	if command != expected {
		t.Errorf("Wrong command %s, expected was %s", command, expected)
	}
}
