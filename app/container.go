package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

type iAppLifeCycle interface {
	register(name string, obj interface{}, fn func() error) error
	onStart(func() error)
	onStop(func() error)
	start() []error
	stop() []error
}

type IContainerMonitor interface {
	Logger() *slog.Logger
}

type IContainer interface {
	IContainerMonitor
	Register(name string, obj interface{}) error
	Invoke(name string) (interface{}, error)
	InvokeP(name string) interface{}
	InvokeByType(obj interface{}) (interface{}, error)
	InvokeByTypeP(obj interface{}) interface{}
}

type IApp interface {
	iAppLifeCycle
	IContainer
	Run()
}

type ApplicationContainer struct {
	// Monitoring
	log *slog.Logger

	// signal channel
	sigs chan os.Signal

	// dynamic dependencies
	deps map[string]interface{}

	// hooks
	onStartHook []func() error
	onStopHook  []func() error
}

func NewContainer(opts ContainerOpts) *ApplicationContainer {

	container := &ApplicationContainer{
		log: slog.New(slog.NewTextHandler(opts.Log.W, &slog.HandlerOptions{
			Level: opts.Log.L,
		})),
		sigs:        make(chan os.Signal, 1),
		onStartHook: make([]func() error, 0),
		onStopHook:  make([]func() error, 0),
		deps:        make(map[string]interface{}),
	}

	signal.Notify(container.sigs, syscall.SIGINT, syscall.SIGTERM)

	return container
}

func (c *ApplicationContainer) print(format string, a ...any) {
	fmt.Println("[TPDO]", fmt.Sprintf(format, a...))
}

func (c *ApplicationContainer) printnln(format string, a ...any) {
	fmt.Println("\n[TPDO]", fmt.Sprintf(format, a...))
}

func (c *ApplicationContainer) execProvider(provider IProvider) interface{} {
	if obj, err := provider.Provide(c); err != nil {
		panic(err)
	} else {
		return obj
	}
	return nil
}

func (c *ApplicationContainer) register(name string, obj interface{}, fn func() error) error {
	if err := c.Register(name, obj); err != nil {
		c.print("ERROR registering provider %s", name)
		return err
	}

	if err := fn(); err != nil {
		c.print("ERROR on OnRegister provider hook function for %s", name)
		return err
	}

	return nil
}

func (c *ApplicationContainer) WithProvider(provider IProvider) *ApplicationContainer {
	return c.WithNamedProvider("", provider)
}

func (c *ApplicationContainer) WithNamedProvider(name string, provider IProvider) *ApplicationContainer {
	if obj := c.execProvider(provider); obj != nil {
		if name == "" {
			name = fmt.Sprint(reflect.TypeOf(obj))
		}
		c.register(name, obj, provider.OnRegister())

		// adding LifeCycle hooks
		c.onStart(provider.OnStart())
		c.onStop(provider.OnStop())
	}
	return c
}

func (c *ApplicationContainer) exitWithErrors(msg string, errs []error) {
	fmt.Println(msg)
	for i, err := range errs {
		fmt.Printf("  %d - %s\n", i, err)
	}
	os.Exit(1)
}

func (c *ApplicationContainer) Run() {
	defer func() {
		if errs := c.stop(); len(errs) > 0 {
			c.exitWithErrors("Some container dependencies cannot be stopped due to some errors", errs)
		}
	}()

	if startErrs := c.start(); len(startErrs) > 0 {
		c.exitWithErrors("Dependencies container cannot START due to some deps starting errors", startErrs)
	}

	c.print("Application started successfully... (waiting for signal termination)")
	<-c.sigs
	c.printnln("Application terminated by signal")
}

func (c *ApplicationContainer) onStart(fn func() error) {
	c.onStartHook = append(c.onStartHook, fn)
}

func (c *ApplicationContainer) onStop(fn func() error) {
	c.onStopHook = append(c.onStopHook, fn)
}

func (c *ApplicationContainer) Register(name string, obj interface{}) error {
	if obj == nil {
		return fmt.Errorf("the provided dependency as %s cannot be nil", name)
	}

	if _, exists := c.deps[name]; exists {
		return fmt.Errorf("depedency with name %s already provided", name)
	}

	c.print("Registering dependency for %s", name)
	c.deps[name] = obj
	return nil
}

func (c *ApplicationContainer) Invoke(name string) (interface{}, error) {
	if obj, exists := c.deps[name]; !exists {
		return nil, fmt.Errorf("depedency with name %s has not been provided", name)
	} else {
		return obj, nil
	}
}

func (c *ApplicationContainer) InvokeP(name string) interface{} {
	if obj, exists := c.deps[name]; !exists {
		panic(fmt.Errorf("depedency with name %s has not been provided", name))
	} else {
		return obj
	}
}

func (c *ApplicationContainer) InvokeByType(obj interface{}) (interface{}, error) {
	name := fmt.Sprint(reflect.TypeOf(obj))
	return c.Invoke(name)
}

func (c *ApplicationContainer) InvokeByTypeP(obj interface{}) interface{} {
	name := fmt.Sprint(reflect.TypeOf(obj))
	return c.InvokeP(name)
}

func (c *ApplicationContainer) execHook(fns []func() error) []error {
	depErrs := make([]error, 0)
	for _, fn := range fns {
		if err := fn(); err != nil {
			depErrs = append(depErrs, err)
		}
	}

	return depErrs
}

func (c *ApplicationContainer) start() []error {
	c.print("Starting dependencies...")
	return c.execHook(c.onStartHook)
}

func (c *ApplicationContainer) stop() []error {
	c.print("Stopping dependencies...")
	return c.execHook(c.onStopHook)
}

func (c *ApplicationContainer) Logger() *slog.Logger {
	return c.log
}
