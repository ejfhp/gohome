package gohome_test

import (
	"fmt"
	"testing"

	"github.com/savardiego/gohome"
)

func TestWhatWithValue(t *testing.T) {
	w := gohome.What("3")
	fmt.Printf("Value %s\n", w)
}
