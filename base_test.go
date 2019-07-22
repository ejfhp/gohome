package gohome_test

import (
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
	what, err := who.WhatFromDesc("turn_on")
	if err != nil {
		t.Errorf("What not found: %v", err)
	}
	where, err := plant.WhereFromDesc("kitchen.table")
	if err != nil {
		t.Errorf("Where not found: %v", err)
	}
	frame := gohome.NewCommand(who, what, where).Frame()
	expectedFrame := "*1*1*11##"

	if frame != expectedFrame {
		t.Errorf("Wrong command %v, expected was %v", frame, expectedFrame)
	}
}

func TestWhereFromFrame(t *testing.T) {
	plant := makeTestPlant(t)
	messages := [][]string{
		{"*5*1*23##", "23"},
		{"*#5*1*#43*8##", "1"},
		{"*5*0*13##", "13"},
		{"*5*0*1##", "1"},
		{"*#5*21##", "21"},
		{"*#5*21*2##", "21"},
		{"*#5*1##", "1"},
		{"*5*2##", ""},
		{"*5**2##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := plant.ParseFrame(m[0])
		expWhereCode := m[1]
		if msg.Where.Code != expWhereCode {
			t.Errorf("%d - Wrong WHERE decoded: exp:%s actual:'%s'", i, expWhereCode, msg.Where.Code)
		}
	}
}

func TestWhoFromFrame(t *testing.T) {
	messages := [][]string{
		{"*5*1*23##", "5"},
		{"*#5*1*#43*8##", "5"},
		{"*5*0*13##", "5"},
		{"*5*0*1##", "5"},
		{"*#5*21##", "5"},
		{"*#5*21*2##", "5"},
		{"*#5*1##", "5"},
		{"*5*2##", ""},
		{"*5**2##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := gohome.ParseFrame(m[0])
		expWho := gohome.Who(m[1])
		if msg.Who != expWho {
			t.Errorf("%d - Wrong WHO decoded: exp:%s actual:%s", i, expWho, msg.Who)
		}
	}
}

func TestWhatFromFrame(t *testing.T) {
	messages := [][]string{
		{"*5*1*23##", "1"},
		{"*#5*1*#43*8##", ""},
		{"*5*0*13##", "0"},
		{"*5*10*13##", "10"},
		{"*5*0*1##", "0"},
		{"*#5*21##", ""},
		{"*#5*21*2##", ""},
		{"*#5*1##", ""},
		{"*5*2##", ""},
		{"*5**2##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := gohome.ParseFrame(m[0])
		expWhat := gohome.What(m[1])
		if msg.What != expWhat {
			t.Errorf("%d - Wrong WHAT decoded: exp:%s actual:%s", i, expWhat, msg.What)
		}
	}
}
func TestDecodeWhoFromFrame(t *testing.T) {
	messages := [][]string{
		{"*1*1*23##", "LIGHT"},
		{"*1*0*13##", "LIGHT"},
		{"*1*11*1##", "LIGHT"},
		{"*1*18*21##", "LIGHT"},
		{"21##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := gohome.ParseFrame(m[0])
		expWho := m[1]
		decoded, err := gohome.DecodeWho(msg.Who)
		if decoded != expWho {
			t.Errorf("%d - Wrong WHO decoded: exp:%s actual:%s", i, expWho, decoded)
		}
		if err != nil {
			t.Logf("Cannot decode who from message: %v ", err)
		}
	}
}

func TestDecodeWhatFromFrame(t *testing.T) {
	messages := [][]string{
		{"*1*1*23##", "TURN_ON"},
		{"*1*0*13##", "TURN_OFF"},
		{"*1*11*1##", "ON_1_MIN"},
		{"*1*18*21##", "ON_0_5_SEC"},
		{"21##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := gohome.ParseFrame(m[0])
		expWhat := m[1]
		decoded, err := msg.Who.DecodeWhat(msg.What)
		if decoded != expWhat {
			t.Errorf("%d - Wrong WHERE decoded: exp:%s actual:%s", i, expWhat, decoded)
		}
		if err != nil {
			t.Logf("Cannot decode what from message: %v ", err)
		}
	}
}

func TestMessageIsValid(t *testing.T) {
	messages := map[string][]int{
		"*1*1*23##":      []int{1, 0},
		"*1*0*13##":      []int{1, 0},
		"*1*11*1##":      []int{1, 0},
		"*1*18*21##":     []int{1, 0},
		"*#1*2##":        []int{1, 1},
		"*#1*18*10##":    []int{1, 3},
		"*#1*18*#10*5##": []int{1, 4},
		"*#*1##":         []int{1, 2},
		"*99*1##":        []int{1, 2},
		"*99*9##":        []int{1, 2},
		"21##":           []int{0, -1},
		"*##":            []int{0, -1},
		"*#":             []int{0, -1},
		"*":              []int{0, -1},
		"#":              []int{0, -1},
		"*1*6*d##":       []int{0, -1},
		"":               []int{0, -1},
	}
	for m, e := range messages {
		if v, k := gohome.IsValid(m); v != (e[0] != 0) || k != e[1] {
			t.Errorf("Wrong validity or vrong kind: %s, got valid:%t kind:%d", m, v, k)
		}
	}
}

func TestMessageIsRequest(t *testing.T) {
	messages := map[string]int{
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
