package loader

type Loader interface {
	LoadConfig(dir string) (map[string]string, error)
}
