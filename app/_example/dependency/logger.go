package dependency

import (
	"github.com/darksubmarine/torpedo-lib-go/app"
	"github.com/darksubmarine/torpedo-lib-go/app/_example/mypkg"
	"github.com/darksubmarine/torpedo-lib-go/conf"
	"log/slog"
	"os"
)

const LOGGER = "LOGGER"

type LoggerProvider struct {
	app.BaseProvider

	// singletonInstance is a pointer to a singleton instance of slog.Logger
	singletonInstance *slog.Logger `torpedo.di:"provide"`

	// builderInstance Instance builder function that returns a pointer to a new instance of slog.Logger
	builderInstance func(c app.IContainer) *slog.Logger `torpedo.di:"provide,name=LOGGER"`

	// private fields initialized by constructor
	cfg conf.Map
}

func NewLoggerProvider(config conf.Map) *LoggerProvider {
	return &LoggerProvider{cfg: config}
}

func (p *LoggerProvider) Provide(c app.IContainer) error {

	// Invoke the registered object
	greeter := c.InvokeP("GREETER").(*mypkg.Greeter)

	p.singletonInstance = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	p.builderInstance = p.builder

	// calling the invoked greeter
	p.singletonInstance.Info(greeter.Message())
	return nil
}

func (p *LoggerProvider) builder(c app.IContainer) *slog.Logger {
	lvl := p.cfg.FetchStringOrElse("info", "level")
	sLvl := conf.SlogLevel(lvl)
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: sLvl}))
}
