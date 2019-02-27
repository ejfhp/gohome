package gohome_test

import (
	"os"
	"strings"
	"testing"

	"github.com/savardiego/gohome"
)

func TestNewPlant(t *testing.T) {
	casa, err := os.Open("casa.gho")
	if err != nil {
		t.Errorf("cannot open json file")
	}
	defer casa.Close()
	plant := gohome.NewPlant(casa)
	if len(plant.Ambients) != 2 || strings.Compare(plant.Name, "home") != 0 || len(plant.Ambients["kitchen"].Lights) != 2 {
		t.Errorf("Import plant configuratin has failed: %v", plant)
	}
	if strings.Compare(plant.Address, "192.168.0.35:20000") != 0 {
		t.Errorf("Import plant configuratin has wrong address")

	}
}

func TestExport(t *testing.T) {
	light := map[string]int{"luce1": 1, "luce2": 2}
	amb1 := gohome.Ambient{Num: 1, Lights: light}
	amb2 := gohome.Ambient{Num: 2, Lights: light}
	casa := &gohome.Plant{Name: "casa", Num: 1, Address: "192.168.0.35:20000"}
	casa.Ambients = map[string]gohome.Ambient{"kitchen": amb1, "living": amb2}
	f, err := os.OpenFile("export.gho", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Errorf("Cannot open file: %v", err)
	}
	defer f.Close()
	err = casa.Export(f)
	if err != nil {
		t.Errorf("Plant export failed: %v", err)
	}

}
