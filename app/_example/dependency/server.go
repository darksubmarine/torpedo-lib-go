package dependency

import (
	"github.com/darksubmarine/torpedo-lib-go/app"
	"github.com/darksubmarine/torpedo-lib-go/conf"
	"log/slog"

	"fmt"
	"io"
	"net/http"
)

type HttpServerProvider struct {
	app.BaseProvider
	cfg conf.Map

	// instance to be provided
	mux *http.ServeMux `torpedo.di:"provide,name=HTTPSERVER"`

	// binds
	logger *slog.Logger `torpedo.di:"bind,name=LOGGER"`
}

func NewHttpServerProvider(config conf.Map) *HttpServerProvider {
	return &HttpServerProvider{cfg: config}
}

func (p *HttpServerProvider) Provide(c app.IContainer) error {
	p.mux = http.NewServeMux()
	p.mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		p.logger.Info("/ping has been called")
		io.WriteString(writer, "pong!")
	})

	return nil
}

func (p *HttpServerProvider) OnStart() func() error {
	return func() error {
		go func() {
			port := p.cfg.FetchIntP("port")
			if err := http.ListenAndServe(fmt.Sprintf(":%d", port), p.mux); err != nil {
				panic(fmt.Sprintf("error starting HTTP server at %d with error %s", port, err))
			}
		}()

		return nil
	}
}
