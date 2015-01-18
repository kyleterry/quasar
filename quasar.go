package quasar

type Service struct {
	chains map[string]FilterChain
}

type Data struct {
}

type Match struct {
}

type FilterChain struct {
	Patterns []string
	DirectOnly bool
	Handler Handler
}

type Handler interface {
	Run(service *Service, data Data, match Match)
}

func (s *Service) AddFilterChain(name string, chain FilterChain) {

} 

func (s *Service) Run() int {
	return 0
}
