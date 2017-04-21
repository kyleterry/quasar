package main

import (
	"fmt"
	"log"

	"github.com/kyleterry/quasar"
)

const HelpText = `Use this to convey what the service does during a help query`

const Description = "Says hello and ping back to a user"

func matchHello(msg quasar.Message) quasar.Result {
	res := make(quasar.Result)
	if msg.Payload != "hello" {
		return nil
	}
	return res
}

var matchName = quasar.NewRegexMatcher("^my name is (?P<name>(.*))$")

func main() {
	// Start by creating a configuration object to pass into the service when it's created.
	config := &quasar.Config{
		Name:    "hello",
		Version: "1.0",
		Service: quasar.ServiceConfig{
			SendAddr: "tcp://localhost:61124",
			RecvAddr: "tcp://localhost:61123",
		},
	}

	service := quasar.New(config)
	service.HelpText = HelpText
	service.Description = Description

	// Handle tables a MsgHandler that will try to match incoming messages against anything you define.
	service.Handle(
		quasar.MsgHandler{
			MatcherFunc: quasar.MatcherFunc(matchHello),
			DirectOnly:  true,
			MatchHandler: quasar.HandlerFunc(func(match quasar.Result, msg quasar.Message, com quasar.Communication) {
				log.Print("Hello handler called")
				com.Send(fmt.Sprintf("Hello, %s!", msg.Nick), msg)
			}),
			// This is a handler level help text. You can use it to document a handler for !help queries.
			HelpText: "hello - tenyks will respond with a hello back",
		},
	)

	service.Handle(
		quasar.MsgHandler{
			MatcherFunc: matchName,
			DirectOnly:  true,
			MatchHandler: quasar.HandlerFunc(func(match quasar.Result, msg quasar.Message, com quasar.Communication) {
				log.Print("Name handler called")
				if name, ok := match["name"]; ok {
					com.Send(fmt.Sprintf("nice to meet you, %s!", name), msg)
				}
			}),
			HelpText: "my name is <name> - tenyks will respond with a personal greeting back",
		},
	)

	// You can set a default handler. This is handler is called if nothing matches an incoming message.
	service.DefaultHandle(
		quasar.MsgHandler{
			DirectOnly: true,
			MatchHandler: quasar.PrivmsgMiddlware(quasar.HandlerFunc(func(match quasar.Result, msg quasar.Message, com quasar.Communication) {
				fmt.Println(msg)
				log.Print("Default handler called")
				com.Send(fmt.Sprintf("you rang, %s?", msg.Nick), msg)
			})),
			HelpText: "if nothing matches, tenyks will respond to any message asking if you rang",
		},
	)

	log.Print("Starting hello service")
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
