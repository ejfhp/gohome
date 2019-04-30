package gohome_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/savardiego/gohome"
)

func makeTestPlant(t *testing.T) *gohome.Plant {
	buf := bytes.NewBufferString("{ \"name\": \"home\", \"address\": \"192.168.28.35:20000\", \"num\": 1, \"ambients\": { \"kitchen\": { \"num\": 1, \"Lights\": { \"table\": 1, \"main\": 2 } }, \"living\": { \"num\": 2, \"Lights\": { \"sofa\": 1, \"tv\": 2 } } } }")
	p, err := gohome.LoadPlant(buf)
	if err != nil {
		t.Errorf("LoadPlant failed: %v", err)
	}
	return p
}

func TestLoadPlant(t *testing.T) {
	config, err := os.Open("testdata/casa.json")
	if err != nil {
		t.Errorf("cannot open json file")
	}
	defer config.Close()
	plant, err := gohome.LoadPlant(config)
	if err != nil {
		t.Errorf("cannot load plant from config file")
	}
	if plant.ServerAddress() != "192.168.0.35:20000" {
		t.Errorf("Import plant configuration has wrong address: '%s', len:%d", plant.ServerAddress(), len(plant.ServerAddress()))
	}
	exp := map[string]string{
		"11": "kitchen.table",
		"12": "kitchen.main",
		"21": "living.sofa",
		"22": "living.tv",
		"1":  "kitchen",
		"2":  "living",
	}
	for k, v := range exp {
		w, err := plant.NewWhere(v)
		if k != string(w) || err != nil {
			t.Errorf("Wrong where %s instead of %s (err: %v)", w, k, err)
		}
	}
}

func TestNewWhere(t *testing.T) {
	plant := makeTestPlant(t)
	exp := map[string]string{
		"11": "kitchen.table",
		"12": "kitchen.main",
		"21": "living.sofa",
		"22": "living.tv",
		"1":  "kitchen",
		"2":  "living",
	}
	for k, v := range exp {
		w, err := plant.NewWhere(v)
		if k != string(w) || err != nil {
			t.Errorf("Wrong where %s instead of %s (err: %v)", w, k, err)
		}
	}
}

func TestDecode(t *testing.T) {
	config, err := os.Open("testdata/casa.json")
	if err != nil {
		t.Errorf("cannot open json file")
	}
	defer config.Close()
	plant, err := gohome.LoadPlant(config)
	if err != nil {
		t.Errorf("cannot load plant from config file")
	}
	exp := map[string]string{
		"11": "kitchen.table",
		"12": "kitchen.main",
		"21": "living.sofa",
		"22": "living.tv",
		"1":  "kitchen",
		"2":  "living",
	}
	for w, e := range exp {
		wh := gohome.Where(w)
		dec, err := plant.Decode(wh)
		fmt.Printf("Where:%s decoded:%s\n", wh, dec)
		if dec != e || err != nil {
			t.Errorf("Where not decoded correctly, exp:%s  decoded:%s", wh, dec)
		}
	}
}

func TestExport(t *testing.T) {
	plant := makeTestPlant(t)
	f, err := os.OpenFile("testdata/export.plant", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Errorf("Cannot open file: %v", err)
	}
	defer f.Close()
	err = plant.ExportPlant(f)
	if err != nil {
		t.Errorf("Plant export failed: %v", err)
	}

}
