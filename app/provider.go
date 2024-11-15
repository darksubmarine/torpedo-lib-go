package app

var nilFn = func() error { return nil }

// IProvider dependency provider interface that cover the app life-cycle events
type IProvider interface {
	Provide(c IContainer) error
	OnStart() func() error
	OnStop() func() error
}

// BaseProvider no operation provider to create your own provider
type BaseProvider struct{}

// OnStart returns a hook function to be executed at start step
func (p *BaseProvider) OnStart() func() error { return nilFn }

// OnStop returns a hook function to be executed at stop step
func (p *BaseProvider) OnStop() func() error { return nilFn }
