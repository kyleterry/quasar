package quasar

type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	UUID    string `json:"uuid"`
	Service ServiceConfig
}

type ServiceConfig struct {
	SenderBind   string `json:"sender_bind"`
	ReceiverBind string `json:"receiver_bind"`
}

func GetConfig() *Config {
	return &Config{
		Name:    "hello",
		Version: "1.0",
		Service: ServiceConfig{
			SenderBind:   "tcp://localhost:61124",
			ReceiverBind: "tcp://localhost:61123",
		},
	}
}
