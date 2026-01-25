package config

import "github.com/spf13/pflag"

type ConfigModule interface {
	SetConfigName(name string)
	ReadInConfig() error
	WriteConfig() error

	SetEnvPrefix(prefix string)

	SetDefault(key string, value interface{})
	Set(key string, value interface{})

	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool

	AllSettings() map[string]interface{}
	AllKeys() []string

	BindPFlag(key string, flag *pflag.Flag) error
}
