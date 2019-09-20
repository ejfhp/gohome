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
	Desc string
}

//GENERAL is the Where that refers to the entire plant
var GENERAL Where = Where{Code: "0", Desc: "GENERAL"}

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

//NewPlant load a plant configuration from a json file. Return a pointer to the Plant that will be used.
func NewPlant(config io.Reader) (*Plant, error) {
	if config == nil {
		return nil, errors.New("Plant configuration is nil")
	}
	decoder := json.NewDecoder(config)
	plant := Plant{}
	err := decoder.Decode(&plant)
	if err != nil {
		return nil, err
	}
	return &plant, nil
}

//NewWhere returns a
func (p *Plant) WhereFromDesc(text string) (Where, error) {
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

//Decode returns where defined by the ambient and light names in the plant config file: <ambient>[.<light>]
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

//ParseFrame parse a OWN frame and returns a structured message.
func (p *Plant) ParseFrame(frame string) Message {
	fmt.Printf("Checking frame: %s\n", frame)
	message := Message{}
	valid, msgkind := IsValid(frame)
	if !valid {
		fmt.Printf("Frame not valid: %s\n", frame)
		message.Kind = INVALID
		return message
	}
	if msgkind == REQUEST {
		t := regexpRequest.FindStringSubmatch(string(frame))
		fmt.Printf("Frame (%s) recognized as REQUEST: %v\n", frame, t)
		message.Who = NewWho(t[1])
		where, err := p.WhereFromCode(t[2])
		if err != nil {
			log.Printf("Frame not valid: %s due to: %v\n", frame, err)
			message.Kind = INVALID
			return message
		}
		message.Where = where
		message.Kind = REQUEST
		return message
	}
	if msgkind == COMMAND {
		t := regexpCommand.FindStringSubmatch(string(frame))
		log.Printf("Frame (%s) recognized as COMMAND: %v\n", frame, t)
		message.Who = NewWho(t[1])
		what, err := message.Who.WhatFromCode(t[2])
		if err != nil {
			log.Printf("Frame what not valid: %s due to: %v\n", frame, err)
			message.Kind = INVALID
			return message
		}
		message.What = what
		where, err := p.WhereFromCode(t[3])
		if err != nil {
			log.Printf("Frame where not valid: %s due to: %v\n", frame, err)
			message.Kind = INVALID
			return message
		}
		message.Where = where
		message.Kind = COMMAND
		return message
	}
	if msgkind == DIMENSIONGET {
		t := regexpDimensionGet.FindStringSubmatch(string(frame))
		fmt.Printf("Frame (%s) recognized as DIMENSIONGET: %v\n", frame, t)
		message.Who = NewWho(t[1])
		where, err := p.WhereFromCode(t[2])
		if err != nil {
			log.Printf("Frame where not valid: %s due to: %v\n", frame, err)
			message.Kind = INVALID
			return message
		}
		message.Where = where
		message.Kind = DIMENSIONGET
		return message
	}
	if msgkind == DIMENSIONSET {
		t := regexpDimensionSet.FindStringSubmatch(string(frame))
		fmt.Printf("Frame (%s) recognized as DIMENSIONSET: %v\n", frame, t)
		message.Who = NewWho(t[1])
		where, err := p.WhereFromCode(t[2])
		if err != nil {
			log.Printf("Frame where not valid: %s due to: %v\n", frame, err)
			message.Kind = INVALID
			return message
		}
		message.Where = where
		message.Kind = DIMENSIONSET
		return message
	}
	return message
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

func (p *Plant) ParseFromJSON(jsonMessage string) Message {
	var mapMsg map[string]interface{}
	msg := Message{}
	json.Unmarshal([]byte(jsonMessage), &mapMsg)
	who, ok := mapMsg["who"].(string)
	if !ok {
		fmt.Printf("who not found\n")
		return Message{Kind: INVALID}
	}
	msg.Who = NewWho(who)
	var err error
	var wa What
	var we Where
	what, ok := mapMsg["what"].(string)
	if ok {
		wa, err = msg.Who.WhatFromDesc(what)
		msg.What = wa
	}
	where, ok := mapMsg["where"].(string)
	if !ok {
		fmt.Printf("where not found\n")
		return Message{Kind: INVALID}
	}
	we, err = p.WhereFromDesc(where)
	msg.Where = we
	kind := mapMsg["kind"].(string)
	msg.Kind = kind
	if err != nil {
		fmt.Printf("cannot parse message from JSON due to %v\n", err)
		return Message{Kind: INVALID}
	}
	fmt.Println(mapMsg)
	return msg
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
