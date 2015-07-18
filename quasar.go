package quasar

type Data struct {
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

type Handler func(ParsedMatch, Data)

func AddChain(fc FilterChain) {

}

func Send(message string, data Data) {

}
