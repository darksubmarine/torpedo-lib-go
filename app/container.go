package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

// iAppLifeCycle interface that defines the app life cycle
type iAppLifeCycle interface {
	register(name string, obj interface{}, fn func() error) error
	onStart(func() error)
	onStop(func() error)
	start() []error
	stop() []error
}

// IContainerMonitor interface that defines the container monitoring capabilities
type IContainerMonitor interface {
	Logger() *slog.Logger
}

// IContainer exposed interface that defines the container methods
type IContainer interface {
	IContainerMonitor
	Register(name string, obj interface{}) error
	Invoke(name string) (interface{}, error)
	InvokeP(name string) interface{}
	InvokeByType(obj interface{}) (interface{}, error)
	InvokeByTypeP(obj interface{}) interface{}
}

// IApp application interface
type IApp interface {
	iAppLifeCycle
	IContainer
	Run()
}

// ApplicationContainer implementation of the IApp interface. Handles the app life cycle and the
// main dependency container.
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

// NewApplicationContainer app container constructor
func NewApplicationContainer(opts ContainerOpts) *ApplicationContainer {
	return NewContainer(opts)
}

// NewContainer app container constructor
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

// print verbose method to print to std out
func (c *ApplicationContainer) print(format string, a ...any) {
	fmt.Println("[TPDO]", fmt.Sprintf(format, a...))
}

// printnln verbose method to printnln to std out
func (c *ApplicationContainer) printnln(format string, a ...any) {
	fmt.Println("\n[TPDO]", fmt.Sprintf(format, a...))
}

// execProvider call the provider.Provide method used to register a new dependency
func (c *ApplicationContainer) execProvider(provider IProvider) interface{} {
	if obj, err := provider.Provide(c); err != nil {
		panic(err)
	} else {
		return obj
	}
	return nil
}

// register the provided dependency object and execute the hook OnRegister
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

// WithProvider useful to provide an object by type
func (c *ApplicationContainer) WithProvider(provider IProvider) *ApplicationContainer {
	return c.WithNamedProvider("", provider)
}

// WithNamedProvider useful to provide an object by name
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

// exitWithErrors finalize the app with an os.Exit(1)
func (c *ApplicationContainer) exitWithErrors(msg string, errs []error) {
	fmt.Println(msg)
	for i, err := range errs {
		fmt.Printf("  %d - %s\n", i, err)
	}
	os.Exit(1)
}

// Run app method to starts the application life cycle
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

// onStart register the OnStart hook from providers
func (c *ApplicationContainer) onStart(fn func() error) {
	c.onStartHook = append(c.onStartHook, fn)
}

// onStart register the OnStop hook from providers
func (c *ApplicationContainer) onStop(fn func() error) {
	c.onStopHook = append(c.onStopHook, fn)
}

// Register exposed method to register object by name without linked hooks
func (c *ApplicationContainer) Register(name string, obj interface{}) error {
	if obj == nil {
		return fmt.Errorf("%w {name=%s}", ErrNilDependency, name)
	}

	if _, exists := c.deps[name]; exists {
		return fmt.Errorf("%w {name=%s}", ErrDependencyAlreadyProvided, name)
	}

	c.print("Registering dependency for %s", name)
	c.deps[name] = obj
	return nil
}

// Invoke fetch a named dependency from the container or error if not exists.
func (c *ApplicationContainer) Invoke(name string) (interface{}, error) {
	if obj, exists := c.deps[name]; !exists {
		return nil, fmt.Errorf("%w {name=%s}", ErrDependencyNotProvided, name)
	} else {
		return obj, nil
	}
}

// InvokeP fetch a named dependency from the container and panic if not exists.
func (c *ApplicationContainer) InvokeP(name string) interface{} {
	if obj, exists := c.deps[name]; !exists {
		panic(fmt.Errorf("depedency with name %s has not been provided", name))
	} else {
		return obj
	}
}

// InvokeByType fetch an object registered by type or error if not exists.
func (c *ApplicationContainer) InvokeByType(obj interface{}) (interface{}, error) {
	name := fmt.Sprint(reflect.TypeOf(obj))
	return c.Invoke(name)
}

// InvokeByTypeP fetch an object registered by type or panic if not exists.
func (c *ApplicationContainer) InvokeByTypeP(obj interface{}) interface{} {
	name := fmt.Sprint(reflect.TypeOf(obj))
	return c.InvokeP(name)
}

// execHook execute the given hooks.
func (c *ApplicationContainer) execHook(fns []func() error) []error {
	depErrs := make([]error, 0)
	for _, fn := range fns {
		if err := fn(); err != nil {
			depErrs = append(depErrs, err)
		}
	}

	return depErrs
}

// start executes all the provided OnStart hooks
func (c *ApplicationContainer) start() []error {
	c.print("Starting dependencies...")
	return c.execHook(c.onStartHook)
}

// stop executes all the provided OnStart hooks
func (c *ApplicationContainer) stop() []error {
	c.print("Stopping dependencies...")
	return c.execHook(c.onStopHook)
}

// Logger returns the container logger useful to log something into your providers
func (c *ApplicationContainer) Logger() *slog.Logger {
	return c.log
}
