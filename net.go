package gohome

import (
	"log"
	"net"

	"github.com/pkg/errors"
)

var ErrNoConnection = errors.New("NO CONNECTION")
var ErrNoData = errors.New("NO DATA")
var ErrNAK = errors.New("NAK")
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
func (h *Home) Do(command Command) error {
	log.Printf("Home.Do")
	return h.cable.SendCommand(string(command))
}

//Ask the system
func (h *Home) Ask(request Command) ([]string, error) {
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
		return nil, errors.Wrapf(ErrConnectionFailed, "server address: %s", c.address)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		return nil, errors.Wrapf(ErrNoConnection, "no dial to server address: %s", c.address)
	}
	if err = c.acked(conn); err != nil {
		return nil, errors.Wrapf(ErrConnectionFailed, "NAK from server address: %s", c.address)
	}
	return conn, nil
}

func (c *Cable) SendCommand(message string) error {
	log.Printf("Cable.SendCommmand message:%s", message)
	conn, err := c.connect()
	if err != nil {
		return errors.Wrap(err, "cannot connect")
	}
	defer conn.Close()
	log.Printf("Cable.open_command_session")
	if c.send(conn, SystemMessages["OPEN_COMMAND_SESSION"]) != nil {
		return errors.Wrap(err, "cannot open session")
	}
	log.Printf("Cable.check ack")
	if err = c.acked(conn); err != nil {
		return errors.Wrapf(ErrConnectionFailed, "NAK from server address: %s", c.address)
	}
	log.Printf("Sending nessage")
	if c.send(conn, message) != nil {
		return errors.Wrapf(err, "cannot send message %s, ", message)
	}
	log.Printf("Cable.check ack")
	if err = c.acked(conn); err != nil {
		return errors.Wrapf(ErrConnectionFailed, "NAK from server address: %s", c.address)
	}
	return nil
}

func (c *Cable) SendRequest(request string) ([]string, error) {
	log.Printf("Cable.SendRequest request:%s", request)
	conn, err := c.connect()
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect")
	}
	defer conn.Close()
	if c.send(conn, SystemMessages["OPEN_COMMAND_SESSION"]) != nil {
		return nil, errors.Wrapf(err, "cannot open session")
	}
	if err = c.acked(conn); err != nil {
		return nil, errors.Wrapf(ErrConnectionFailed, "NAK from server address: %s", c.address)
	}
	if c.send(conn, request) != nil {
		return nil, errors.Wrapf(err, "cannot send request %s, ", request)
	}
	answers := make([]string, 0, 10)
	for {
		a, err := c.receive(conn)
		if err != nil {
			return answers, errors.Wrapf(err, "failed to receive answer for request: %s", request)
		}
		if a == SystemMessages["ACK"] {
			break
		}
		answers = append(answers, a)
	}
	return answers, nil
}

func (c *Cable) send(conn *net.TCPConn, message string) error {
	log.Printf("Cable.send message:%s", message)
	_, err := conn.Write([]byte(message))
	if err != nil {
		return errors.Wrap(ErrConnectionFailed, "failed to send")
	}
	return nil
}

func (c *Cable) acked(conn *net.TCPConn) error {
	log.Printf("Cable.acked")
	msg, err := c.receive(conn)
	if err != nil {
		return errors.Wrap(err, "cannot check ACK")
	}
	return isAck(msg)
}

//returns answer, ok
func (c *Cable) receive(conn *net.TCPConn) (string, error) {
	log.Printf("Cable.receive")
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
	log.Printf("Cable.receive: msg:%s, ok:%t", ans, ok)
	return ans, nil
}

func isAck(m string) error {
	log.Printf("isAck")
	if m != SystemMessages["ACK"] {
		return ErrNAK
	}
	return nil
}
