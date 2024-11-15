package main

import (
	"github.com/darksubmarine/torpedo-lib-go/app"
	"github.com/darksubmarine/torpedo-lib-go/app/_example/dependency"
	"github.com/darksubmarine/torpedo-lib-go/app/_example/mypkg"
	"github.com/darksubmarine/torpedo-lib-go/conf"
	"log/slog"
	"os"
)

func main() {
	// 1. App configuration
	config := conf.Load(false, conf.NewYamlLoader("app/_example/config.yml"))

	// 2. Dependencies
	opts := app.ContainerOpts{Log: app.ContainerLogsOpts{W: os.Stdout, L: slog.LevelInfo}}

	application := app.NewApplicationContainer(opts)

	application.WithProvider(
		dependency.NewLoggerProvider(config.FetchSubMapP("log")))

	application.WithProvider(
		dependency.NewHttpServerProvider(config.FetchSubMapP("server")))

	application.WithProvider(
		dependency.NewHelloProvider())

	// a custom object created and registered into the application container.
	greeter := mypkg.NewGreeter()
	if err := application.Register("GREETER", greeter); err != nil {
		panic("something happened registering custom object")
	}

	// 3. Run your application!
	application.Run()

}
