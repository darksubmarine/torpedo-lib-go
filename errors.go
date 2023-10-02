package torpedo_lib

import "errors"

var (
	// ErrEmptyId the entity ID is empty
	ErrEmptyId = errors.New("the entity ID is empty")

	// ErrIdConvertion the generated id by the repository cannot be converted to string
	ErrIdConvertion = errors.New("the generated id by the repository cannot be converted to string")

	// ErrIdNotFound entity not found with the given id
	ErrIdNotFound = errors.New("entity not found with the given id")

	//ErrBindingDTO the given object cannot be bind to DTO
	ErrBindingDTO = errors.New("the given object cannot be bind to DTO")

	// ErrNilEntity the given entity is nil
	ErrNilEntity = errors.New("the given entity is nil")
)

func IsError(e error) bool {
	if e != nil {
		return true
	}
	return false
}

type Err struct {
	err error
}

func Error(err error) *Err {
	return &Err{err: err}
}

func (e *Err) Error() string {
	return e.err.Error()
}
