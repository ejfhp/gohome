package gohome_test

import (
	"fmt"
	"testing"

	"github.com/savardiego/gohome"
)

/*
TO AVOID THESE TESTS WHILE NOT PRESENT A MYHOME SERVER ON THE NETWORK RUN:

go test -short

*/

//gcloud pubsub topics publish calling_home --message={}
func TestPubSubListen(t *testing.T) {
	plant := makeTestPlant(t)
	h := gohome.NewHome(plant)
	if h == nil {
		t.Logf("New Home contruction failed.")
		t.Fail()
	}
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	pubsub, err := gohome.NewPubSub()
	if err != nil {
		t.Errorf("Failed to create pubsub %v", err)
	}
	incoming, errors := pubsub.Listen(h)
	for true {
		select {
		case inMsg := <-incoming:
			fmt.Printf("Got message %s \n", inMsg.Frame())
			break
		case err := <-errors:
			t.Errorf("Got error: %v \n", err)
		}
	}

}
