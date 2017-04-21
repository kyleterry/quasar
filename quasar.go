package quasar

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
)

const (
	// DefaultSend is the default tcp address for sending data to the server.
	DefaultSend = "tcp://localhost:61124"
	// DefaultRecv is the default tcp address for recieving data from the server.
	DefaultRecv = "tcp://localhost:61123"
)

// Service is the state of the service.
// It holds the handlers and matchers and helps relay
// messages to the right handler when they come in from
// the Tenyks server.
type Service struct {
	Name        string
	UUID        string
	Description string
	HelpText    string
	Config      *Config

	defaultHandler MsgHandler
	handlers       []MsgHandler
	conn           *Connection
	responseCh     chan Message
	ctrlCh         chan struct{}
}

// Matcher is an interface that has a match method that receives a
// Message and tries to match it against a set of rules.
// If no match is found, then you should return nil.
type Matcher interface {
	Match(Message) Result
}

// MatcherFunc registers the wrapped function as a Matcher.
type MatcherFunc func(Message) Result

func (fn MatcherFunc) Match(msg Message) Result {
	return fn(msg)
}

// Handler is an interface that has a Handle method that takes a Result and Message as arguments.
// What you do with these is up to you. This is a function that will be triggered if it's
// matcher returns a valid match.
type Handler interface {
	HandleMatch(Result, Message, Communication)
}

type HandlerFunc func(Result, Message, Communication)

func (fn HandlerFunc) HandleMatch(r Result, m Message, c Communication) {
	fn(r, m, c)
}

type Result map[string]string

type MsgHandler struct {
	MatcherFunc  Matcher
	DirectOnly   bool
	PrivateOnly  bool
	MatchHandler Handler
}

type Communication struct {
	ch chan<- Message
}

func (c *Communication) Send(line string, message Message) {
	message.Payload = line
	c.ch <- message
}

func (s *Service) responseReceiver() {
	for {
		message := <-s.responseCh
		j, err := json.Marshal(message)
		if err != nil {
			continue
		}
		s.publish(string(j))
	}
}

func (s *Service) publish(msg string) {
	s.conn.out <- msg
}

var NoopHandler = MsgHandler{MatchHandler: HandlerFunc(func(r Result, m Message, c Communication) {})}

func New(config *Config) *Service {
	s := &Service{
		Config: config,

		responseCh: make(chan Message, 1000),
		ctrlCh:     make(chan struct{}),
	}
	// Add noop handler for default
	s.DefaultHandle(NoopHandler)

	return s
}

func (s *Service) Handle(handler MsgHandler) {
	// Choose one, sucker
	if handler.DirectOnly && handler.PrivateOnly {
		// TODO: handle this with an error
		log.Panicln("Cannot have both DirectOnly and PrivateOnly set to true")
	}
	s.handlers = append(s.handlers, handler)
}

func (s *Service) DefaultHandle(handler MsgHandler) {
	s.defaultHandler = handler
}

func (s *Service) registerCommunicationAndCallHandler(handler Handler, result Result, msg Message) {
	com := Communication{s.responseCh}
	go handler.HandleMatch(result, msg, com)
}

func (s *Service) findMsgMatch(msg Message) {
	for _, mh := range s.handlers {
		res := mh.MatcherFunc.Match(msg)
		if res == nil {
			continue
		} else {
			s.registerCommunicationAndCallHandler(mh.MatchHandler, res, msg)
			return
		}
	}
	s.registerCommunicationAndCallHandler(s.defaultHandler.MatchHandler, nil, msg)
}

func (s *Service) deserializeAndDispatch(rawmsg string) {
	msg := Message{}
	if err := json.Unmarshal([]byte(rawmsg), &msg); err != nil {
		// We don't care about malformed messages, so we discard and return.
		return
	}
	s.findMsgMatch(msg)
}

func (s *Service) deligator() {
	for {
		select {
		case msg := <-s.conn.in:
			s.deserializeAndDispatch(msg)
		case <-s.ctrlCh:
			return
		}
	}
}

// Run runs forever or until a signal stops the program
func (s *Service) Run() error {
	if s.Config.Service.SendAddr == "" {
		s.Config.Service.SendAddr = DefaultSend
	}

	if s.Config.Service.RecvAddr == "" {
		s.Config.Service.RecvAddr = DefaultRecv
	}

	conn, err := NewConnection(s.Config)
	if err != nil {
		return errors.Wrap(err, "failed to create connection")
	}
	s.conn = conn
	err = conn.start()
	if err != nil {
		return errors.Wrap(err, "failed to start connection")
	}

	go s.responseReceiver()
	go s.deligator()

	sigchn := make(chan os.Signal, 1)
	signal.Notify(sigchn, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

LOOP:
	for {
		select {
		case <-sigchn:
			break LOOP
		}
	}

	s.Cleanup()
	return nil
}

// Cleanup will close open connections and clean up anything that shouldn't linger after shutdown.
func (s *Service) Cleanup() {
	close(s.ctrlCh)
	s.conn.close()
}
