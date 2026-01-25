package config

type ConfigRegistry struct {
	modules map[string]ConfigModule
}

func NewConfigRegistry() *ConfigRegistry {
	return &ConfigRegistry{
		modules: make(map[string]ConfigModule),
	}
}

func (cr *ConfigRegistry) RegisterModule(name string, module ConfigModule) {
	module.SetConfigName(name)
	cr.modules[name] = module
}

func (cr *ConfigRegistry) GetModule(name string) ConfigModule {
	return cr.modules[name]
}
