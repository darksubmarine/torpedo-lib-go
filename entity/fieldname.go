package entity

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// FieldNameToCode given a field name returns the coded attribute field name.
// Used to map a field name to its attribute in DMO, DTO and QRO objects.
func FieldNameToCode(name string) string {
	return fmt.Sprintf("%s_", cases.Title(language.English, cases.NoLower).String(name))
}
