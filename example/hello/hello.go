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

func matchHello(msg quasar.Message) (quasar.Result, error) {
	res := make(quasar.Result)
	if msg.Payload != "hello" {
		return nil, quasar.ErrNoMatch
	}
	return res, nil
}

var matchName = quasar.NewRegexMatcher("^my name is (?P<name>(.*))$")

func main() {
	config := &quasar.Config{
		Name:    "hello",
		Version: "1.0",
		Service: quasar.ServiceConfig{
			SendAddr: "tcp://localhost:61124",
			RecvAddr: "tcp://localhost:61123",
		},
	}
	service := quasar.NewService(config)
	service.HelpText = HelpText
	service.Description = Description
	service.Handle(
		quasar.MsgHandler{
			MatcherFunc: quasar.MatcherFunc(matchHello),
			DirectOnly:  true,
			MatchHandler: quasar.HandlerFunc(func(match quasar.Result, msg quasar.Message) {
				log.Print("Hello handler called")
				if err := service.Send(fmt.Sprintf("Hello, %s!", msg.Nick), msg); err != nil {
					log.Print(err)
				}
			}),
		},
	)

	service.Handle(
		quasar.MsgHandler{

			MatcherFunc: matchName,
			DirectOnly:  true,
			MatchHandler: quasar.HandlerFunc(func(match quasar.Result, msg quasar.Message) {
				log.Print("Name handler called")
				if name, ok := match["name"]; ok {
					if err := service.Send(fmt.Sprintf("hello, %s!", name), msg); err != nil {
						log.Print(err)
					}
				}
			}),
		})

	service.DefaultHandle(
		quasar.MsgHandler{
			DirectOnly: true,
			MatchHandler: quasar.HandlerFunc(func(match quasar.Result, msg quasar.Message) {
				log.Print("Default handler called")
				service.Send(fmt.Sprintf("you rang, %s?", msg.Nick), msg)
			}),
		})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
