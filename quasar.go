package quasar

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	Name        string
	UUID        string
	Description string
	HelpText    string
	Config      *Config

	defaultHandler MsgHandler
	handlers       []MsgHandler
	conn           *Connection
}

// Matcher is a type that has a match method that can receive a Message and try to match it
// against a set of rules.
// If no match is found, then you should return nil, ErrNoMatch
type Matcher interface {
	Match(Message) (Result, error)
}

type MatcherFunc func(Message) (Result, error)

func (fn MatcherFunc) Match(msg Message) (Result, error) {
	return fn(msg)
}

// Handler is a type that has a Handle method that takes a Result and Message as arguments.
// What you do with these is up to you. This is a function that will be triggered if it's
// matcher returns a valid match.
type Handler interface {
	HandleMatch(Result, Message)
}

type HandlerFunc func(Result, Message)

func (fn HandlerFunc) HandleMatch(r Result, m Message) {
	fn(r, m)
}

type Result map[string]string

type MsgHandler struct {
	MatcherFunc  Matcher
	DirectOnly   bool
	PrivateOnly  bool
	MatchHandler Handler
}

func (s *Service) Send(line string, message Message) error {
	message.Payload = line
	j, err := json.Marshal(message)
	if err != nil {
		return err
	}
	s.publish(string(j))
	return nil
}

func (s *Service) publish(msg string) {
	s.conn.out <- msg
}

var NoopHandler = MsgHandler{MatchHandler: HandlerFunc(func(r Result, m Message) {})}

func NewService(config *Config) *Service {
	s := &Service{
		Config: config,
	}
	// Add noop handler for default
	s.DefaultHandle(NoopHandler)

	return s
}

func (s *Service) Handle(handler MsgHandler) {
	// Choose one, sucker
	if handler.DirectOnly && handler.PrivateOnly {
		log.Panicln("Cannot have both DirectOnly and PrivateOnly set to true")
	}
	s.handlers = append(s.handlers, handler)
}

func (s *Service) DefaultHandle(handler MsgHandler) {
	s.defaultHandler = handler
}

func (s *Service) findMatch(msg Message) {
	for _, mh := range s.handlers {
		res, err := mh.MatcherFunc.Match(msg)
		if IsNoMatch(err) {
			continue
		} else {
			mh.MatchHandler.HandleMatch(res, msg)
			return
		}
	}
	s.defaultHandler.MatchHandler.HandleMatch(nil, msg)
}

func (s *Service) dispatch(rawmsg string) {
	msg := Message{}
	if err := json.Unmarshal([]byte(rawmsg), &msg); err != nil {
		// ?
	}
	s.findMatch(msg)
}

func (s *Service) deligator() {
	for {
		select {
		case msg := <-s.conn.in:
			go s.dispatch(msg)
		}
	}
}

// Runs forever or until a signal stops the program
func (s *Service) Run() error {
	sigch := make(chan os.Signal, 1)
	go func() {
		<-sigch
		s.Cleanup()
	}()

	conn, err := NewConn(s.Config, s)
	if err != nil {
		panic(err)
	}
	s.conn = conn
	err = conn.start()
	if err != nil {
		panic(err)
	}

	go s.deligator()

	sigchn := make(chan os.Signal, 1)
	signal.Notify(sigchn, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		select {
		case <-sigchn:
			conn.close()
		}
	}
	return nil
}

func (s *Service) Cleanup() {
	s.conn.close()
}
