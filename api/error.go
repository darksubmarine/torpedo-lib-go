package api

const (
	// E4001 Error binding JSON
	E4001 = "4001"

	// E4002 Partial entity incomplete some field is missing
	E4002 = "4002"

	// E4003 Entity building from DTO
	E4003 = "4003"

	// E4004 Entity not found
	E4004 = "4004"

	// E4005 - Entity Query (TQL) error
	E4005 = "4005"

	// E5001 - Entity creation error
	E5001 = "5001"

	// E5002 - Entity update error
	E5002 = "5002"

	// E5003 - Entity read error
	E5003 = "5003"

	// E5004 - Entity remove error
	E5004 = "5004"

	// E5005 - Entity Query (TQL) error
	E5005 = "5005"
)

// Error api error message
type Error struct {
	Code string `json:"code"`
	Msg  string `json:"error"`
}

// NewError returns an API error message
func NewError(code string, err error) Error {
	return Error{Code: code, Msg: err.Error()}
}

// ErrorBindingJSON returns a 4001 error
func ErrorBindingJSON(err error) Error {
	return NewError(E4001, err)
}

// ErrorPartialEntityIncomplete returns a 4002 error
func ErrorPartialEntityIncomplete(err error) Error {
	return NewError(E4002, err)
}

// ErrorBuildingEntityFromDTO returns a 4003 error
func ErrorBuildingEntityFromDTO(err error) Error {
	return NewError(E4003, err)
}

// ErrorNotFound returns a 4004 error
func ErrorNotFound(err error) Error {
	return NewError(E4004, err)
}

// ErrorEntityQueryByUser returns a 4005 error
func ErrorEntityQueryByUser(err error) Error {
	return NewError(E4005, err)
}

// ErrorEntityCreation returns a 5001 error
func ErrorEntityCreation(err error) Error {
	return NewError(E5001, err)
}

// ErrorEntityUpdate returns a 5002 error
func ErrorEntityUpdate(err error) Error {
	return NewError(E5002, err)
}

// ErrorEntityRead returns a 5003 error
func ErrorEntityRead(err error) Error {
	return NewError(E5003, err)
}

// ErrorEntityRemove returns a 5004 error
func ErrorEntityRemove(err error) Error {
	return NewError(E5004, err)
}

// ErrorEntityQuery returns a 5005 error
func ErrorEntityQuery(err error) Error {
	return NewError(E5005, err)
}
