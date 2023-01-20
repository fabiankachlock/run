package loader

type Loader interface {
	GetScope() string
	LoadConfig(dir string) map[string]string
}
