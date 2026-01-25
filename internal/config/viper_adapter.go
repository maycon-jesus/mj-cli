// config/viper_adapter.go
package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ViperAdapter implementa ConfigModuleInterface usando o Viper
type ViperAdapter struct {
	v *viper.Viper
}

func setDefaultViperConfig(v *viper.Viper, appName string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configDirHome := path.Join(homeDir, fmt.Sprintf(".%s", appName))
	err = os.MkdirAll(configDirHome, os.ModePerm)
	if err != nil {
		panic(err)
	}

	v.AddConfigPath(".")
	v.AddConfigPath(configDirHome)
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}

// NewViperAdapter cria uma nova instância do adapter
func NewViperAdapter(appName string) ConfigModule {
	v := viper.New()

	setDefaultViperConfig(v, appName)

	return &ViperAdapter{
		v,
	}
}

// Implementação dos métodos da interface
func (va *ViperAdapter) SetConfigName(name string) {
	va.v.SetConfigName(name)
}

func (va *ViperAdapter) SetEnvPrefix(prefix string) {
	va.v.SetEnvPrefix(prefix)
}

func (va *ViperAdapter) ReadInConfig() error {
	return va.v.ReadInConfig()
}

func (va *ViperAdapter) WriteConfig() error {
	err := va.v.WriteConfig()
	if err != nil {
		return va.v.SafeWriteConfig()
	}
	return nil
}

func (va *ViperAdapter) SetDefault(key string, value interface{}) {
	va.v.SetDefault(key, value)
}

func (va *ViperAdapter) Set(key string, value interface{}) {
	va.v.Set(key, value)
}

func (va *ViperAdapter) Get(key string) interface{} {
	return va.v.Get(key)
}

func (va *ViperAdapter) GetString(key string) string {
	return va.v.GetString(key)
}

func (va *ViperAdapter) GetInt(key string) int {
	return va.v.GetInt(key)
}

func (va *ViperAdapter) GetBool(key string) bool {
	return va.v.GetBool(key)
}

func (va *ViperAdapter) AllSettings() map[string]interface{} {
	return va.v.AllSettings()
}

func (va *ViperAdapter) AllKeys() []string {
	return va.v.AllKeys()
}

func (va *ViperAdapter) BindPFlag(key string, flag *pflag.Flag) error {
	return va.v.BindPFlag(key, flag)
}
