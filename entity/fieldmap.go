package entity

type FieldMap map[string]string

func (fm FieldMap) HasField(name string) bool {
	_, ok := fm[name]
	return ok
}

func (fm FieldMap) FieldType(name string) (string, bool) {
	if val, ok := fm[name]; ok {
		return val, true
	}

	return "", false
}
