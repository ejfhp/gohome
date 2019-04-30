package gohome_test

import (
	"testing"

	"github.com/savardiego/gohome"
)

func TestNewWho(t *testing.T) {
	expected := "1"
	who := gohome.NewWho("light")
	if who.Text() != expected {
		t.Errorf("Wrong WHO")
	}
}

func TestNewWhat(t *testing.T) {
	expected := "0"
	who := gohome.NewWho("light")
	what := who.NewWhat("TURN_OFF")
	if what.Text() != expected {
		t.Errorf("Wrong WHAT")
	}
}

func TestNewCommand(t *testing.T) {
	plant := makeTestPlant(t)
	who := gohome.NewWho("light")
	what := who.NewWhat("turn_on")
	where, err := plant.NewWhere("kitchen.table")
	if err != nil {
		t.Errorf("Where not found: %v", err)
	}
	command := gohome.NewCommand(who, what, where)
	expected := gohome.Message("*1*1*11##")
	if command != expected {
		t.Errorf("Wrong command %s, expected was %s", command, expected)
	}
}

func TestWhereFromMessage(t *testing.T) {
	messages := [][]string{
		{"*1*1*23##", "23"},
		{"*1*0*13##", "13"},
		{"*1*0*1##", "1"},
		{"*1*21##", "21"},
		{"*1*2##", "2"},
	}
	for i, m := range messages {
		msg := gohome.Message(m[0])
		expWhere := gohome.Where(m[1])
		wher := msg.Where()
		if wher != expWhere {
			t.Errorf("%d - Wrong where decoded: exp:%s actual:%s", i, expWhere, wher)
		}
	}
}
