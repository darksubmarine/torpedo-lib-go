package rels

// BelongsToFn function to fetch paginated items
type BelongsToFn[E any] func(string, int64, int64) ([]E, error)

// BelongsTo function to fetch all items that belongs to a collection
func BelongsTo[ENTITY any](fn BelongsToFn[ENTITY], id string, batch int64) ([]ENTITY, error) {
	var toRet = make([]ENTITY, 0)
	// starting from i=1 because it is first page.
	for i := 1; ; i++ {
		if entities, err := fn(id, batch, int64(i)); err == nil && len(entities) > 0 {
			toRet = append(toRet, entities...)
		} else if err == nil && len(entities) == 0 {
			return toRet, nil
		} else if err != nil {
			return nil, err
		} else {
			return nil, nil
		}
	}
}
