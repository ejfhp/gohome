package gohome

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

//GENERAL is the Where that refers to the entire plant
const GENERAL Where = "0"

//ErrAmbientNotFound is returned when the desired where is not found in the conf file
var ErrAmbientNotFound = errors.New("ambient not found")

//ErrLightNotFound is returned when the desired loight is not found in the conf file
var ErrLightNotFound = errors.New("light not found")

type Ambient struct {
	Num    int            `json:"num"`
	Lights map[string]int `json:"lights"`
}

type Plant struct {
	Name     string             `json:"name"`
	Num      int                `json:"num"`
	Address  string             `json:"address"`
	Ambients map[string]Ambient `json:"ambients"`
}

//NewWhere returns a
func (p *Plant) NewWhere(where string) (Where, error) {
	var noWhere Where
	if where == "general" {
		where := Where("0")
		return where, nil
	}
	split := strings.Split(where, ".")
	if len(split) == 2 {
		amb, ok := p.Ambients[split[0]]
		if !ok {
			return noWhere, ErrAmbientNotFound
		}
		lig := amb.Lights[split[1]]
		if !ok {
			return noWhere, ErrLightNotFound
		}
		where := Where(fmt.Sprintf("%d%d", amb.Num, lig))
		return where, nil
	}
	if len(split) == 1 {
		amb, ok := p.Ambients[split[0]]
		if !ok {
			return noWhere, ErrAmbientNotFound
		}
		where := Where(fmt.Sprintf("%d", amb.Num))
		return where, nil
	}
	return noWhere, ErrLightNotFound
}

//ServerAddress returns the server address for the loaded configuration
func (p *Plant) ServerAddress() string {
	return p.Address
}

//ExportPlant the current plant configuration to the given file
func (p *Plant) ExportPlant(f io.Writer) error {
	encoder := json.NewEncoder(f)
	return encoder.Encode(p)
}

//LoadPlant load a plant configuration from a json file. Return a pointer to the Plant that will be used.
func LoadPlant(config io.Reader) (*Plant, error) {
	decoder := json.NewDecoder(config)
	plant := Plant{}
	err := decoder.Decode(&plant)
	if err != nil {
		return nil, err
	}
	return &plant, nil
}
