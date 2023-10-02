package tql

func NewResult(data []interface{}, count int, prev interface{}, next interface{}) *Result {
	return &Result{
		Meta: ResultMeta{
			Count: int64(count),
			//Prev:  prev,
			Next: next,
		},
		Data: data,
	}
}

type Result struct {
	Meta ResultMeta    `json:"meta"`
	Data []interface{} `json:"data"`
}

type ResultMeta struct {
	Count int64 `json:"count"`
	//Prev  interface{} `json:"prev"`
	Next interface{} `json:"nextToken"`
}
