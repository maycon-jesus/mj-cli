package config

import (
	"fmt"
	"strings"
)

type Entry struct {
	Description  string
	Value        any
	DefaultValue any
}

func (e Entry) effective() any {
	if e.Value != nil {
		return e.Value
	}
	return e.DefaultValue
}

type ConfigManager struct {
	entries   map[string]Entry
	persister PersistModule
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		entries:   make(map[string]Entry),
		persister: nil,
	}
}

func (cm *ConfigManager) WithPersister(p PersistModule) *ConfigManager {
	cm.persister = p
	return cm
}

func (cm *ConfigManager) AddEntry(key string, description string, defaultValue any) error {
	// Ancestor check: Verify that no ancestor of the new key is a leaf
	parts := strings.Split(key, ".")
	for i := 1; i < len(parts); i++ {
		subKey := strings.Join(parts[:i], ".")
		if _, exists := cm.entries[subKey]; exists {
			return fmt.Errorf("ancestor %q is a leaf; cannot register %q",
				subKey, key)
		}
	}

	// Descendant check: Verify that no descendant of the new key exists
	prefix := key + "."
	for k := range cm.entries {
		if strings.HasPrefix(k, prefix) {
			return fmt.Errorf("descendant %q exists; cannot register leaf %q", k, key)
		}
	}

	_, exists := cm.entries[key]
	if exists {
		return fmt.Errorf("entry with key %q already exists", key)
	}

	if cm.persister != nil {
		if value, ok := cm.persister.Get(key); ok {
			cm.entries[key] = Entry{
				Description:  description,
				Value:        value,
				DefaultValue: defaultValue,
			}
			return nil
		}
	}

	cm.entries[key] = Entry{
		Description:  description,
		Value:        nil,
		DefaultValue: defaultValue,
	}
	return nil
}

func (cm *ConfigManager) Get(key string) (any, bool) {
	entry, exists := cm.entries[key]
	if !exists {
		return nil, false
	}
	return entry.effective(), true
}

func (cm *ConfigManager) Set(key string, value any) bool {
	entry, exists := cm.entries[key]
	if !exists {
		return false
	}
	entry.Value = value
	cm.entries[key] = entry
	return true
}

func (cm *ConfigManager) Has(key string) bool {
	_, exists := cm.entries[key]
	return exists
}

func (cm *ConfigManager) Save() error {
	if cm.persister == nil {
		return fmt.Errorf("no persister configured")
	}

	// snapshot := map[string]any{}
	// for key, entry := range cm.entries {
	// 	snapshot[key] = entry.Value
	// }

	return cm.persister.Save(cm.entries)
}

func (cm *ConfigManager) Load() error {
	if cm.persister == nil {
		return fmt.Errorf("no persister configured")
	}

	return cm.persister.Load()
}
