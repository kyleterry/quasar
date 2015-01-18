package quasar

type Service struct {

}

type Data struct {

}

type Match struct {

}

type Handler interface {
	Run(data Data, match Match)
}

func (s *Service) Filter(filterKey string, pattern string, handlerFn Handler) {

} 

func (s *Service) Run() {

}
