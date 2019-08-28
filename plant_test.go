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

func TestParseFrame(t *testing.T) {
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
		if !msg.IsValid() {
			if msg.Kind != ts[3] {
				t.Errorf("decoded/expected values for message '%s' are wrong: %s/%s", m, msg.Kind, ts[3])
			}
			continue
		}
		if msg.Who.Desc != ts[0] || msg.What.Desc != ts[1] || msg.Where.Desc != ts[2] || msg.Kind != ts[3] {
			t.Errorf("decoded/expected values for message '%s' are wrong: %s %s %s %s", m, msg.Who.Desc, msg.What.Desc, msg.Where.Desc, msg.Kind)
		}
	}
}

func TestWhereFromFrame(t *testing.T) {
	plant := makeTestPlant(t)
	messages := [][]string{
		{"*1*1*23##", "23"},
		{"*#1*1*#43*8##", "1"},
		{"*1*0*13##", "13"},
		{"*1*0*1##", "1"},
		{"*#1*21##", "21"},
		{"*#1*21*2##", "21"},
		{"*#1*1##", "1"},
		{"*1*2##", ""},
		{"*1**2##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := plant.ParseFrame(m[0])
		expWhereCode := m[1]
		if !msg.IsValid() {
			t.Logf("Frame '%s' is invalid", m)
			continue
		}
		if msg.Where.Code != expWhereCode {
			t.Errorf("%d - Wrong WHERE decoded: exp:%s actual:'%s'", i, expWhereCode, msg.Where.Code)
		}
	}
}

func TestWhoFromFrame(t *testing.T) {
	plant := makeTestPlant(t)
	messages := [][]string{
		{"*1*1*23##", "1"},
		{"*#1*1*#43*8##", "1"},
		{"*1*0*13##", "1"},
		{"*1*0*1##", "1"},
		{"*#1*21##", "1"},
		{"*#1*21*2##", "1"},
		{"*#1*1##", "1"},
		{"*1*2##", ""},
		{"*1**2##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := plant.ParseFrame(m[0])
		expWho := gohome.Who{Code: m[1]}
		if msg.IsValid() {
			t.Logf("Frame '%s' is invalid", m)
			break
		}
		if msg.Who.Code != expWho.Code {
			t.Errorf("%d - Wrong WHO decoded: exp:%s actual:%s", i, expWho, msg.Who)
		}
	}
}

func TestWhatFromFrame(t *testing.T) {
	plant := makeTestPlant(t)
	messages := [][]string{
		{"*1*1*23##", "1"},
		{"*#1*1*#43*8##", ""},
		{"*1*0*13##", "0"},
		{"*1*10*13##", "10"},
		{"*1*0*1##", "0"},
		{"*#1*21##", ""},
		{"*#1*21*2##", ""},
		{"*#1*1##", ""},
		{"*1*2##", ""},
		{"*1**2##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := plant.ParseFrame(m[0])
		expWhat := gohome.What{Code: m[1]}
		if msg.IsValid() {
			t.Logf("Frame '%s' is invalid", m)
			break
		}
		if msg.What.Code != expWhat.Code {
			t.Errorf("%d - Wrong WHAT decoded: exp:%s actual:%s", i, expWhat, msg.What)
		}
	}
}
func TestDecodeWhoFromFrame(t *testing.T) {
	plant := makeTestPlant(t)
	messages := [][]string{
		{"*1*1*23##", "LIGHT"},
		{"*1*0*13##", "LIGHT"},
		{"*1*11*1##", "LIGHT"},
		{"*1*18*21##", "LIGHT"},
		{"21##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := plant.ParseFrame(m[0])
		expWho := m[1]
		if msg.IsValid() {
			t.Logf("Frame '%s' is invalid", m)
			break
		}
		if msg.Who.Desc != expWho {
			t.Errorf("%d - Wrong WHO decoded: exp:%s actual:%s", i, expWho, msg.Who)
		}
	}
}

func TestDecodeWhatFromFrame(t *testing.T) {
	plant := makeTestPlant(t)
	messages := [][]string{
		{"*1*1*23##", "TURN_ON"},
		{"*1*0*13##", "TURN_OFF"},
		{"*1*11*1##", "ON_1_MIN"},
		{"*1*18*21##", "ON_0_5_SEC"},
		{"21##", ""},
		{"", ""},
	}
	for i, m := range messages {
		msg := plant.ParseFrame(m[0])
		expWhat := m[1]
		if msg.IsValid() {
			t.Logf("Frame '%s' is invalid", m)
			break
		}
		if msg.What.Desc != expWhat {
			t.Errorf("%d - Wrong WHERE decoded: exp:%s actual:%s", i, expWhat, msg.What.Desc)
		}
	}
}
func TestFormatToJSON(t *testing.T) {
	plant := makeTestPlant(t)
	exp := map[string]string{
		"*1*1*11##": "{\"who\":\"LIGHT\",\"what\":\"TURN_ON\",\"where\":\"kitchen.table\",\"kind\":\"COMMAND\"}",
		"*1*1*12##": "{\"who\":\"LIGHT\",\"what\":\"TURN_ON\",\"where\":\"kitchen.main\",\"kind\":\"COMMAND\"}",
		"*#1*1*##":  "{\"who\":\"\",\"what\":\"\",\"where\":\"\",\"kind\":\"INVALID\"}",
		"*1*1*22##": "{\"who\":\"LIGHT\",\"what\":\"TURN_ON\",\"where\":\"living.tv\",\"kind\":\"COMMAND\"}",
		"*#1*12##":  "{\"who\":\"LIGHT\",\"what\":\"\",\"where\":\"kitchen.main\",\"kind\":\"REQUEST\"}",
		"*1*1*2##":  "{\"who\":\"LIGHT\",\"what\":\"TURN_ON\",\"where\":\"living\",\"kind\":\"COMMAND\"}",
		"*3*2##":    "{\"who\":\"\",\"what\":\"\",\"where\":\"\",\"kind\":\"INVALID\"}",
		"*1*2##":    "{\"who\":\"\",\"what\":\"\",\"where\":\"\",\"kind\":\"INVALID\"}",
		"*1*1##":    "{\"who\":\"\",\"what\":\"\",\"where\":\"\",\"kind\":\"INVALID\"}",
		"":          "{\"who\":\"\",\"what\":\"\",\"where\":\"\",\"kind\":\"INVALID\"}",
	}
	for m, ts := range exp {
		json := plant.FormatToJSON(plant.ParseFrame(m))
		if json != ts {
			t.Errorf("decoded JSON for message '%s' is wrong: %s!=%s", m, json, ts)
		}
	}
}

func TestWhereFromCode(t *testing.T) {
	plant := makeTestPlant(t)
	wheres := map[string]string{
		"11": "kitchen.table",
		"12": "kitchen.main",
		"1":  "kitchen",
		"22": "living.tv",
	}
	for k, v := range wheres {
		ambient, err := plant.WhereFromCode(k)
		if err != nil {
			t.Errorf("failed to decod where (%s) due to %v", v, err)
		}
		if ambient.Desc != v {
			t.Errorf("decoded where (%s) is not the expected (%s)", ambient.Desc, v)
		}
	}
}

func TestWhereFromDesc(t *testing.T) {
	plant := makeTestPlant(t)
	wheres := map[string]string{
		"kitchen.table": "11",
		"kitchen.main":  "12",
		"kitchen":       "1",
		"living.tv":     "22",
	}
	for k, v := range wheres {
		ambient, err := plant.WhereFromDesc(k)
		if err != nil {
			t.Errorf("failed to decod where (%s) due to %v", v, err)
		}
		if ambient.Code != v {
			t.Errorf("decoded where (%s) is not the expected (%s)", ambient.Code, v)
		}
	}

}
