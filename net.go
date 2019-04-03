package gohome

import (
	"errors"
	"fmt"
	"log"
	"net"
)

var connErr = map[string]string{
	"RESFAIL": "FAILED TO RESOLVE ADDRESS",
	"NOCONN":  "FAILED TO CONNECT",
	"NOSEND":  "FAILED TO SEND TO SERVER",
	"NOREAD":  "FAILED TO READ SERVER ANSWER",
	"NOOPCMD": "FAILED TO SEND OPEN COMMAND SESSION",
	"EMPTY":   "EMPTY MESSAGE",
	"NAK":     "MESSAGE NOT ACEPTED",
}

//ConnectionError wraps socket communicatin errors
type ConnectionError struct {
	code   string
	server string
	cause  error
}

var ErrConnectionFailed = errors.New("Connection failed.")

func (ce ConnectionError) Error() string {
	return fmt.Sprintf("CONNECTION ERROR: %s, ADDRESS: %s, CAUSE: %v", connErr[ce.code], ce.server, ce.cause)
}

//HomeError wraps OWN errors
type HomeError struct {
	code     string
	incoming string
	answer   string
}

func (he HomeError) Error() string {
	return fmt.Sprintf("MY HOME ERROR: '%s', SENT MESSAGE: '%s', SERVER ANSWER: '%s'", connErr[he.code], he.incoming, he.answer)
}

// SystemMessages contains the OpenWebNet codes for various system messages
var SystemMessages = map[string]string{
	"ACK":                   "*#*1##",
	"NACK":                  "*#*0##",
	"OPEN_COMMAND_SESSION":  "*99*0##", // OpenWebNet command to ask for a command session
	"OPEN_EVENT_SESSION":    "*99*1##",
	"OPEN_SCENARIO_SESSION": "*99*9##",
}

type Cable struct {
	address string
	conn    *net.TCPConn
}

//NewCable return a Cable connected to the OpenWebNet server at the given address
func NewCable(address string) (*Cable, bool) {
	log.Printf("NewCable; address: %s", address)
	c := Cable{address: address}
	ok := c.acked()
	return &c, ok
}

func (c *Cable) connect() error {
	tcpAddress, err := net.ResolveTCPAddr("tcp4", c.address)
	if err != nil {
		logError(ConnectionError{"RESFAIL", c.address, err})
		return ErrConnectionFailed
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		logError(ConnectionError{"NOCONN", c.address, err})
		return ErrConnectionFailed
	}
	if c.acked() {
		c.conn = conn
	}
	return nil
}

func (c *Cable) disconnect() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return ErrConnectionFailed
		}
	}
	return nil
}

func (c *Cable) sendCommand(message string) bool {
	ok := c.send(SystemMessages["OPEN_COMMAND_SESSION"])
	if !ok {
		return false
	}
	ok = c.acked()
	if !ok {
		return false
	}
	ok = c.send(message)
	ok = c.acked()
	return ok
}

func (c *Cable) sendRequest(request string) []string {
	answers := make([]string, 0, 10)
	ok := c.send(SystemMessages["OPEN_COMMAND_SESSION"])
	if !ok {
		return answers
	}
	ok = c.acked()
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

func (c *Cable) send(message string) bool {
	log.Printf("cable.send() message:%s", message)
	c.connect()
	defer c.disconnect()
	ok := true
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		logError(ConnectionError{"NOSEND", c.address, err})
		ok = false
	}
	return ok
}

func (c *Cable) acked() bool {
	msg, ok := c.receive()
	if !ok {
		return false
	}
	return isAck(msg)
}

//returns answer, ok
func (c *Cable) receive() (string, bool) {
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
	cable *Cable
	plant *Plant
}

//NewHome creates a new Home connected through the given Cable
func NewHome(cable *Cable, plant *Plant) *Home {
	return &Home{cable: cable, plant: plant}
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
