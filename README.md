# Quasar
Library for implementing Tenyks services in Go

## Install

`go get github.com/kyleterry/quasar`

## Usage

```go
package main

import (
	"fmt"
	"github.com/kyleterry/quasar"
	"flag"
	"os"
)

type OurHandler struct {}

func (h *OurHandler) Run(service *quasar.Service, data quasar.Data, match quasar.Match) {
	name := match["name"]
	service.Send(fmt.Sprintf("Hey %s, I'm Tenyks", name))
}

func main() {
	var config = flag.String("config", "", "Path to configuration file")
	flag.Parse()
	service := quasar.NewService()
	service.InitConfigFromFile(config)

	// Filters
	chain := FilterChain{
		Patterns: "^(?i)(hi|hello|sup|hey), I'm (?P<name>(.*))$",
		DirectOnly: true,
		Handler Ourhandler{}
	}
	service.AddFilterChain("say_hello", chain)

	os.Exit(service.Run())
}
```
