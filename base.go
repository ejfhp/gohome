package gohome

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type Who struct {
	Code string
	Text string
}

func WhoFromText(text string) Who {
	return ListWho[strings.ToUpper(who)]
}

func WhoFromCode(code string) (Who, error) {
	if code == "" {
		return Who{}, nil
	}
	for k, v := range ListWho {
		if v.Code == code {
			return Who{code, k}, nil
		}
	}
	return Who{}, ErrWhoNotFound
}

type What struct {
	Code string
	Text string
}

func (w Who) WhatFromText(text string) What {
	return WhoWhat[w][strings.ToUpper(text)]
}

func (w Who) WhatFromCode(code string) (What, error) {
	if code == "" {
		return What{}, nil
	}
	for k, v := range WhoWhat[w] {
		if v.Code == code {
			return What{code, k}, nil
		}
	}
	return What{}, ErrWhatNotFound
}

type Dimension string
type Value string
type Message struct {
	Who     Who
	What    What
	Where   Where
	Kind    int
	special string
}

func (m Message) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"WHO\": \"")
	buffer.WriteString(m.Who.Text)
	buffer.WriteString("\", ")
	buffer.WriteString("\"WHAT\": \"")
	buffer.WriteString(m.What.Text)
	buffer.WriteString("\", ")
	buffer.WriteString("\"WHERE\": \"")
	buffer.WriteString(m.Where.Text)
	buffer.WriteString("\"")
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

const COMMAND = 0
const REQUEST = 1
const SPECIAL = 2
const DIMENSIONGET = 3
const DIMENSIONSET = 4
const INVALID = -1

var ErrWhatNotFound = errors.New("WHAT not found")
var ErrWhoNotFound = errors.New("WHO not found")

var Light = Who{"1", "LIGHT"}

var ListWho = map[string]Who{
	"LIGHT": Light,
}

var listLightningWhat = map[string]What{
	"TURN_OFF":         What("0"),
	"TURN_ON":          What("1"),
	"SET_20":           What("2"),
	"SET_30":           What("3"),
	"SET_40":           What("4"),
	"SET_50":           What("5"),
	"SET_60":           What("6"),
	"SET_70":           What("7"),
	"SET_80":           What("8"),
	"SET_90":           What("9"),
	"SET_100":          What("10"),
	"ON_1_MIN":         What("11"),
	"ON_2_MIN":         What("12"),
	"ON_3_MIN":         What("13"),
	"ON_4_MIN":         What("14"),
	"ON_5_MIN":         What("15"),
	"ON_15_MIN":        What("16"),
	"ON_30_SEC":        What("17"),
	"ON_0_5_SEC":       What("18"),
	"BLINK_ON_0_5_SEC": What("20"),
	"BLINK_ON_1_SEC":   What("21"),
	"BLINK_ON_1_5_SEC": What("22"),
	"BLINK_ON_2_SEC":   What("23"),
	"BLINK_ON_2_5_SEC": What("24"),
	"BLINK_ON_3_SEC":   What("25"),
	"BLINK_ON_3_5_SEC": What("26"),
	"BLINK_ON_4_SEC":   What("27"),
	"BLINK_ON_4_5_SEC": What("28"),
	"BLINK_ON_5_SEC":   What("29"),
	"UP_ONE_LEVEL":     What("30"),
	"DOWN_ONE_LEVEL":   What("31"),
	"JOLLY":            What("1000"),
}

//WhoWhat maps all commands for every kind of application (light, automation..)
var WhoWhat = map[Who]map[string]What{
	Light: listLightningWhat,
}
var regexpCommand = regexp.MustCompile(`^\*([0-9]{1,2})\*([0-9]{1,2})\*([0-9]{1,2})##`)
var regexpRequest = regexp.MustCompile(`^\*#([0-9]{1,2})\*([0-9]{1,2})##`)
var regexpDimensionGet = regexp.MustCompile(`^\*#([0-9]{1,2})\*([0-9]{1,2})\*([0-9]{1,2})##`)
var regexpDimensionSet = regexp.MustCompile(`^\*#([0-9]{1,2})\*([0-9]{1,2})\*#([0-9]{1,2})(\*[0-9]{1,2})+##`)

// var regexpRequest = regexp.MustCompile(`(^\*#[0-9])`)

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
func NewCommand(who Who, what What, where Where) Message {
	return Message{Who: who, What: what, Where: where, Kind: COMMAND}
}

func ParseFrame(frame string) Message {
	message := Message{}
	valid, msgkind := IsValid(frame)
	if !valid {
		fmt.Printf("Frame not valid: %s\n", frame)
		return Message{Kind: INVALID}
	}
	if msgkind == REQUEST {
		t := regexpRequest.FindStringSubmatch(string(frame))
		fmt.Printf("Parse reques (%s): %v\n", frame, t)
		message.Who = Who(t[1])
		message.Where = Where(t[2])
		message.Kind = REQUEST
		return message
	}
	if msgkind == COMMAND {
		t := regexpCommand.FindStringSubmatch(string(frame))
		fmt.Printf("Parse command (%s): %v\n", frame, t)
		message.Who = Who(t[1])
		message.What = What(t[2])
		message.Where = Where(t[3])
		message.Kind = COMMAND
		return message
	}
	if msgkind == DIMENSIONGET {
		t := regexpDimensionGet.FindStringSubmatch(string(frame))
		fmt.Printf("Parse dimget (%s): %v\n", frame, t)
		message.Who = Who(t[1])
		message.Where = Where(t[2])
		message.Kind = DIMENSIONGET
		return message
	}
	if msgkind == DIMENSIONSET {
		t := regexpDimensionSet.FindStringSubmatch(string(frame))
		fmt.Printf("Parse dimset (%s): %v\n", frame, t)
		message.Who = Who(t[1])
		message.Where = Where(t[2])
		message.Kind = DIMENSIONSET
		return message
	}
	return message
}

//NewRequest build a new Request to send to the home plant
func NewRequest(who Who, what What, where Where) Message {
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
