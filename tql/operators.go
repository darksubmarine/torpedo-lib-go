package tql

const (
	FilterTypeAny = "any"
	FilterTypeAll = "all"

	OpNEQ = "!="
	OpEQ  = "=="
	OpGT  = ">"
	OpGTE = ">="
	OpLT  = "<"
	OpLTE = "<="

	OpBTLimits     = "[n]"
	OpBTNoLimits   = "(n)"
	OpBTLeftLimit  = "[n)"
	OpBTRightLimit = "(n]"

	OpIN = "[?]"

	OpPrefix   = "s.."
	OpSuffix   = "..s"
	OpContains = ".s."
)

func isSimpleOperator(op string) bool {
	switch op {
	case OpNEQ, OpEQ, OpGT, OpGTE, OpLT, OpLTE:
		return true
	}
	return false
}

func isBetweenOperator(op string) bool {
	switch op {
	case OpBTLimits, OpBTNoLimits, OpBTLeftLimit, OpBTRightLimit:
		return true
	}
	return false
}

func isInListOperator(op string) bool {
	if op == OpIN {
		return true
	}
	return false
}

func isStringOperator(op string) bool {
	switch op {
	case OpPrefix, OpSuffix, OpContains:
		return true
	}
	return false
}
