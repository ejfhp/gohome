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

func TestDecodeWhere(t *testing.T) {
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
		"":   "",
	}
	for w, e := range exp {
		wh := gohome.Where(w)
		dec, err := plant.DecodeWhere(wh)
		fmt.Printf("Where:%s decoded:%s\n", wh, dec)
		if dec != e {
			t.Errorf("Where not decoded correctly, exp:%s  decoded:%s", wh, dec)
		}
		if err != nil {
			t.Logf("Where not decoded correctly, exp:%s  decoded:%s", wh, dec)
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

func TestParse(t *testing.T) {
	plant := makeTestPlant(t)
	exp := map[string][]string{
		"*1*1*11##": []string{"LIGHT", "TURN_ON", "kitchen.table"},
		"*1*1*12##": []string{"LIGHT", "TURN_ON", "kitchen.main"},
		"*1*1*21##": []string{"LIGHT", "TURN_ON", "living.sofa"},
		"*1*1*22##": []string{"LIGHT", "TURN_ON", "living.tv"},
		"*1*1*1##":  []string{"LIGHT", "TURN_ON", "kitchen"},
		"*1*1*2##":  []string{"LIGHT", "TURN_ON", "living"},
		"*3*2##":    []string{"", "", ""},
		"*1**2##":   []string{"LIGHT", "", "living"},
		"*1*1*##":   []string{"LIGHT", "", ""},
		"":          []string{"", "", ""},
	}
	for m, ts := range exp {
		ot, tt, et, err := plant.Parse(gohome.Message(m))
		if err != nil {
			t.Errorf("failed to decode message '%s' due to: %v", m, err)
		}
		if string(ot) != ts[0] || string(tt) != ts[1] || string(et) != ts[2] {
			t.Errorf("decoded values fom message '%s' are wrong: %s!=%s  %s!=%s %s!=%s", m, ot, ts[0], tt, ts[1], et, ts[2])

		}
	}

}
