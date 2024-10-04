package app

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

const (
	tagPrefix   = "torpedo.di"
	tagProvider = "provide"
	tagBinder   = "bind"
)

// iAppLifeCycle interface that defines the app life cycle
type iAppLifeCycle interface {
	//register(name string, obj interface{}, fn func() error) error
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

type RegisteredProvider struct {
	instance IProvider
	provides map[string]struct{}
	binds    map[string]struct{}
}

// ApplicationContainer implementation of the IApp interface. Handles the app life cycle and the
// main dependency container.
type ApplicationContainer struct {
	// Monitoring
	log *slog.Logger

	// signal channel
	sigs chan os.Signal

	// dynamic dependencies
	providers map[string]*RegisteredProvider
	values    sync.Map //map[string]reflect.Value

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
		providers:   map[string]*RegisteredProvider{},
		values:      sync.Map{}, //map[string]reflect.Value{},
	}

	signal.Notify(container.sigs, syscall.SIGINT, syscall.SIGTERM)

	return container
}

// WithProvider provides dependencies
func (c *ApplicationContainer) WithProvider(provider IProvider) *ApplicationContainer {
	var providerName = reflect.TypeOf(provider).String()

	// registering providers
	if _, ok := c.providers[providerName]; !ok {
		c.providers[providerName] = &RegisteredProvider{
			instance: provider,
			provides: map[string]struct{}{},
			binds:    map[string]struct{}{},
		}
	} else {
		panic(fmt.Sprintf("another provider has been registered with the name: %s", providerName))
	}

	for i := 0; i < reflect.TypeOf(provider).Elem().NumField(); i++ {
		field := reflect.TypeOf(provider).Elem().Field(i)
		// skipping register field.
		if field.Name == "BaseProvider" {
			continue
		}

		if tagVal, ok := field.Tag.Lookup(tagPrefix); ok && tagVal != "" {
			tagParts := strings.Split(tagVal, ",")
			instanceType := field.Type.String()
			if len(tagParts) == 2 {
				instanceType = strings.Replace(tagParts[1], "name=", "", 1)
			}

			if tagParts[0] == tagProvider { // it is a provider
				if _, ok := c.values.Load(instanceType); /*c.values[instanceType]*/ !ok {
					val := reflect.ValueOf(provider).Elem().Field(i)
					//c.values[instanceType] = reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem()
					c.values.Store(instanceType, reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem())
					c.providers[providerName].provides[instanceType] = struct{}{}
				} else {
					panic(fmt.Sprintf("provider %s is trying to provide an already registered dependency with name/type: %s", providerName, instanceType))
				}
			} else if tagParts[0] == tagBinder { // it is a bound instance
				c.providers[providerName].binds[instanceType] = struct{}{}
			}
		}
	}
	return c
}

func (c *ApplicationContainer) lookupProviderFor(bind string) string {
	for pName, provider := range c.providers {
		if _, ok := provider.provides[bind]; ok {
			return pName
		}
	}

	return ""
}

func (c *ApplicationContainer) provideDependencies() {
	depGraph := NewDependencyGraph()

	// adding vertexes
	for pName, _ := range c.providers {
		depGraph.AddVertex(pName)
	}

	for pName, provider := range c.providers {
		for bind, _ := range provider.binds {
			provDependsOn := c.lookupProviderFor(bind)
			if provDependsOn == "" {
				panic(fmt.Errorf("%w whithin %s for %s", ErrDependencyNotProvided, pName, bind))
			}
			_ = depGraph.AddEdge(pName, provDependsOn)
		}
	}

	if lst, err := depGraph.TopologicalSort(); err == nil {
		for i := 0; i < len(lst); i++ {
			provName := lst[i]

			// ensure all binds are set.
			pInstance := c.providers[provName].instance
			for i := 0; i < reflect.TypeOf(pInstance).Elem().NumField(); i++ {
				field := reflect.TypeOf(pInstance).Elem().Field(i)
				// skipping register field.
				if field.Name == "BaseProvider" {
					continue
				}

				if tagVal, ok := field.Tag.Lookup(tagPrefix); ok && tagVal != "" {
					tagParts := strings.Split(tagVal, ",")
					if tagParts[0] == tagBinder {
						instanceType := field.Type.String()
						if len(tagParts) == 2 {
							instanceType = strings.Replace(tagParts[1], "name=", "", 1)
						}

						if val, exists := c.values.Load(instanceType); /*c.values[instanceType]*/ exists {
							value := val.(reflect.Value)
							if (value.Kind() == reflect.Pointer || value.Kind() == reflect.Func) && value.IsNil() {
								panic(fmt.Sprintf("The binded field named '%s' in your provider %s cannot be nil. \n> "+
									"hint: Check if the provider has been initialized into the Provider() method in your %s", field.Name, provName, c.lookupProviderFor(instanceType)))
							}

							fv := reflect.ValueOf(pInstance).Elem().Field(i)
							// fv is not writable because it is no-exportable. So, using unsafe.Pointer to hack it!
							fv = reflect.NewAt(fv.Type(), unsafe.Pointer(fv.UnsafeAddr())).Elem()

							if value.Kind() == reflect.Func {
								out := value.Call([]reflect.Value{reflect.ValueOf(c)})
								fv.Set(reflect.ValueOf(out[0].Interface()))
							} else {
								fv.Set(reflect.ValueOf(value.Interface()))
							}
						} else {
							panic(fmt.Sprintf("provider not found for %s at %s", instanceType, provName))
						}
					}
				}
			}

			if err := pInstance.Provide(c); err != nil {
				panic(fmt.Sprintf("error providing %s with %s", provName, err))
			} else {
				// adding LifeCycle hooks
				c.onStart(pInstance.OnStart())
				c.onStop(pInstance.OnStop())
			}
		}
	} else {
		panic(fmt.Sprintf("Dependency cycle: %s\n", err))
	}
}

// Run app method to starts the application life cycle
func (c *ApplicationContainer) Run() {
	defer func() {
		if errs := c.stop(); len(errs) > 0 {
			c.exitWithErrors("Some container dependencies cannot be stopped due to some errors", errs)
		}
	}()

	// check cycle dependencies and register providers.
	c.provideDependencies()

	if startErrs := c.start(); len(startErrs) > 0 {
		c.exitWithErrors("Dependencies container cannot START due to some deps starting errors", startErrs)
	}

	c.print("Application started successfully... (waiting for signal termination)")
	<-c.sigs
	c.printnln("Application terminated by signal")
}

// Register exposed method to register object by name without linked hooks
func (c *ApplicationContainer) Register(name string, obj interface{}) error {
	tpy := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	instanceType := tpy.String()
	if name != "" {
		instanceType = name
	}
	//c.values[instanceType] = val //reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem()
	c.values.Store(instanceType, val)
	return nil
}

// Invoke fetch a named dependency from the container or error if not exists.
func (c *ApplicationContainer) Invoke(name string) (interface{}, error) {
	if obj, exists := c.values.Load(name); /*c.values[name]*/ !exists {
		return nil, fmt.Errorf("%w {name=%s}", ErrDependencyNotProvided, name)
	} else {
		_obj := obj.(reflect.Value)
		return _obj.Interface(), nil
	}
}

// InvokeP fetch a named dependency from the container and panic if not exists.
func (c *ApplicationContainer) InvokeP(name string) interface{} {
	if obj, exists := c.values.Load(name); /*c.values[name]*/ !exists {
		panic(fmt.Errorf("depedency with name %s has not been provided", name))
	} else {
		_obj := obj.(reflect.Value)
		return _obj.Interface()
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

// Logger returns the container logger useful to log something into your providers
func (c *ApplicationContainer) Logger() *slog.Logger {
	return c.log
}

// ------------------------
// -- Unexported methods --
// ------------------------

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

// exitWithErrors finalize the app with an os.Exit(1)
func (c *ApplicationContainer) exitWithErrors(msg string, errs []error) {
	fmt.Println(msg)
	for i, err := range errs {
		fmt.Printf("  %d - %s\n", i, err)
	}
	os.Exit(1)
}

// onStart register the OnStart hook from providers
func (c *ApplicationContainer) onStart(fn func() error) {
	c.onStartHook = append(c.onStartHook, fn)
}

// onStop register the OnStop hook from providers
func (c *ApplicationContainer) onStop(fn func() error) {
	c.onStopHook = append(c.onStopHook, fn)
}

// print verbose method to print to std out
func (c *ApplicationContainer) print(format string, a ...any) {
	fmt.Println("[TPDO]", fmt.Sprintf(format, a...))
}

// printnln verbose method to printnln to std out
func (c *ApplicationContainer) printnln(format string, a ...any) {
	fmt.Println("\n[TPDO]", fmt.Sprintf(format, a...))
}
