package gohome

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

const COMMAND = 0
const REQUEST = 1
const SPECIAL = 2
const DIMENSIONGET = 3
const DIMENSIONSET = 4
const INVALID = -1

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
	Kind    int
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
	case COMMAND:
		frame := fmt.Sprintf("*#%s*%s##", m.Who, m.Where)
		return frame
	case REQUEST:
		frame := fmt.Sprintf("*%s*%s*%s##", m.Who, m.What, m.Where)
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

func IsValid(msg string) (bool, int) {
	if len(msg) < 5 {
		return false, -1
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
	return false, -1
}

func ExplainKind(kind int) string {
	switch kind {
	case COMMAND:
		return "COMMAND"
	case REQUEST:
		return "REQUEST"
	case DIMENSIONGET:
		return "DIMENSIONGET"
	case DIMENSIONSET:
		return "DIMENSIONSET"
	case INVALID:
		return "INVALID"
	}
	return ""
}
