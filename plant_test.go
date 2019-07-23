package gohome_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/savardiego/gohome"
)

func makeTestPlant(t *testing.T) *gohome.Plant {
	buf := bytes.NewBufferString("{ \"name\": \"home\", \"address\": \"192.168.28.35:20000\", \"num\": 1, \"ambients\": { \"kitchen\": { \"num\": 1, \"Lights\": { \"table\": 1, \"main\": 2 } }, \"living\": { \"num\": 2, \"Lights\": { \"sofa\": 1, \"tv\": 2 } }, \"camera\": { \"num\": 5, \"Lights\": { \"bed\": 8, \"main\": 6 } } } }")
	p, err := gohome.NewPlant(buf)
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
	plant, err := gohome.NewPlant(config)
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
		w, err := plant.WhereFromDesc(v)
		if k != w.Code || err != nil {
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
		w, err := plant.WhereFromDesc(v)
		if k != w.Code || err != nil {
			t.Errorf("Wrong where '%s' instead of '%s' (err: %v)", w, k, err)
		}
	}
}
func TestNewWrongWhere(t *testing.T) {
	plant := makeTestPlant(t)
	exp := []string{
		"",
		"livingsofa",
		"living.",
		"living.wrong",
		"kitchen.sofa.wrong",
	}
	for i, v := range exp {
		w, err := plant.WhereFromDesc(v)
		if w.Desc != "" || err == nil {
			t.Errorf("Wrong where for '%d': '%s' (err: %v)", i, w, err)
		}
	}
}

func TestDecodeWhere(t *testing.T) {
	config, err := os.Open("testdata/casa.json")
	if err != nil {
		t.Errorf("cannot open json file")
	}
	defer config.Close()
	plant, err := gohome.NewPlant(config)
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
		wh, err := plant.WhereFromCode(w)
		if wh.Desc != e {
			t.Errorf("Where not decoded correctly, exp:%s  decoded:%s", e, wh.Desc)
		}
		if err != nil {
			t.Logf("Where not decoded correctly, exp:%s  decoded:%s", e, wh.Desc)
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

func TestParseParams(t *testing.T) {
	plant := makeTestPlant(t)
	exp := map[string][]string{
		"*1*1*11##": []string{"LIGHT", "TURN_ON", "kitchen.table", "COMMAND"},
		"*1*1*12##": []string{"LIGHT", "TURN_ON", "kitchen.main", "COMMAND"},
		"*1*1*21##": []string{"LIGHT", "TURN_ON", "living.sofa", "COMMAND"},
		"*1*1*22##": []string{"LIGHT", "TURN_ON", "living.tv", "COMMAND"},
		"*1*1*1##":  []string{"LIGHT", "TURN_ON", "kitchen", "COMMAND"},
		"*1*1*2##":  []string{"LIGHT", "TURN_ON", "living", "COMMAND"},
		"*3*2##":    []string{"", "", "", "INVALID"},
		"*1**2##":   []string{"", "", "", "INVALID"},
		"*1*1*##":   []string{"", "", "", "INVALID"},
		"*#1*1##":   []string{"LIGHT", "", "kitchen", "REQUEST"},
		"":          []string{"", "", "", "INVALID"},
	}
	for m, ts := range exp {
		msg := plant.ParseFrame(m)
		if msg.Who.Desc != ts[0] || msg.What.Desc != ts[1] || msg.Where.Desc != ts[2] || msg.Kind != ts[3] {
			t.Errorf("decoded values for message '%s' are wrong: %s!=%s  %s!=%s %s!=%s %s!=%s", m, ot, ts[0], tt, ts[1], et, ts[2], k, ts[3])
		}
	}
}

func TestFormatToJSON(t *testing.T) {
	plant := makeTestPlant(t)
	exp := map[string]string{
		"*1*1*11##": "{\"who\": \"LIGHT\", \"what\": \"TURN_ON\", \"where\": \"kitchen.table\", \"kind\": \"0\"}",
		"*1*1*12##": "{\"who\": \"LIGHT\", \"what\": \"TURN_ON\", \"where\": \"kitchen.main\", \"kind\": \"0\"}",
		"*#1*1*##":  "{\"who\": \"LIGHT\", \"what\": \"\", \"where\": \"living\", \"kind\": \"1\"}",
		"*1*1*22##": "{\"who\": \"LIGHT\", \"what\": \"TURN_ON\", \"where\": \"living.tv\", \"kind\": \"0\"}",
		"*#1*12##":  "{\"who\": \"LIGHT\", \"what\": \"\", \"where\": \"kitchen.main\", \"kind\": \"1\"}",
		"*1*1*2##":  "{\"who\": \"LIGHT\", \"what\": \"TURN_ON\", \"where\": \"living\", \"kind\": \"0\"}",
		"*3*2##":    "{\"who\": \"\", \"what\": \"\", \"where\": \"\", \"kind\": \"-1\"}",
		"*1*2##":    "{\"who\": \"LIGHT\", \"what\": \"\", \"where\": \"\", \"kind\": \"-1\"}",
		"*1*1##":    "{\"who\": \"LIGHT\", \"what\": \"\", \"where\": \"\", \"kind\": \"-1\"}",
		"":          "{\"who\": \"\", \"what\": \"\", \"where\": \"\", \"kind\": \"-1\"}",
	}
	for m, ts := range exp {
		json := plant.FormatToJSON(gohome.ParseFrame(m))
		if json != ts {
			t.Errorf("decoded JSON for message '%s' is wrong: %s!=%s", m, json, ts)
		}
	}
}
