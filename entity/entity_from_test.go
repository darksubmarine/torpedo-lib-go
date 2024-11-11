package entity_test

import (
	"encoding/json"
	"github.com/darksubmarine/torpedo-lib-go/entity"
	"github.com/darksubmarine/torpedo-lib-go/ptr"
	"github.com/stretchr/testify/assert"
	"testing"
)

var fullDTO = []byte(`{
	"id": "qwerty-123456",
	"created": 1697649758238,
	"updated": 1697649758123,
	"string": "some-string-value",
	"number": 515,
	"boolean": true,
	"slice": [1,2,3,4,5,6,7,8,9,0],
	"name": "some-name-value",

	"hasone": {"id": "qwerty-123456", "created": 1697649758238, "updated": 1697649758123},
	"hasmany": [{"id": "qwerty-123456", "created": 1697649758238, "updated": 1697649758123}]
}`)

var partialDTO = []byte(`{
	"id": "qwerty-123456",
	"created": 1697649758238,
	
	"string": "some-string-value",
	"number": 515,
	"boolean": true
}`)

var validateFieldInvalidDTO = []byte(`{"inlist": "invalid"}`)
var validateFieldValidDTO = []byte(`{"inlist": "valid"}`)

func getDTO(data []byte) DTOEntity {
	var dto DTOEntity
	_ = json.Unmarshal(data, &dto)
	return dto
}

func TestDTOToEntity_FullDTO(t *testing.T) {

	var dto = getDTO(fullDTO)
	var ety = NewEntity()
	assert.Nil(t, entity.From(&dto, ety))

	assert.EqualValues(t, ety.Id(), ptr.ToString(dto.Id_))
	assert.EqualValues(t, ety.Created(), ptr.ToInt64(dto.Created_))
	assert.EqualValues(t, ety.Updated(), ptr.ToInt64(dto.Updated_))
	assert.EqualValues(t, ety.String(), ptr.ToString(dto.String_))
	assert.EqualValues(t, ety.Int(), ptr.ToInt(dto.Int_))
	assert.EqualValues(t, ety.Boolean(), ptr.ToBool(dto.Boolean_))
	assert.EqualValues(t, ety.Slice(), dto.Slice_)
	assert.EqualValues(t, ety.Name(), ptr.ToString(dto.Name_))

	// Relationship is a decorated field only populated at service.Read method
	// This one should not be set via entity.From
	assert.Nil(t, ety.HasOne())
	assert.Nil(t, ety.HasMany())
}

func TestDTOToEntity_PartialDTO(t *testing.T) {

	var dto = getDTO(partialDTO)
	var ety = NewEntity()
	assert.Nil(t, entity.From(&dto, ety))

	assert.EqualValues(t, ety.Id(), ptr.ToString(dto.Id_))
	assert.EqualValues(t, ety.Created(), ptr.ToInt64(dto.Created_))
	assert.EqualValues(t, ety.Updated(), ptr.ToInt64(dto.Updated_))
	assert.EqualValues(t, ety.String(), ptr.ToString(dto.String_))
	assert.EqualValues(t, ety.Int(), ptr.ToInt(dto.Int_))
	assert.EqualValues(t, ety.Boolean(), ptr.ToBool(dto.Boolean_))
	assert.EqualValues(t, ety.Slice(), dto.Slice_)
	assert.EqualValues(t, ety.Name(), ptr.ToString(dto.Name_))
}

func TestDTOToEntity_ValidateFieldError(t *testing.T) {

	var dto = getDTO(validateFieldInvalidDTO)
	var ety = NewEntity()
	err := entity.From(&dto, ety)
	assert.Error(t, err)
}

func TestDTOToEntity_ValidateFieldOK(t *testing.T) {

	var dto = getDTO(validateFieldValidDTO)
	var ety = NewEntity()
	assert.Nil(t, entity.From(&dto, ety))
	assert.EqualValues(t, ety.Inlist(), ptr.ToString(dto.Inlist_))
}
