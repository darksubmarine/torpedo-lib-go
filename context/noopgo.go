package context

import "time"

// NoopGoContext satisfy the go builtin context.Context interface
type NoopGoContext struct{}

func (d *NoopGoContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (d *NoopGoContext) Done() <-chan struct{} {
	return nil
}

func (d *NoopGoContext) Err() error {
	return nil
}

func (d *NoopGoContext) Value(key any) any {
	return nil
}
