package entity_test

import (
	"encoding/json"
	"github.com/darksubmarine/torpedo-lib-go/crypto"
)

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

	Id_      string `json:"id"`
	Created_ int64  `json:"created"`
	Updated_ int64  `json:"updated"`

	EntityDMOBasePartial
}

type EntityDMOBasePartial struct {
	String_  string `json:"string,omitempty" read_method:"String"`
	Int_     int    `json:"number,omitempty"`
	Boolean_ bool   `json:"boolean,omitempty"`
	Slice_   []int  `json:"slice"`
}

type EntityDMOJSON struct {
	EntityDMOBase
	Name_ string `json:"name,omitempty"`
}

var fullDMO = []byte(`{
	"id": "qwerty-123456",
	"created": 1697649758238,
	"updated": 1697649758123,
	"string": "some-string-value",
	"number": 515,
	"boolean": true,
	"slice": [1,2,3,4,5,6,7,8,9,0],
	"name": "some-name-value"
}`)

var partialDMO = []byte(`{
	"id": "qwerty-123456",
	"created": 1697649758238,
	
	"string": "some-string-value",
	"number": 515,
	"boolean": true
}`)

func getDMO(data []byte) EntityDMOJSON {
	var dmo = EntityDMOJSON{EntityDMOBase: EntityDMOBase{EntityDMO: NewEntityDMO([]byte("some =-key")), EntityDMOBasePartial: EntityDMOBasePartial{}}}
	_ = json.Unmarshal(data, &dmo)
	return dmo
}

func getEmptyDMO() EntityDMOJSON {
	return EntityDMOJSON{EntityDMOBase: EntityDMOBase{EntityDMO: NewEntityDMO([]byte("some =-key")), EntityDMOBasePartial: EntityDMOBasePartial{}}}
}
