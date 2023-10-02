package entity

import "github.com/darksubmarine/torpedo-lib-go/trn"

// TRN returns a TRN of entity kind
func TRN(name string, id string) *trn.TRN {
	return trn.New("entity", name, id)
}
