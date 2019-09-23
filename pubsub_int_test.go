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
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	pubsub, err := gohome.NewPubSub()
	if err != nil {
		t.Errorf("Failed to create pubsub %v", err)
	}
	msgChan := pubsub.Listen()
	for m := range msgChan {
		fmt.Printf("Got message %s, ", m)
	}

}
