package mongodb_utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdate(t *testing.T) {
	bsonD := ToBSONDocument(&fullDMO, "Id_", "Created_")

	toCheck := map[string]interface{}{
		"updated": 1697649758123,
		"string":  "some-string-value",
		"number":  515,
		"boolean": true,
		"slice":   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		"custom":  "some-name-value",
	}

	assert.Len(t, bsonD, 6)
	for _, d := range bsonD {
		assert.NotEqual(t, "_id", d.Key)
		assert.NotEqual(t, "created", d.Key)
		assert.EqualValues(t, toCheck[d.Key], d.Value)
	}

}
