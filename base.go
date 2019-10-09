package gohome

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

const COMMAND = "COMMAND"
const REQUEST = "REQUEST"
const SPECIAL = "SPECIAL"
const DIMENSIONGET = "DIMENSIONGET"
const DIMENSIONSET = "DIMENSIONSET"
const INVALID = "INVALID"

var ErrWhatNotFound = errors.New("WHAT not found")
var ErrWhoNotFound = errors.New("WHO not found")

var regexpCommand = regexp.MustCompile(`^\*([0-9]{1,2})\*([0-9]{1,2})\*([0-9]{1,2})##`)
var regexpRequest = regexp.MustCompile(`^\*#([0-9]{1,2})\*([0-9]{1,2})##`)
var regexpDimensionGet = regexp.MustCompile(`^\*#([0-9]{1,2})\*([0-9]{1,2})\*([0-9]{1,2})##`)
var regexpDimensionSet = regexp.MustCompile(`^\*#([0-9]{1,2})\*([0-9]{1,2})\*#([0-9]{1,2})(\*[0-9]{1,2})+##`)

type Dimension string
type Value string

type What struct {
	Code string
	Desc string
}

type Message struct {
	Who     *Who   `json:"who"`
	What    What   `json:"what"`
	Where   Where  `json:"where"`
	Kind    string `json:"kind"`
	special string
}

func (m Message) MarshalJSON() ([]byte, error) {
	var whoD, whatD, whereD string
	if m.Who != nil {
		whoD = m.Who.Desc
	}
	if m.What != (What{}) {
		whatD = m.What.Desc
	}
	if m.Where != (Where{}) {
		whereD = m.Where.Desc
	}
	mj := struct {
		Who   string `json:"who"`
		What  string `json:"what"`
		Where string `json:"where"`
		Kind  string `json:"kind"`
	}{
		Who:   whoD,
		What:  whatD,
		Where: whereD,
		Kind:  m.Kind,
	}
	js, err := json.Marshal(&mj)
	return js, err
}

func (m Message) Frame() string {
	if m.IsSpecial() {
		return m.special
	}
	switch m.Kind {
	case REQUEST:
		frame := fmt.Sprintf("*#%s*%s##", m.Who.Code, m.Where.Code)
		return frame
	case COMMAND:
		frame := fmt.Sprintf("*%s*%s*%s##", m.Who.Code, m.What.Code, m.Where.Code)
		return frame
	}
	return ""
}

func (m Message) IsSpecial() bool {
	if m.special != "" {
		return true
	}
	return false
}

func (m Message) IsValid() bool {
	if m.Kind != INVALID {
		return true
	}
	return false
}

//NewCommand build a new Command to send to the home plant
func NewCommand(who *Who, what What, where Where) Message {
	return Message{Who: who, What: what, Where: where, Kind: COMMAND}
}

//NewRequest build a new Request to send to the home plant
func NewRequest(who *Who, what What, where Where) Message {
	return Message{Who: who, What: what, Where: where, Kind: REQUEST}
}

func IsValid(msg string) (bool, string) {
	if len(msg) < 5 {
		return false, INVALID
	}
	for _, m := range SystemMessages {
		if msg == m.Frame() {
			return true, SPECIAL
		}
	}
	switch {
	case regexpCommand.MatchString(msg):
		return true, COMMAND
	case regexpRequest.MatchString(msg):
		return true, REQUEST
	case regexpDimensionGet.MatchString(msg):
		return true, DIMENSIONGET
	case regexpDimensionSet.MatchString(msg):
		return true, DIMENSIONSET
	}
	return false, INVALID
}
