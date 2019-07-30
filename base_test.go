package gohome_test

import (
	"strconv"
	"testing"

	"github.com/savardiego/gohome"
)

func TestNewWho(t *testing.T) {
	expected := "LIGHT"
	who := gohome.NewWho("1")
	if who.Desc != expected {
		t.Errorf("Wrong WHO")
	}
}

func TestNewWhat(t *testing.T) {
	expected := "0"
	who := gohome.NewWho("1")
	what, err := who.WhatFromDesc("TURN_OFF")
	if err != nil {
		t.Errorf("What not found: %v", err)
	}
	if what.Code != expected {
		t.Errorf("Wrong WHAT")
	}
}

func TestNewCommand(t *testing.T) {
	plant := makeTestPlant(t)
	who := gohome.NewWho("1")
	what, err := who.WhatFromDesc("TURN_ON")
	if err != nil {
		t.Errorf("What not found: %v", err)
	}
	where, err := plant.WhereFromDesc("kitchen.table")
	if err != nil {
		t.Errorf("Where not found: %v", err)
	}
	command := gohome.NewCommand(who, what, where)
	frame := command.Frame()
	expectedFrame := "*1*1*11##"

	if frame != expectedFrame {
		t.Errorf("Wrong command %s, expected was %s", frame, expectedFrame)
	}
}

func TestMessageIsValid(t *testing.T) {
	messages := map[string][]string{
		"*1*1*23##":      []string{"TRUE", "COMMAND"},
		"*1*0*13##":      []string{"TRUE", "COMMAND"},
		"*1*11*1##":      []string{"TRUE", "COMMAND"},
		"*1*18*21##":     []string{"TRUE", "COMMAND"},
		"*#1*2##":        []string{"TRUE", "REQUEST"},
		"*#1*18*10##":    []string{"TRUE", "DIMENSIONGET"},
		"*#1*18*#10*5##": []string{"TRUE", "DIMENSIONSET"},
		"*#*1##":         []string{"TRUE", "SPECIAL"},
		"*99*1##":        []string{"TRUE", "SPECIAL"},
		"*99*9##":        []string{"TRUE", "SPECIAL"},
		"21##":           []string{"FALSE", "INVALID"},
		"*##":            []string{"FALSE", "INVALID"},
		"*#":             []string{"FALSE", "INVALID"},
		"*":              []string{"FALSE", "INVALID"},
		"#":              []string{"FALSE", "INVALID"},
		"*1*6*d##":       []string{"FALSE", "INVALID"},
		"":               []string{"FALSE", "INVALID"},
	} /// validity
	for m, e := range messages {
		valid, _ := strconv.ParseBool(e[0])
		if val, kind := gohome.IsValid(m); val != valid || kind != e[1] {
			t.Errorf("Wrong validity or vrong kind: %s, got valid:%t kind:%s", m, val, kind)
		}
	}
}

func TestMessageIsRequest(t *testing.T) {
	messages := map[string]string{
		"*1*1*23##":         gohome.COMMAND,
		"*1*0*13##":         gohome.COMMAND,
		"*1*11*1##":         gohome.COMMAND,
		"*#1*18*21##":       gohome.DIMENSIONGET,
		"*#1*18*#21*4*78##": gohome.DIMENSIONSET,
		"*#*1##":            gohome.SPECIAL,
		"*99*1##":           gohome.SPECIAL,
		"*1*9##":            gohome.INVALID,
	}
	for m, e := range messages {
		exp := e
		if _, k := gohome.IsValid(m); k != exp {
			t.Errorf("Request failed to be recognized: %s", m)
		}
	}
}
