package dependency

import (
	"github.com/darksubmarine/torpedo-lib-go/app"
	"io"
	"log/slog"

	"net/http"
)

type HelloProvider struct {
	app.BaseProvider

	logger *slog.Logger   `torpedo.di:"bind,name=LOGGER"`
	mux    *http.ServeMux `torpedo.di:"bind,name=HTTPSERVER"`
}

func NewHelloProvider() *HelloProvider { return new(HelloProvider) }

func (p *HelloProvider) Provide(c app.IContainer) error {
	p.mux.HandleFunc("/hello", p.sayHello)

	// in this case we won't be providing anything...
	// only using the Provide method to inject a hello controller
	// as part of the mux server
	return nil
}

func (p *HelloProvider) sayHello(w http.ResponseWriter, r *http.Request) {
	p.logger.Info("got /hello request")
	io.WriteString(w, "Hello, HTTP!\n")
}
