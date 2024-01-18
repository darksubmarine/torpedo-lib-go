package tql

import "errors"

var (
	// ErrTQLNotSupported Torpedo Query Language is not supported by this repository
	ErrTQLNotSupported = errors.New("Torpedo Query Language is not supported by this repository")

	// ErrEmptyFilter filter cannot be empty
	ErrEmptyFilter = errors.New("filter cannot be empty")

	// ErrInvalidOperator the given operator is invalid
	ErrInvalidOperator = errors.New("the given operator is invalid")

	// ErrInvalidValue the given value is invalid
	ErrInvalidValue = errors.New("the given value is invalid")

	// ErrInvalidValueBTLen the given value must be a list of 2 elements (left and right limits)
	ErrInvalidValueBTLen = errors.New("the given value must be a list of 2 elements (left and right limits)")

	// ErrOperationNotSupported the given operation is not supported for the given field
	ErrOperationNotSupported = errors.New("the given operation is not supported for the given field")

	// ErrInvalidFieldName the given field name is invalid
	ErrInvalidFieldName = errors.New("the given field name is invalid")

	// ErrInvalidFieldNameAtProjection the given field name into projection list is invalid
	ErrInvalidFieldNameAtProjection = errors.New("the given field name into projection list is invalid")

	// ErrInvalidFieldNameAtSort the given field name into sort list is invalid
	ErrInvalidFieldNameAtSort = errors.New("the given field name into sort list is invalid")

	// ErrQueryResultObjectBuild something happens building a QRO from the entity
	ErrQueryResultObjectBuild = errors.New("something happens building a QRO from the entity")

	// ErrInvalidPaginationType the pagination must be 'cursor' or 'offset'
	ErrInvalidPaginationType = errors.New("the pagination must be 'cursor' or 'offset'")

	// ErrInvalidSortFieldNotProjectionMember sort field must be a member of the projection field list
	ErrInvalidSortFieldNotProjectionMember = errors.New("sort field must be a member of the projection field list")
)
