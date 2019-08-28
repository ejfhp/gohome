package gohome_test

import (
	"testing"

	"github.com/savardiego/gohome"
)

/*
TO AVOID THESE TESTS WHILE NOT PRESENT A MYHOME SERVER ON THE NETWORK RUN:

go test -short

*/

func SkipTestPubSubListen(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
	}
	pubsub, err := gohome.NewPubSub()
	if err != nil {
		t.Errorf("Failed to create pubsub %v", err)
	}
	pubsub.Listen()
}
