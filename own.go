package gohome

import (
	"fmt"
	"log"
	"net"
)

var errors = map[string]string{
	"RESFAIL": "FAILED TO RESOLVE ADDRESS",
	"NOCONN":  "FAILED TO CONNECT",
	"NOSEND":  "FAILED TO SEND TO SERVER",
	"NOREAD":  "FAILED TO READ SERVER ANSWER",
	"NOOPCMD": "FAILED TO SEND OPEN COMMAND SESSION",
	"EMPTY":   "EMPTY MESSAGE",
	"NAK":     "MESSAGE NOT ACEPTED",
}

type What string
type Who string
type Dimension string
type Value string
type Where string
type Command string

//ConnectionError wraps socket communicatin errors
type ConnectionError struct {
	code   string
	server string
	cause  error
}

func (ce ConnectionError) Error() string {
	return fmt.Sprintf("CONNECTION ERROR: %s, ADDRESS: %s, CAUSE: %v", errors[ce.code], ce.server, ce.cause)
}

//HomeError wraps OWN errors
type HomeError struct {
	code     string
	incoming string
	answer   string
}

func (he HomeError) Error() string {
	return fmt.Sprintf("MY HOME ERROR: '%s', SENT MESSAGE: '%s', SERVER ANSWER: '%s'", errors[he.code], he.incoming, he.answer)
}

// SystemMessages contains the OpenWebNet codes for various system messages
var SystemMessages = map[string]string{
	"ACK":                   "*#*1##",
	"NACK":                  "*#*0##",
	"OPEN_COMMAND_SESSION":  "*99*0##", // OpenWebNet command to ask for a command session
	"OPEN_EVENT_SESSION":    "*99*1##",
	"OPEN_SCENARIO_SESSION": "*99*9##",
}

//Cable connects to the plant
type Cable interface {
	sendCommand(message string) bool
	sendRequest(request string) []string
}

type cable struct {
	address string
	conn    *net.TCPConn
}

//NewCable return a Cable connected to the OpenWebNet server at the given address
func NewCable(address string) (Cable, bool) {
	log.Printf("NewCable; address: %s", address)
	tcpAddress, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		logError(ConnectionError{"RESFAIL", address, err})
		return nil, false
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		logError(ConnectionError{"NOCONN", address, err})
		return nil, false
	}
	c := cable{address, conn}
	ok := c.ack()
	return &c, ok
}

func (c *cable) sendCommand(message string) bool {
	ok := c.send(SystemMessages["OPEN_COMMAND_SESSION"])
	if !ok {
		return false
	}
	ok = c.ack()
	if !ok {
		return false
	}
	ok = c.send(message)
	ok = c.ack()
	return ok
}

func (c *cable) sendRequest(request string) []string {
	answers := make([]string, 0, 10)
	ok := c.send(SystemMessages["OPEN_COMMAND_SESSION"])
	if !ok {
		return answers
	}
	ok = c.ack()
	if !ok {
		return answers
	}
	ok = c.send(request)
	var end bool
	for ok && !end {
		var a string
		a, ok = c.receive()
		end = isAck(a)
		if !end {
			answers = append(answers, a)
		}
	}
	return answers
}

func (c *cable) send(message string) bool {
	log.Printf("cable.send() message:%s", message)
	ok := true
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		logError(ConnectionError{"NOSEND", c.address, err})
		ok = false
	}
	return ok
}

func (c *cable) ack() bool {
	msg, ok := c.receive()
	if !ok {
		return false
	}
	return isAck(msg)
}

//returns answer, ok
func (c *cable) receive() (string, bool) {
	msg := make([]byte, 0, 20)
	b := make([]byte, 1)
	var end bool
	ok := true
	for !end {
		n, err := c.conn.Read(b)
		if err != nil {
			ok = false
			logError(ConnectionError{"NOREAD", "", nil})
		}
		if n == 0 {
			ok = false
			logError(ConnectionError{"EMPTY", "", nil})
		}
		msg = append(msg, b[0])
		if len(msg) > 1 && msg[len(msg)-1] == '#' && msg[len(msg)-2] == '#' {
			end = true
		}
	}
	ans := string(msg)
	log.Printf("cable.receive: msg:%s, ok:%t", ans, ok)
	return ans, ok
}

func isAck(m string) bool {
	var a bool
	if m == SystemMessages["ACK"] {
		a = true
	} else if m == SystemMessages["NACK"] {
		a = false
	}
	return a
}

//Home is a Btcino MyHome plant that can be controlled with a OpenWebNet enabled device (F452 ecc)
type Home struct {
	cable Cable
}

//NewHome creates a new Home connected through the given Cable
func NewHome(cable Cable) *Home {
	return &Home{cable: cable}
}

//Do some action with your home
func (h *Home) Do(command Command) bool {
	return h.cable.sendCommand(string(command))
}

//Ask the system
func (h *Home) Ask(request Command) []string {
	log.Printf("Home.Ask")
	return h.cable.sendRequest(string(request))
}

func logError(err error) {
	log.Printf("ERROR: %v", err)
}
