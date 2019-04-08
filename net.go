package gohome

import (
	"log"
	"net"

	"github.com/pkg/errors"
)

var ErrNoConnection = errors.New("NO CONNECTION")
var ErrNoData = errors.New("NO DATA")
var ErrServerNotFound = errors.New("SERVER NOT FOUND")
var ErrConnectionFailed = errors.New("CONNECTION FAILED")

//HomeError wraps OWN errors
type HomeError struct {
	code     string
	incoming string
	answer   string
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
}

//Home is a Btcino MyHome plant that can be controlled with a OpenWebNet enabled device (F452 ecc)
type Home struct {
	cable *Cable
	plant *Plant
}

//NewHome creates a new Home connected through the given Cable
func NewHome(plant *Plant) *Home {
	log.Printf("NewHome")
	address := plant.ServerAddress()
	cable := NewCable(address)
	return &Home{cable: cable, plant: plant}
}

//Do some action with your home
func (h *Home) Do(command Command) bool {
	log.Printf("Home.Do")
	return h.cable.SendCommand(string(command))
}

//Ask the system
func (h *Home) Ask(request Command) []string {
	log.Printf("Home.Ask")
	return h.cable.SendRequest(string(request))
}

//NewCable return a Cable connected to the OpenWebNet server at the given address
func NewCable(address string) *Cable {
	c := Cable{address: address}
	return &c
}

func (c *Cable) connect() (*net.TCPConn, error) {
	log.Printf("Cable.connect address:%s", c.address)
	tcpAddress, err := net.ResolveTCPAddr("tcp4", c.address)
	if err != nil {
		return nil, errors.Wrapf(ErrConnectionFailed, "server address: %d", c.address)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		return nil, errors.Wrapf(ErrNoConnection, "no dial to server address: %d", c.address)
	}
	if c.acked(conn) {
		return nil, errors.Wrapf(ErrConnectionFailed, "NAK from server address: %d", c.address)
	}
	return conn, nil
}

func (c *Cable) SendCommand(message string) error {
	log.Printf("Cable.SendCommmand")
	conn, err := c.connect()
	if err != nil {
		return errors.Wrap(err, "cannot connect")
	}
	defer conn.Close()
	ans, err := c.send(conn, SystemMessages["OPEN_COMMAND_SESSION"])
	if err != nil {
		return errors.Wrap(err, "cannot open session")
	}
	if c.acked(conn) {
		return errors.Wrapf(ErrConnectionFailed, "NAK from server address: %d", c.address)
	}
	ans, err = c.send(conn, message)
	if err != nil {
		return errors.Wrapf(err, "cannot send message %s, ", message)
	}
	if c.acked(conn) {
		return errors.Wrapf(ErrConnectionFailed, "NAK from server address: %d", c.address)
	}
	return nil
}

func (c *Cable) SendRequest(request string) []string {
	c.connect()
	defer c.disconnect()
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

func (c *Cable) send(conn *net.TCPConn, message string) (string, error) {
	if c.conn == nil {
		log.Printf("Connection is nil\n")
		return "", false
	}
	log.Printf("cable.send() message:%s", message)
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		return "", ConnectionError{"NOSEND", c.address, err}
	}
}

func (c *Cable) acked(conn *net.TCPConn) bool {
	msg, ok := c.receive()
	if !ok {
		return false
	}
	return isAck(msg)
}

//returns answer, ok
func (c *Cable) receive(conn *net.TCPConn) (string, error) {
	if conn == nil {
		return "", errors.Wrap(ErrNoConnection, "cannot receive from nil connection")
	}
	msg := make([]byte, 0, 20)
	b := make([]byte, 1)
	var end bool
	ok := true
	for !end {
		n, err := conn.Read(b)
		if err != nil {
			return "", ErrNoData
		}
		if n == 0 {
			return "", ErrNoData
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

func logError(err error) {
	log.Printf("ERROR: %v", err)
}
