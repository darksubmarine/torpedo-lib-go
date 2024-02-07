package main

import (
	"github.com/darksubmarine/torpedo-lib-go/app"
	"github.com/darksubmarine/torpedo-lib-go/app/_example/dependency"
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

	application.WithNamedProvider(
		dependency.LOGGER,
		dependency.NewLoggerProvider(config.FetchSubMapP("log")))

	application.WithNamedProvider(
		dependency.HTTP_SERVER,
		dependency.NewHttpServerProvider(config.FetchSubMapP("server")))

	application.WithProvider(dependency.NewHelloProvider())

	// 3. Run your application!
	application.Run()

}
