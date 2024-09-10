package mypkg

// Greeter is a custom struct
type Greeter struct {
	message string
}

// NewGreeter constructor function returns a pointer to Greeter
func NewGreeter() *Greeter {
	return &Greeter{message: "Hello from Torpedo.di module"}
}

// Message exported method to get a message
func (m *Greeter) Message() string { return m.message }
