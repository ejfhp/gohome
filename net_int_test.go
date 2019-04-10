package gohome_test

import (
	"fmt"
	"testing"

	"github.com/savardiego/gohome"
)

func TestNewHome(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	plant := makeTestPlant(t)
	h := gohome.NewHome(plant)
	if h == nil {
		t.Logf("New Home contruction failed.")
		t.Fail()
	}
}
func TestDoTurnOn(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	plant := makeTestPlant(t)
	h := gohome.NewHome(plant)
	if h == nil {
		t.Logf("New Home contruction failed.")
	}
	//const cmd = "*1*18*71##"
	const cmd = "*1*0*31##"
	if err := h.Do(cmd); err != nil {
		t.Errorf("Send message failed failed: %v", err)
	}
}

func TestAsk(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	plant := makeTestPlant(t)
	h := gohome.NewHome(plant)
	if h == nil {
		t.Logf("New Home contruction failed.")
		t.Fail()
	}
	const query = "*#1*56##"
	answer, err := h.Ask(query)
	if err != nil {
		t.Errorf("Ask failed: %v", err)
	}
	if len(answer) < 1 {
		t.Errorf("Query failed.")
	}
	fmt.Println(answer)
}

func TestAskMany(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	plant := makeTestPlant(t)
	h := gohome.NewHome(plant)
	if h == nil {
		t.Logf("New Home contruction failed.")
		t.Fail()
	}
	const query = "*#1*0##"
	answer, err := h.Ask(query)
	if err != nil {
		t.Errorf("Ask failed: %v", err)
	}
	if len(answer) < 1 {
		t.Logf("Query failed.")
		t.Fail()
	}
	fmt.Println(answer)
}
