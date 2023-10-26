package data_struct

func SkipMap(skip ...string) map[string]struct{} {
	sm := map[string]struct{}{}
	for _, v := range skip {
		sm[v] = struct{}{}
	}
	return sm
}
