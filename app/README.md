# Application Container

The `ApplicationContainer` handles the app life cycle letting you provide dependencies applying 
the control of inversion principle. 

The dependencies have access to the application life cycle at:

 - **Register**: this is the first app status, before to start, when you are providing your dependencies.
 - **Start**: Once that all dependencies are registered, the app is executed by `.Run()` method, here is when the app is set to `start` status and all depencies `OnStart` hooks are executed.
 - **Stop**: This is the last app status, when a signal termination is received the dependencies `OnStop` hooks are executed.

Please check the [example](_example) application.

## Provider

In order to inject a dependency the `IProvider` interface must be implemented. For instance a good example could be a
HTTP Server:

```go
package dependency

import (
	"github.com/darksubmarine/torpedo-lib-go/app"
	"github.com/darksubmarine/torpedo-lib-go/conf"
	
	"fmt"
	"net/http"
	"io"
)

const HTTP_SERVER = "HTTPSERVER"

type HttpServerProvider struct {
	app.BaseProvider
	
	cfg conf.Map
	mux  *http.ServeMux
}

func NewHttpServerProvider(config conf.Map) *HttpServerProvider {
	return &HttpServerProvider{cfg: config}
}

func (p *HttpServerProvider) Provide(c app.IContainer) (interface{}, error)  {

	p.mux = http.NewServeMux()

	p.mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		io.WriteString(writer, "pong!")
	})
	
	return p.mux, nil
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
```

Once that we have the HTTP server provider we could create new provider to add an endpoint:

```go
package dependency

import (
	"github.com/darksubmarine/torpedo-lib-go/app"
	"github.com/darksubmarine/torpedo-lib-go/conf"

	"fmt"
	"net/http"
	"io"
)

type HelloProvider struct {
	app.BaseProvider
}

func NewHelloProvider() *HelloProvider { return new(HelloProvider) }

func (p *HelloProvider) Provide(c app.IContainer) (interface{}, error)  {
    mux := c.InvokeP(HTTP_SERVER).(*http.ServeMux)
	mux.HandleFunc("/hello", sayHello)

	// in this case we won't be providing anything... 
	// only using the Provide method to inject a hello controller 
	// as part of the mux server
	return nil, nil
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

```

## Running the application container

The main application must be initialized with an instance of `IApp` like:

```go
package main

import (
	"github.com/darksubmarine/torpedo-lib-go/app"
	"github.com/darksubmarine/torpedo-lib-go/conf"
	"log/slog"
	"os"
)

func main() {

	// 1. App configuration
	config := conf.Map{}
	config.Add(3000, "port")

	// 2. Depdencies
	opts := app.ContainerOpts{Log: app.ContainerLogsOpts{W: os.Stdout, L: slog.LevelInfo}}

	application := app.NewApplicationContainer(opts)
	application.WithNamedProvider(dependency.HTTP_SERVER, dependency.NewHttpServerProvider(config))
	application.WithProvider(dependency.NewHelloProvider())
	
	// 3. Run your application!
	application.Run()

}

```