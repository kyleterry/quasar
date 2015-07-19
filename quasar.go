package quasar

import (
	"log"
	"os"
	"os/signal"

	"github.com/garyburd/redigo/redis"
)

type Service struct {
	Name        string
	UUID        string
	Description string
	HelpText    string
	Config      Config

	defaultHandler FilterChain
	handlers       []FilterChain
	conn           *Connection
}

type Message struct {
	Target      string
	Command     string
	Mask        string
	Direct      bool
	Nick        string
	Host        string
	FullMessage string
	User        string
	FromChannel bool
	Connection  string
	Payload     string
	Meta        Meta
}

type Meta struct {
	Name    string
	Version string
}

type ParsedMatch struct {
}

type FilterChain struct {
	Filters    []string
	DirectOnly bool
	Handler    Handler
}

type Handler func(ParsedMatch, Message)

func (s *Service) Send(line string, message Message) error {
	return nil
}

func NewServiceFromConfig(config Config) *Service {
	return NewService(config.Name, config.UUID)
}

func NewService(name, uuid string) *Service {
	return &Service{
		Name: name,
		UUID: uuid,
	}
}

func (s *Service) AddChain(chain FilterChain) {
	s.handlers = append(s.handlers, chain)
}

func (s *Service) AddDefaultHandler(chain FilterChain) {
	s.defaultHandler = chain
}

func (s *Service) findMatch(msg redis.Message) {

}

func (s *Service) dispatch(msg []byte) {

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
	conn := NewConn(s.Config, s)
	s.conn = conn
	err := conn.Start()
	if err != nil {
		log.Fatal(err)
	}

	go s.deligator()

	sigChn := make(chan os.Signal, 1)
	signal.Notify(sigChn, os.Interrupt)
	for {
		select {
		case <-sigChn:
			conn.Shutdown()
		}
	}
	return nil
}
