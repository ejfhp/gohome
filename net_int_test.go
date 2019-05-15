package gohome_test

import (
	"fmt"
	"testing"
	"time"

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
	// const query = "*#1*0##"
	const query = "*#5##"
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

func TestListen(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	plant := makeTestPlant(t)
	h := gohome.NewHome(plant)
	if h == nil {
		t.Logf("New Home contruction failed.")
		t.Fail()
	}
	listen, stop, errs := h.Listen()
	go func() {
		dur := time.Duration(10 * time.Second)
		fmt.Printf(">>>>> Waiting %f seconds.. \n", dur.Seconds())
		time.Sleep(dur)
		// var s struct{}
		fmt.Printf(">>>>> Sending stop.. \n")
		stop <- struct{}{}
	}()
	fmt.Printf(">>>>> Ready to listen.. \n")
	ok := true
	var e error
	var m gohome.Message
	for ok == true {
		select {
		case e, ok = <-errs:
			fmt.Printf(">>>>> error received (ok? %t): %v\n", ok, e)
		case m, ok = <-listen:
			who, what, where, err := plant.Parse(m)
			fmt.Printf(">>>>> received (ok? %t): %s %s %s  -- err: %v\n", ok, who, what, where, err)
		}
	}
}
