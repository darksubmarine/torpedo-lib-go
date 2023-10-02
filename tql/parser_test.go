package tql

import (
	"encoding/json"
	"testing"
)

func TestFilterItem(t *testing.T) {

	_json := `
{
    "filter": [
      {
        "field": "age",
        "operator": ">=",
        "value": 18
      },
      {
        "field": "plan",
        "operator": "[?]", 
        "value": ["silver","gold","platinum"]
      },
      {
        "field": "date",
        "operator": ">?<",
        "value": ["silver","gold"]
      }
    ],
    
    "sort": [
      {"field": "name", "type": "asc"},
      {"field": "social_number", "type": "desc"}
    ],
    
    "projection": ["id", "name"],
    
    "pagination": {
      "items": 10,
      "page": 5,
      "cursor": "QWER!@#$ASDF"
    }
}`

	var q Query

	if err := json.Unmarshal([]byte(_json), &q); err != nil {
		t.Error(err)
	}
}
