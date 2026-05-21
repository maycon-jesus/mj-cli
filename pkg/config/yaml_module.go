package config

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"
)

type YamlModule struct {
	filePath string
	data     map[string]any
}

func NewYamlModule(filePath string) *YamlModule {
	return &YamlModule{
		filePath: filePath,
		data:     map[string]any{},
	}
}

func (m *YamlModule) Save(configs map[string]Entry) error {
	keys := make([]string, 0, len(configs))
	for k := range configs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	root := &yaml.Node{Kind: yaml.MappingNode}
	for _, key := range keys {
		entry := configs[key]
		if err := insertEntry(root, strings.Split(key, "."), entry); err != nil {
			return fmt.Errorf("encode %q: %w", key, err)
		}
	}

	bs, err := yaml.Marshal(root)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	if err := os.WriteFile(m.filePath, bs, 0644); err != nil {
		return fmt.Errorf("write yaml: %w", err)
	}
	return nil
}

func insertEntry(parent *yaml.Node, parts []string, entry Entry) error {
	if len(parts) == 1 {
		keyNode := &yaml.Node{
			Kind:        yaml.ScalarNode,
			Value:       parts[0],
			HeadComment: entry.Description,
		}
		valueNode := &yaml.Node{}
		if err := valueNode.Encode(entry.effective()); err != nil {
			return err
		}
		parent.Content = append(parent.Content, keyNode, valueNode)
		return nil
	}

	for i := 0; i < len(parent.Content); i += 2 {
		if parent.Content[i].Value == parts[0] {
			return insertEntry(parent.Content[i+1], parts[1:], entry)
		}
	}

	// not has current key, create new node
	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: parts[0]}
	childNode := &yaml.Node{Kind: yaml.MappingNode}
	parent.Content = append(parent.Content, keyNode, childNode)
	return insertEntry(childNode, parts[1:], entry)
}

func (m *YamlModule) Load() error {
	bs, err := os.ReadFile(m.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			m.data = map[string]any{}
			return nil
		}
		return fmt.Errorf("read yaml: %w", err)
	}

	nested := map[string]any{}
	if err := yaml.Unmarshal(bs, &nested); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	m.data = nested
	return nil
}

func (m *YamlModule) Get(key string) (any, bool) {
	parts := strings.Split(key, ".")
	var current any = m.data
	for _, part := range parts {
		node, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = node[part]
		if !ok {
			return nil, false
		}
	}
	return current, true
}
