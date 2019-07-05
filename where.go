package gohome

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/pkg/errors"

)

type Where struct {
	Code string
	Text string
}

//GENERAL is the Where that refers to the entire plant
var GENERAL Where = Where{"0", "GENERAL"}

//ErrAmbientNotFound is returned when the desired where is not found in the conf file
var ErrAmbientNotFound = errors.New("ambient not found")

//ErrLightNotFound is returned when the desired loight is not found in the conf file
var ErrLightNotFound = errors.New("light not found")

//ErrWhereNotInPlant is returned whene a where numeric code is not found in the current plant configuration
var ErrWhereNotInPlant = errors.New("WHERE not found in the current plant configuration")

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
func (p *Plant) WhereFromText(text string) (Where, error) {
	var noWhere Where
	if strings.ToUpper(text) == "GENERAL" {
		where := GENERAL
		return where, nil
	}
	split := strings.Split(text, ".")
	if len(split) == 2 {
		amb, ok := p.Ambients[split[0]]
		if !ok {
			return noWhere, ErrAmbientNotFound
		}
		lig, ok := amb.Lights[split[1]]
		if !ok {
			return noWhere, ErrLightNotFound
		}
		where := Where{fmt.Sprintf("%d%d", amb.Num, lig), text}
		return where, nil
	}
	if len(split) == 1 {
		amb, ok := p.Ambients[split[0]]
		if !ok {
			return noWhere, ErrAmbientNotFound
		}
		where := Where{fmt.Sprintf("%d", amb.Num), text}
		return where, nil
	}
	return noWhere, ErrLightNotFound
}

//Decode returns where defined by the ambient an light names in the plant config file: <ambient>[.<light>]
func (p *Plant) WhereFromCode(code string) (Where, error) {
	if code == "" {
		return Where{}, nil
	}
	var wtext string
	if len(code) < 1 || len(code) > 2 {
		return Where{}, ErrWhereNotInPlant
	}
	amb, err := strconv.Atoi(string(code[0:1]))
	if err != nil {
		return Where{}, errors.Wrapf(ErrWhereNotInPlant, "where: %v", code)
	}
	for ka, a := range p.Ambients {
		if a.Num == amb {
			wtext = ka
			if len(code) == 2 {
				lig, err := strconv.Atoi(string(code[1:2]))
				if err != nil {
					return Where{}, errors.Wrapf(ErrWhereNotInPlant, "where: %v", code)
				}
				for kl, pl := range a.Lights {
					if pl == lig {
						wtext = wtext + "." + kl
					}
				}
			}
		}
	}
	return Where{code, wtext}, nil
}

//Parse return the who, what, where of a message as three different params, and if is a request
func (p *Plant) Explain(msg Message) (string, string, string, string) {
	var decwho, decwhat, decwhere string
	var err error
	decwho, err = DecodeWho(msg.Who)
	if err != nil {
		log.Printf("Plant.Parse - cannot decode WHO of message '%v' due to: %v", msg, err)
		return decwho, decwhat, decwhere, "INVALID"
	}
	decwhere, err = p.DecodeWhere(msg.Where)
	if err != nil {
		log.Printf("Plant.Parse - cannot decode WHERE of message '%v' due to: %v", msg, err)
		return decwho, decwhat, decwhere, "INVALID"
	}
	if msg.Kind == COMMAND {
		decwhat, err = msg.Who.DecodeWhat(msg.What)
		if err != nil {
			log.Printf("Plant.Parse - cannot decode WHAT of message '%v' due to: %v", msg, err)
			return decwho, decwhat, decwhere, "INVALID"
		}
	}
	return decwho, decwhat, decwhere, ExplainKind(msg.Kind)
}

//FormatToJSON returns the who, what, where of a message in a JSON formatted string
func (p *Plant) FormatToJSON(msg Message) string {
	j, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Plant.FormatToJSON - cannot format message '%v' due to: %v", msg, err)
		return "{ERROR: }"
	}
	return string(j)
}

func (p *Plant) ParseFromJSON(json string) Message {
	return Message{}
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
