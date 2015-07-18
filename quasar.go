package quasar

type Service struct {
	Name        string
	UUID        string
	Description string
	HelpText    string

	defaultHandler  FilterChain
	messagesSent    int
	messagesMatched int
}

type Payload struct {
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

type Handler func(ParsedMatch, Payload)

func (s *Service) Send(message string, payload Payload) error {
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

}

func (s *Service) AddDefaultHandler(chain FilterChain) {

}

// Runs forever or until a signal stops the program
func (s *Service) Run() error {
	return nil
}
