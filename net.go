package gohome

import (
	"log"
	"net"
	"time"

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
var SystemMessages = map[string]Message{
	"ACK":                   Message{special: "*#*1##"},
	"NACK":                  Message{special: "*#*0##"},
	"QUERY_ALL":             Message{special: "*#5##"},
	"OPEN_COMMAND_SESSION":  Message{special: "*99*0##"}, // OpenWebNet command to ask for a command session
	"OPEN_EVENT_SESSION":    Message{special: "*99*1##"},
	"OPEN_SCENARIO_SESSION": Message{special: "*99*9##"},
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
	cable := newCable(address)
	return &Home{cable: cable, plant: plant}
}

//Do some action with your home
func (h *Home) Do(command Message) error {
	log.Printf("Home.Do")
	if command.Kind != COMMAND {
		return errors.Errorf("Message is not a command: %v", command)
	}
	return h.cable.sendCommand(command)
}

//Ask the system
func (h *Home) Ask(request Message) ([]Message, error) {
	log.Printf("Home.Ask")
	if request.Kind != REQUEST {
		return nil, errors.Errorf("Message is not a request: %v", request)
	}
	frames, err := h.cable.sendRequest(request)
	if err != nil {
		return []Message{}, errors.Wrapf(err, "cannot send request frame '%v'", request)
	}
	//TODO parse di tutte le frame
	res := make([]Message, len(frames))
	for i, f := range frames {
		res[i] = h.plant.ParseFrame(f)
	}
	return res, nil
}

func (h *Home) Listen() (<-chan string, chan<- struct{}, <-chan error) {
	msgChan := make(chan string, 1)
	signChan := make(chan struct{})
	errChan := make(chan error)
	go h.cable.listen(msgChan, signChan, errChan)
	return msgChan, signChan, errChan
}

func newCable(address string) *Cable {
	c := Cable{address: address}
	return &c
}

func (c *Cable) connect() (*net.TCPConn, error) {
	log.Printf("Cable.connect address:%s", c.address)
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", c.address, timeout)
	if conn == nil {
		return nil, errors.Wrapf(ErrNoConnection, "no dial to server address: %s", c.address)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "no dial to server address: %s", c.address)
	}
	connTCP := conn.(*net.TCPConn)
	if err != nil {
		return nil, errors.Wrapf(ErrNoConnection, "no dial to server address: %s", c.address)
	}
	if !c.acked(connTCP) {
		return nil, ErrNAK
	}
	return connTCP, nil
}

func (c *Cable) sendCommand(command Message) error {
	log.Printf("Cable.SendCommmand message:%v", command)
	conn, err := c.connect()
	if err != nil {
		return errors.Wrap(err, "cannot connect")
	}
	defer conn.Close()
	if c.send(conn, SystemMessages["OPEN_COMMAND_SESSION"].Frame()) != nil {
		return errors.Wrap(err, "cannot open session")
	}
	if !c.acked(conn) {
		return ErrNAK
	}
	if c.send(conn, command.Frame()) != nil {
		return errors.Wrapf(err, "cannot send message %v, ", command)
	}
	if !c.acked(conn) {
		return ErrNAK
	}
	return nil
}

func (c *Cable) sendRequest(request Message) ([]string, error) {
	log.Printf("Cable.SendRequest request:%v", request)
	conn, err := c.connect()
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect")
	}
	defer conn.Close()
	if c.send(conn, SystemMessages["OPEN_COMMAND_SESSION"].Frame()) != nil {
		return nil, errors.Wrapf(err, "cannot open command session")
	}
	if !c.acked(conn) {
		return nil, ErrNAK
	}
	if c.send(conn, request.Frame()) != nil {
		return nil, errors.Wrapf(err, "cannot send request %v, ", request)
	}
	answers := make([]string, 0, 10)
	for {
		a, err := c.receive(conn, false)
		if err != nil {
			return answers, errors.Wrapf(err, "failed to receive answer for request: %v", request)
		}
		if a == SystemMessages["ACK"].Frame() {
			break
		}
		answers = append(answers, a)
	}
	return answers, nil
}

func (c *Cable) listen(out chan<- string, in <-chan struct{}, errs chan<- error) {
	log.Printf("Cable.listen")
	conn, err := c.connect()
	if err != nil {
		errs <- errors.Wrap(err, "cannot connect")
		close(out)
		return
	}
	defer conn.Close()
	if c.send(conn, SystemMessages["OPEN_EVENT_SESSION"].Frame()) != nil {
		errs <- errors.Wrapf(err, "cannot open event session")
		close(out)
		return
	}
	if !c.acked(conn) {
		errs <- ErrNAK
		close(out)
		return
	}
Listen:
	for {
		select {
		case <-in:
			close(out)
			break Listen
		default:
			frame, err := c.receive(conn, true)
			if err != nil {
				errs <- errors.Wrapf(err, "failed to receive events")
			}
			if ok, _ := IsValid(frame); ok {
				out <- frame
			}
		}
	}
}

func (c *Cable) send(conn *net.TCPConn, frame string) error {
	log.Printf("Cable.send frame:%s", frame)
	_, err := conn.Write([]byte(frame))
	if err != nil {
		return errors.Wrap(ErrConnectionFailed, "failed to send")
	}
	return nil
}

func (c *Cable) acked(conn *net.TCPConn) bool {
	msg, err := c.receive(conn, false)
	if err != nil {
		log.Printf("Cannot check ACK: %+v", err)
		return false
	}
	if msg != SystemMessages["ACK"].Frame() {
		return false
	}
	return true
}

//returns answer, ok
func (c *Cable) receive(conn *net.TCPConn, noTimeout bool) (string, error) {
	if conn == nil {
		return "", errors.Wrap(ErrNoConnection, "cannot receive from nil connection")
	}
	frame := make([]byte, 0, 20)
	b := make([]byte, 1)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 1))
		n, err := conn.Read(b)
		if err, ok := err.(net.Error); ok && err.Timeout() && noTimeout {
			log.Printf("...timeout")
			break
		}
		if err != nil {
			return "", errors.Wrap(err, "cannot read from connection")
		}
		if n == 0 {
			return "", ErrNoData
		}
		frame = append(frame, b[0])
		//TODO sostituire con regexp
		if len(frame) > 1 && frame[len(frame)-1] == '#' && frame[len(frame)-2] == '#' {
			break
		}
	}
	log.Printf("Cable.receive '%s'", frame)
	return string(frame), nil
}
