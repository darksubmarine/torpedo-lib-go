package app

var nilFn = func() error { return nil }

type IProvider interface {
	Provide(lf IContainer) (interface{}, error)
	OnRegister() func() error
	OnStart() func() error
	OnStop() func() error
}

type BaseProvider struct{}

func (p *BaseProvider) OnRegister() func() error { return nilFn }
func (p *BaseProvider) OnStart() func() error    { return nilFn }
func (p *BaseProvider) OnStop() func() error     { return nilFn }
