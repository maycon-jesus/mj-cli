package config

type ReadOnlyConfig interface {
	Get(key string) (any, bool)
	Has(key string) bool
}
