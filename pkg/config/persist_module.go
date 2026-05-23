package config

type PersistModule interface {
	Save(configs map[string]Entry) error
	Load() error
	Get(key string) (any, bool)
}
