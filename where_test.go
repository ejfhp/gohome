package gohome_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/savardiego/gohome"
)

func makeTestPlant(t *testing.T) *gohome.Plant {
	buf := bytes.NewBufferString("{ \"name\": \"home\", \"address\": \"192.168.0.35:20000\", \"num\": 1, \"ambients\": { \"kitchen\": { \"num\": 1, \"Lights\": { \"table\": 1, \"main\": 2 } }, \"living\": { \"num\": 2, \"Lights\": { \"sofa\": 1, \"tv\": 2 } } } }")
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
	gohome.LoadPlant(config)
	if gohome.ServerAddress() != "192.168.0.35:20000" {
		t.Errorf("Import plant configuratin has wrong address: '%s', len:%d", gohome.ServerAddress(), len(gohome.ServerAddress()))
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
		w, err := gohome.NewWhere(v)
		if k != string(w) || err != nil {
			t.Errorf("Wrong where %s instead of %s (err: %v)", w, k, err)
		}
	}
}

func TestNewWhere(t *testing.T) {
	makeTestPlant(t)
	exp := map[string]string{
		"11": "kitchen.table",
		"12": "kitchen.main",
		"21": "living.sofa",
		"22": "living.tv",
		"1":  "kitchen",
		"2":  "living",
	}
	for k, v := range exp {
		w, err := gohome.NewWhere(v)
		if k != string(w) || err != nil {
			t.Errorf("Wrong where %s instead of %s (err: %v)", w, k, err)
		}
	}
}

func TestExport(t *testing.T) {
	makeTestPlant(t)
	f, err := os.OpenFile("export.plant", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Errorf("Cannot open file: %v", err)
	}
	defer f.Close()
	err = gohome.ExportPlant(f)
	if err != nil {
		t.Errorf("Plant export failed: %v", err)
	}

}
