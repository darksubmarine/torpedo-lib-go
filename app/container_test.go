package app

import (
	"fmt"
	"github.com/darksubmarine/torpedo-lib-go/conf"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// HttpServerProvider provides a http server
type HttpServerProvider struct {
	BaseProvider
	cfg conf.Map

	// instance to be provided
	mux *http.ServeMux `torpedo.di:"provide,name=HTTPSERVER"`
}

func NewHttpServerProvider(config conf.Map) *HttpServerProvider {
	return &HttpServerProvider{cfg: config}
}

func (p *HttpServerProvider) Provide(c IContainer) error {
	p.mux = http.NewServeMux()
	p.mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
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

// HelloProvider adds a /hello endpoint to the bound http server
type HelloProvider struct {
	BaseProvider

	mux *http.ServeMux `torpedo.di:"bind,name=HTTPSERVER"`
}

func NewHelloProvider() *HelloProvider { return new(HelloProvider) }

func (p *HelloProvider) Provide(c IContainer) error {
	p.mux.HandleFunc("/hello", p.sayHello)
	return nil
}

func (p *HelloProvider) sayHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, HTTP!\n")
}

func initApplicationContainer() *ApplicationContainer {
	// 1. App configuration
	config := conf.Map{}
	config.Add(3000, "port")

	// 2. Depdencies
	opts := ContainerOpts{Log: ContainerLogsOpts{W: os.Stdout, L: slog.LevelInfo}}

	application := NewApplicationContainer(opts)
	application.WithProvider(NewHttpServerProvider(config))
	application.WithProvider(NewHelloProvider())

	// 3. Run your application!
	go func() { application.Run() }()

	time.Sleep(5 * time.Second)

	return application
}

func TestApplicationContainer_Run(t *testing.T) {
	application := initApplicationContainer()

	httpServeMux, err := application.Invoke("HTTPSERVER")
	assert.Nil(t, err)

	request := httptest.NewRequest("GET", "/ping", nil)
	recorder := httptest.NewRecorder()

	httpServeMux.(*http.ServeMux).ServeHTTP(recorder, request)

	assert.EqualValues(t, 200, recorder.Result().StatusCode)
	assert.EqualValues(t, "pong!", recorder.Body.String())
}
