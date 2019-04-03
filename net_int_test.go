package gohome_test

import (
	"fmt"
	"testing"

	"github.com/savardiego/gohome"
)

func TestNewCable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	_, ok := gohome.NewCable("192.168.28.35:20000")
	if !ok {
		t.Logf("New Cable contruction failed.")
		t.Fail()
	}
}

func TestNewHome(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	c, ok := gohome.NewCable("192.168.28.35:20000")
	if !ok {
		t.Logf("New Cable contruction failed.")
		t.Fail()
	}
	plant := makeTestPlant(t)
	h := gohome.NewHome(c, plant)
	if h == nil {
		t.Logf("New Home contruction failed.")
		t.Fail()
	}
}
func TestDoTurnOn(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	c, ok := gohome.NewCable("192.168.28.35:20000")
	if !ok {
		t.Logf("New Cable contruction failed.")
		t.Fail()
	}
	plant := makeTestPlant(t)
	h := gohome.NewHome(c, plant)
	if h == nil {
		t.Logf("New Home contruction failed.")
		t.Fail()
	}
	//const cmd = "*1*18*71##"
	const cmd = "*1*0*11##"
	if !h.Do(cmd) {
		t.Logf("Send message failed failed.")
		t.Fail()
	}
}

func TestAsk(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	c, ok := gohome.NewCable("192.168.28.35:20000")
	if !ok {
		t.Logf("New Cable contruction failed.")
		t.Fail()
	}
	plant := makeTestPlant(t)
	h := gohome.NewHome(c, plant)
	if h == nil {
		t.Logf("New Home contruction failed.")
		t.Fail()
	}
	const query = "*#1*56##"
	answer := h.Ask(query)
	if len(answer) < 1 {
		t.Logf("Query failed.")
		t.Fail()
	}
	fmt.Println(answer)
}

func TestAskMany(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	c, ok := gohome.NewCable("192.168.28.35:20000")
	if !ok {
		t.Logf("New Cable contruction failed.")
		t.Fail()
	}
	plant := makeTestPlant(t)
	h := gohome.NewHome(c, plant)
	if h == nil {
		t.Logf("New Home contruction failed.")
		t.Fail()
	}
	const query = "*#1*0##"
	answer := h.Ask(query)
	if len(answer) < 1 {
		t.Logf("Query failed.")
		t.Fail()
	}
	fmt.Println(answer)
}
