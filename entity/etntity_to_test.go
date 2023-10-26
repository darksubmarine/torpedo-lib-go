package entity_test

import (
	"github.com/darksubmarine/torpedo-lib-go/entity"
	"github.com/darksubmarine/torpedo-lib-go/ptr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntityTo_DTO(t *testing.T) {

	var dto = DTOEntity{}
	var ety = NewEntity()

	ety.SetId("qwerty0987")
	ety.SetCreated(123123123)
	ety.SetUpdated(9898988)

	ety.SetString("some string value")
	ety.SetInt(9)
	ety.SetBoolean(true)
	ety.SetSlice([]int{1, 2, 9, 8})
	ety.SetName("some name value")

	assert.Nil(t, entity.To(ety, &dto))

	assert.EqualValues(t, ety.Id(), ptr.ToString(dto.Id_))
	assert.EqualValues(t, ety.Created(), ptr.ToInt64(dto.Created_))
	assert.EqualValues(t, ety.Updated(), ptr.ToInt64(dto.Updated_))
	assert.EqualValues(t, ety.String(), ptr.ToString(dto.String_))
	assert.EqualValues(t, ety.Int(), ptr.ToInt(dto.Int_))
	assert.EqualValues(t, ety.Boolean(), ptr.ToBool(dto.Boolean_))
	assert.EqualValues(t, ety.Name(), ptr.ToString(dto.Name_))
	assert.EqualValues(t, ety.Slice(), dto.Slice_)

}

func TestEntityTo_DMO(t *testing.T) {
	var dmo = getEmptyDMO()
	var ety = NewEntity()

	ety.SetId("qwerty0987")
	ety.SetCreated(123123123)
	ety.SetUpdated(9898988)

	ety.SetString("some string value")
	ety.SetInt(9)
	ety.SetBoolean(true)
	ety.SetSlice([]int{1, 2, 9, 8})
	ety.SetName("some name value")

	assert.Nil(t, entity.To(ety, &dmo))

	assert.EqualValues(t, ety.Id(), dmo.Id_)
	assert.EqualValues(t, ety.Created(), dmo.Created_)
	assert.EqualValues(t, ety.Updated(), dmo.Updated_)
	assert.EqualValues(t, ety.String(), dmo.String_)
	assert.EqualValues(t, ety.Int(), dmo.Int_)
	assert.EqualValues(t, ety.Boolean(), dmo.Boolean_)
	assert.EqualValues(t, ety.Name(), dmo.Name_)
	assert.EqualValues(t, ety.Slice(), dmo.Slice_)
}
