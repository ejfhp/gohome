package own_test

import (
	"testing"
)

func TestNewCommand(t *testing.T) {
	plant := getTestPlant()
	command := gohome.NewCommand(gohome.Lightning, gohome.TURN_ON, plant.AddressOfLight("kitchen", "table"))
	expected := gohome.Command("*1*1*11##")
	if command != expected {
		t.Errorf("Wrong command %s, expected was %s", command, expected)
	}
}
