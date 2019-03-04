package gohome_test

import (
	"fmt"
	"testing"

	"github.com/savardiego/gohome"
)

func TestNewCommand(t *testing.T) {
	plant := getTestPlant()
	command := gohome.NewCommand(gohome.Lightning, gohome.ON_30_SEC, plant.AddressOfLight("kitchen", "table"))
	fmt.Printf("Value %s\n", command)
}
