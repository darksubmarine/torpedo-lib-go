package conf

type Loader interface {
	Load(m Map) Map
}
