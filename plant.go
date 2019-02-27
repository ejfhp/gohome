package gohome

import (
	"encoding/json"
	"io"
	"os"
)

type Ambient struct {
	Num    int            `json:"num"`
	Lights map[string]int `json:lights`
}

type Plant struct {
	Name     string             `json:"name"`
	Num      int                `json:"num"`
	Address  string             `json:"address"`
	Ambients map[string]Ambient `json:"ambients"`
}

func NewPlant(config io.Reader) *Plant {
	decoder := json.NewDecoder(config)
	plant := Plant{}
	decoder.Decode(&plant)
	return &plant
}

func (p *Plant) Where(ambient string, light string) {
}

//Export the current plant configuration to the given file
func (p *Plant) Export(f *os.File) error {
	encoder := json.NewEncoder(f)
	return encoder.Encode(p)
}
