package torpedo_lib

import (
	"github.com/darksubmarine/torpedo-lib-go/trn"
)

// NewTRN creates a new TRN instance
func NewTRN(kind, name, id string) *trn.TRN {
	return trn.New(kind, name, id)
}
