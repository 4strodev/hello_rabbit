package main

import (
	"log"

	"github.com/4strodev/rabbit_go_publisher/pkg/core/adapters"
	"github.com/4strodev/rabbit_go_publisher/pkg/core/app"
	"github.com/4strodev/rabbit_go_publisher/pkg/shared/infrastructure"
	wiring "github.com/4strodev/wiring/pkg"
)

func main() {
	var err error
	container := wiring.New()

	app := app.NewApp(container)
	err = app.AddComponent(&infrastructure.RabbitComponent{})
	if err != nil {
		panic(err)
	}

	err = app.AddAdapter(&adapters.CliAdapter{})
	if err != nil {
		panic(err)
	}

	err = app.Start()
	if err != nil {
		log.Fatal(err)
	}
}
