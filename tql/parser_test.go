package tql

import (
	"encoding/json"
	"testing"
)

func TestFilterItem(t *testing.T) {

	_json := `
{
  "filter": {
        "type": "all",
        "fields": [
            {
                "field": "created",
                "operator": ">=",
                "value": 1666875856369
            },
            {
                "field": "age",
                "operator": ">=",
                "value": 18
            },
            {
                "field": "plan",
                "operator": "[?]",
                "value": [
                  "silver",
                  "gold",
                  "platinum"
                ]
            },
            {
                "field": "date",
                "operator": "[?]",
                "value": [
                  1665757374948,
                  1664752374987
                ]
            }
        ]
    },

  "projection": ["name", "id","updated"],

  "pagination":{
        "items":30,
        "offset":{
          "page":10,
          "sort": [{"field": "name", "type": "desc"}]
        }
      }

}`

	var q Query

	if err := json.Unmarshal([]byte(_json), &q); err != nil {
		t.Error(err)
	}
}
