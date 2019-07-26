package gohome

import (
	"bytes"
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
	Who     *Who
	What    What
	Where   Where
	Kind    string
	special string
}

func (m Message) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"WHO\": \"")
	buffer.WriteString(m.Who.Desc)
	buffer.WriteString("\", ")
	buffer.WriteString("\"WHAT\": \"")
	buffer.WriteString(m.What.Desc)
	buffer.WriteString("\", ")
	buffer.WriteString("\"WHERE\": \"")
	buffer.WriteString(m.Where.Desc)
	buffer.WriteString("\"")
	buffer.WriteString("}")
	return buffer.Bytes(), nil
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
