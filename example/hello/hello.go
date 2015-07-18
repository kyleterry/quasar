package main

import (
	"fmt"
	"log"

	"github.com/kyleterry/quasar"
)

const HelpText = `This is help text
You can add information on how to use the service so tenyks can tell people
how it should be used`

const Description = "Says hello and ping back to a user"

func main() {
	config := quasar.GetConfig()
	service := quasar.NewService(config.Name, config.UUID)
	service.HelpText = HelpText
	service.Description = Description
	service.AddChain(
		quasar.FilterChain{
			Filters:    []string{"^hello$"},
			DirectOnly: true,
			Handler: func(match quasar.ParsedMatch, payload quasar.Payload) {
				log.Print("Hello handler called")
				if err := service.Send(fmt.Sprintf("Hello, %s!", payload.Nick), payload); err != nil {
					log.Print(err)
				}
			},
		},
	)

	service.AddChain(
		quasar.FilterChain{
			Filters:    []string{"^ping$"},
			DirectOnly: true,
			Handler: func(match quasar.ParsedMatch, payload quasar.Payload) {
				log.Print("Ping handler called")
				service.Send(fmt.Sprintf("Pong, %s!", payload.Nick), payload)
			},
		},
	)

	service.AddDefaultHandler(
		quasar.FilterChain{
			Filters:    []string{},
			DirectOnly: true,
			Handler: func(match quasar.ParsedMatch, payload quasar.Payload) {
				log.Print("Default handler called")
				service.Send(fmt.Sprintf("you rang, %s?", payload.Nick), payload)
			},
		},
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
