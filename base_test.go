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
		{"", ""},
	}
	for i, m := range messages {
		msg := gohome.Message(m[0])
		expWhere := gohome.Where(m[1])
		wher := msg.Where()
		if wher != expWhere {
			t.Errorf("%d - Wrong WHERE decoded: exp:%s actual:%s", i, expWhere, wher)
		}
	}
}

func TestWhoFromMessage(t *testing.T) {
	messages := [][]string{
		{"*5*1*23##", "5"},
		{"*5*0*13##", "5"},
		{"*5*0*1##", "5"},
		{"*5*21##", "5"},
		{"*55*2##", "55"},
		{"", ""},
	}
	for i, m := range messages {
		msg := gohome.Message(m[0])
		expWho := gohome.Who(m[1])
		who := msg.Who()
		if who != expWho {
			t.Errorf("%d - Wrong WHO decoded: exp:%s actual:%s", i, expWho, who)
		}
	}
}

func TestWhatFromMessage(t *testing.T) {
	messages := [][]string{
		{"*5*1*23##", "1"},
		{"*5*0*13##", "0"},
		{"*5*0*1##", "0"},
		{"*5*21##", ""},
		{"*55*2##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := gohome.Message(m[0])
		expWhat := gohome.What(m[1])
		what := msg.What()
		if what != expWhat {
			t.Errorf("%d - Wrong WHAT decoded: exp:%s actual:%s", i, expWhat, what)
		}
	}
}
func TestDecodeWhoFromMessage(t *testing.T) {
	messages := [][]string{
		{"*1*1*23##", "LIGHT"},
		{"*1*0*13##", "LIGHT"},
		{"*1*11*1##", "LIGHT"},
		{"*1*18*21##", "LIGHT"},
		{"21##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := gohome.Message(m[0])
		expWho := m[1]
		who := msg.Who()
		decoded, err := gohome.DecodeWho(who)
		if decoded != expWho {
			t.Errorf("%d - Wrong WHO decoded: exp:%s actual:%s", i, expWho, decoded)
		}
		if err != nil {
			t.Logf("Cannot decode who from message: %v ", err)
		}
	}
}

func TestDecodeWhatFromMessage(t *testing.T) {
	messages := [][]string{
		{"*1*1*23##", "TURN_ON"},
		{"*1*0*13##", "TURN_OFF"},
		{"*1*11*1##", "ON_1_MIN"},
		{"*1*18*21##", "ON_0_5_SEC"},
		{"21##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := gohome.Message(m[0])
		expWhat := m[1]
		who := msg.Who()
		what := msg.What()
		decoded, err := who.DecodeWhat(what)
		if decoded != expWhat {
			t.Errorf("%d - Wrong WHERE decoded: exp:%s actual:%s", i, expWhat, decoded)
		}
		if err != nil {
			t.Logf("Cannot decode what from message: %v ", err)
		}
	}
}
func TestMessageValid(t *testing.T) {
	messages := map[string]bool{
		"*1*1*23##":  true,
		"*1*0*13##":  true,
		"*1*11*1##":  true,
		"*1*18*21##": true,
		"*#*1##":     true,
		"*99*1##":    true,
		"*99*9##":    true,
		"21##":       false,
		"*##":        false,
		"*#":         false,
		"*":          false,
		"#":          false,
		"*1*6*d##":   false,
		"":           false,
	}
	for m, e := range messages {
		msg := gohome.Message(m)
		exp := e
		if msg.IsValid() != exp {
			if exp {
				t.Errorf("Valid message has been recognized invalid: %s", msg)
			} else {
				t.Errorf("Invalid message has been recognized valid: %s", msg)
			}
		}
	}

}
