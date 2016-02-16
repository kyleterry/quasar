package quasar

type Meta struct {
	Name    string
	Version string
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
