package dependency

import (
	"github.com/darksubmarine/torpedo-lib-go/app"
	"github.com/darksubmarine/torpedo-lib-go/conf"
	"log/slog"
	"os"
)

const LOGGER = "LOGGER"

type LoggerProvider struct {
	app.BaseProvider
	cfg conf.Map
}

func NewLoggerProvider(config conf.Map) *LoggerProvider {
	return &LoggerProvider{cfg: config}
}

func (p *LoggerProvider) Provide(c app.IContainer) (interface{}, error) {
	lvl := p.cfg.FetchStringOrElse("info", "level")
	sLvl := conf.SlogLevel(lvl)
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: sLvl})), nil
}
