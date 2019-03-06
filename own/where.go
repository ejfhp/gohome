package own

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const GENERAL Where = "0"

func newWhere(ambient, light int) Where {
	s := "0"
	if ambient > 0 && light < 1 {
		s = fmt.Sprintf("%d", ambient)
	}
	if ambient > 0 && light > 0 {
		s = fmt.Sprintf("%d%d", ambient, light)
	}
	return Where(s)
}

func newWhere0(ambient, light int) Where {
	s := fmt.Sprintf("%02d%02d", ambient, light)
	return Where(s)
}

//Ambient define  aroom, an ambient of the home
type Ambient struct {
	Num    int            `json:"num"`
	Lights map[string]int `json:"lights"`
}

//Plant define the whole plant of the house
type Plant struct {
	Name     string             `json:"name"`
	Num      int                `json:"num"`
	Address  string             `json:"address"`
	Ambients map[string]Ambient `json:"ambients"`
}

//NewPlant build a new plant from a json file
func NewPlant(config io.Reader) *Plant {
	decoder := json.NewDecoder(config)
	plant := Plant{}
	decoder.Decode(&plant)
	return &plant
}

//AddressOfLight return the Where for the given light in the given ambient
func (p *Plant) AddressOfLight(ambient string, light string) Where {
	amb := p.Ambients[ambient]
	lig := amb.Lights[light]
	w := newWhere(amb.Num, lig)
	return w
}

//AddressOfAmb return the Where for the given ambient
func (p *Plant) AddressOfAmb(ambient string) Where {
	amb := p.Ambients[ambient]
	w := newWhere(amb.Num, 0)
	return w
}

//Export the current plant configuration to the given file
func (p *Plant) Export(f *os.File) error {
	encoder := json.NewEncoder(f)
	return encoder.Encode(p)
}
