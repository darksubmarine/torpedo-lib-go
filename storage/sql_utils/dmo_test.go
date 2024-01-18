package sql_utils_test

import "github.com/darksubmarine/torpedo-lib-go/crypto"

// EntityDMO base entity Data Mapper Object
type EntityDMO struct {
	cryptoKey []byte // the-key-has-to-be-32-bytes-long!
}

// NewEntityDMO constructor function of EntityDMO
func NewEntityDMO(key []byte) *EntityDMO {
	return &EntityDMO{cryptoKey: key}
}

// EncryptString encrypts the given string
func (dmo *EntityDMO) EncryptString(value string) (string, error) {
	return crypto.EncodeString(dmo.cryptoKey, value)
}

// DecryptString decrypt the given string
func (dmo *EntityDMO) DecryptString(value string) (string, error) {
	return crypto.DecodeString(dmo.cryptoKey, value)
}

type EntityDMOBase struct {
	*EntityDMO

	Id_      string `db:"id"`
	Created_ int64  `db:"created"`
	Updated_ int64  `db:"updated"`

	EntityDMOBasePartial `db:"inline"`
}

type EntityDMOBasePartial struct {
	String_  string `db:"string"`
	Int_     int    `db:"number"`
	Boolean_ bool   `db:"boolean"`
	Slice_   []int  `db:"slice"`

	DMO `db:"inline"`
}

type DMO struct {
	Custom_ string `db:"custom"`
}

var fullDMO = EntityDMOBase{
	EntityDMO: NewEntityDMO([]byte("some =-key")),
	Id_:       "qwerty-123456",
	Created_:  1697649758238,
	Updated_:  1697649758123,
	EntityDMOBasePartial: EntityDMOBasePartial{
		String_:  "some-string-value",
		Int_:     515,
		Boolean_: true,
		Slice_:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		DMO:      DMO{Custom_: "some-name-value"},
	}}
