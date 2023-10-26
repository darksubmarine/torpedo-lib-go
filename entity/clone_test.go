package entity_test

import (
	"github.com/darksubmarine/torpedo-lib-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntityClone(t *testing.T) {
	var etyDest = NewEntity()
	var ety = NewEntity()

	ety.SetId("qwerty0987")
	ety.SetCreated(123123123)
	ety.SetUpdated(9898988)

	ety.SetString("some string value")
	ety.SetInt(9)
	ety.SetBoolean(true)
	ety.SetSlice([]int{1, 2, 9, 8})
	ety.SetName("some name value")

	assert.Nil(t, entity.Clone(ety, etyDest))

	assert.EqualValues(t, ety.Id(), etyDest.Id())
	assert.EqualValues(t, ety.Created(), etyDest.Created())
	assert.EqualValues(t, ety.Updated(), etyDest.Updated())
	assert.EqualValues(t, ety.String(), etyDest.String())
	assert.EqualValues(t, ety.Int(), etyDest.Int())
	assert.EqualValues(t, ety.Boolean(), etyDest.Boolean())
	assert.EqualValues(t, ety.Name(), etyDest.Name())
	assert.EqualValues(t, ety.Slice(), etyDest.Slice())
}
