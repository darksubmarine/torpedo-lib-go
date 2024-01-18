package entity_test

import (
	"github.com/darksubmarine/torpedo-lib-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDMOToEntity_FullDMO(t *testing.T) {

	var dmo = getDMO(fullDMO)
	var ety = NewEntity()
	assert.Nil(t, entity.From(&dmo, ety))

	assert.Equal(t, ety.Id(), dmo.Id_)
	assert.Equal(t, ety.Created(), dmo.Created_)
	assert.Equal(t, ety.Updated(), dmo.Updated_)
	assert.Equal(t, ety.String(), dmo.String_)
	assert.Equal(t, ety.Int(), dmo.Int_)
	assert.Equal(t, ety.Boolean(), dmo.Boolean_)
	assert.Equal(t, ety.Slice(), dmo.Slice_)
	assert.Equal(t, ety.Name(), dmo.Name_)
}

func TestDMOToEntity_PartialDMO(t *testing.T) {

	var dmo = getDMO(partialDMO)
	var ety = NewEntity()
	assert.Nil(t, entity.From(&dmo, ety))

	assert.Equal(t, ety.Id(), dmo.Id_)
	assert.Equal(t, ety.Created(), dmo.Created_)
	assert.Equal(t, ety.Updated(), dmo.Updated_)
	assert.Equal(t, ety.String(), dmo.String_)
	assert.Equal(t, ety.Int(), dmo.Int_)
	assert.Equal(t, ety.Boolean(), dmo.Boolean_)
	assert.Equal(t, ety.Slice(), dmo.Slice_)
	assert.Equal(t, ety.Name(), dmo.Name_)
}
