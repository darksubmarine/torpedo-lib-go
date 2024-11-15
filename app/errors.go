package app

import (
	"errors"
)

var (
	// ErrDependencyNotProvided dependency has not been provided
	ErrDependencyNotProvided = errors.New("dependency has not been provided")

	// ErrNilDependency the provided dependency cannot be nil
	ErrNilDependency = errors.New("the provided dependency cannot be nil")

	// ErrDependencyAlreadyProvided dependency already provided
	ErrDependencyAlreadyProvided = errors.New("dependency already provided")
)
