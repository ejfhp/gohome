package gohome

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type Who string
type What string
type Dimension string
type Value string
type Where string
type Message string

var ErrWhatNotFound = errors.New("WHAT not found")
var ErrWhoNotFound = errors.New("WHO not found")

var Light = Who("1")

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
var regexpWhere = regexp.MustCompile(`(\*[1])(.*\*)([0-9]{1,2})(##)`)
var regexpWho = regexp.MustCompile(`(^\*)([0-9]{1,2})(.*)`)
var regexpWhat = regexp.MustCompile(`(^\*)([0-9]{1,2})(\*)([0-9]{1,2})(\*)([0-9]{1,2})(##)`)

func NewWho(who string) Who {
	return ListWho[who]
}

func DecodeWho(who Who) (string, error) {
	for k, v := range ListWho {
		if v == who {
			return k, nil
		}
	}
	return "", ErrWhoNotFound
}

func (w Who) Text() string {
	return string(w)
}

func (w Who) NewWhat(what string) What {
	return WhoWhat[w][strings.ToUpper(what)]
}

func (w Who) DecodeWhat(what What) (string, error) {
	for k, v := range WhoWhat[w] {
		if v == what {
			return k, nil
		}
	}
	return "", ErrWhatNotFound
}

func (w What) Text() string {
	return string(w)
}

//NewCommand build a new Command to send to the home plant
func NewCommand(who Who, what What, where Where) Message {
	cmd := fmt.Sprintf("*%s*%s*%s##", who, what, where)
	return Message(cmd)
}

func (m Message) Where() Where {
	w := regexpWhere.FindStringSubmatch(string(m))
	if len(w) < 2 {
		return Where("")
	}
	return Where(w[len(w)-2])
}

func (m Message) Who() Who {
	w := regexpWho.FindStringSubmatch(string(m))
	if len(w) < 3 {
		return Who("")
	}
	return Who(w[2])
}

func (m Message) What() What {
	w := regexpWhat.FindStringSubmatch(string(m))
	if len(w) < 7 {
		return What("")
	}
	return What(w[4])
}

func (m Message) Decode() (Who, What, Where) {
	return m.Who(), m.What(), m.Where()
}
