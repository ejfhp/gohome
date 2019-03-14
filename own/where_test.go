package own_test

import (
	"os"
	"strings"
	"testing"

	"github.com/savardiego/gohome/own"
)

func getTestPlant() *own.Plant {
	lightK := map[string]int{"table": 1, "main": 2}
	lightL := map[string]int{"sofa": 1, "tv": 2}
	lightR := map[string]int{"main": 1, "right": 2, "left": 3}
	ambK := own.Ambient{Num: 1, Lights: lightK}
	ambL := own.Ambient{Num: 2, Lights: lightL}
	ambR := own.Ambient{Num: 3, Lights: lightR}
	casa := own.Plant{Name: "casa", Num: 1, Address: "192.168.0.35:20000"}
	casa.Ambients = map[string]own.Ambient{"kitchen": ambK, "living": ambL, "bedroom": ambR}
	return &casa
}

func TestNewPlant(t *testing.T) {
	config, err := os.Open("casa.plant")
	if err != nil {
		t.Errorf("cannot open json file")
	}
	defer config.Close()
	plant := own.NewPlant(config)
	if len(plant.Ambients) != 2 || strings.Compare(plant.Name, "home") != 0 || len(plant.Ambients["kitchen"].Lights) != 2 {
		t.Errorf("Import plant configuratin has failed: %v", plant)
	}
	if strings.Compare(plant.Address, "192.168.0.35:20000") != 0 {
		t.Errorf("Import plant configuratin has wrong address")
	}
}

func TestAddressOfLight(t *testing.T) {
	casa := getTestPlant()
	exp := map[string][]string{
		"11": []string{"kitchen", "table"},
		"12": []string{"kitchen", "main"},
		"21": []string{"living", "sofa"},
		"22": []string{"living", "tv"},
	}
	for k, v := range exp {
		w := casa.AddressOfLight(v[0], v[1])
		if k != string(w) {
			t.Errorf("Wrong where %s instead of %s", w, k)
		}
	}
}

func TestAddressOfAmb(t *testing.T) {
	casa := getTestPlant()
	exp := map[string][]string{
		"1": []string{"kitchen"},
		"2": []string{"living"},
		"0": []string{"gen"},
	}
	for k, v := range exp {
		w := casa.AddressOfAmb(v[0])
		if k != string(w) {
			t.Errorf("Wrong where %s instead of %s", w, k)
		}
	}
}

func TestExport(t *testing.T) {
	casa := getTestPlant()
	f, err := os.OpenFile("export.plant", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Errorf("Cannot open file: %v", err)
	}
	defer f.Close()
	err = casa.Export(f)
	if err != nil {
		t.Errorf("Plant export failed: %v", err)
	}

}
