package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewYamlModule(t *testing.T) {
	m := NewYamlModule("/tmp/foo.yaml")
	if m == nil {
		t.Fatal("NewYamlModule returned nil")
	}
	if m.filePath != "/tmp/foo.yaml" {
		t.Errorf("filePath = %q, want '/tmp/foo.yaml'", m.filePath)
	}
	if m.data == nil {
		t.Error("data map should be initialized")
	}
	if len(m.data) != 0 {
		t.Errorf("data should be empty initially; got %d entries", len(m.data))
	}
}

func TestYamlModule_Save(t *testing.T) {
	t.Run("escreve YAML com chaves planas", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.yaml")
		m := NewYamlModule(path)

		configs := map[string]Entry{
			"app_name": {Description: "App name", Value: "MyApp", DefaultValue: "Default"},
			"debug":    {Description: "Debug mode", Value: true, DefaultValue: false},
		}
		if err := m.Save(configs); err != nil {
			t.Fatalf("Save returned error: %v", err)
		}

		content := readFile(t, path)
		if !strings.Contains(content, "app_name: MyApp") {
			t.Errorf("content does not contain 'app_name: MyApp':\n%s", content)
		}
		if !strings.Contains(content, "debug: true") {
			t.Errorf("content does not contain 'debug: true':\n%s", content)
		}
	})

	t.Run("aninha chaves com .", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.yaml")
		m := NewYamlModule(path)

		configs := map[string]Entry{
			"git.branch":          {Value: "main", DefaultValue: "main"},
			"git.signing.enabled": {Value: true, DefaultValue: false},
			"git.signing.key":     {Value: "ABC123", DefaultValue: ""},
		}
		if err := m.Save(configs); err != nil {
			t.Fatalf("Save returned error: %v", err)
		}

		content := readFile(t, path)
		// reload via Load + Get to validate structure
		loaded := NewYamlModule(path)
		if err := loaded.Load(); err != nil {
			t.Fatalf("Load returned error: %v", err)
		}

		if v, ok := loaded.Get("git.branch"); !ok || v != "main" {
			t.Errorf("Get(git.branch) = %v, %v; want 'main', true. content:\n%s", v, ok, content)
		}
		if v, ok := loaded.Get("git.signing.enabled"); !ok || v != true {
			t.Errorf("Get(git.signing.enabled) = %v, %v; want true, true", v, ok)
		}
		if v, ok := loaded.Get("git.signing.key"); !ok || v != "ABC123" {
			t.Errorf("Get(git.signing.key) = %v, %v; want 'ABC123', true", v, ok)
		}
	})

	t.Run("Description vira head comment", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.yaml")
		m := NewYamlModule(path)

		configs := map[string]Entry{
			"lang": {Description: "Idioma da UI", Value: "en", DefaultValue: "pt-BR"},
		}
		if err := m.Save(configs); err != nil {
			t.Fatalf("Save returned error: %v", err)
		}

		content := readFile(t, path)
		if !strings.Contains(content, "# Idioma da UI") {
			t.Errorf("content should contain head comment '# Idioma da UI':\n%s", content)
		}
	})

	t.Run("usa DefaultValue quando Value é nil", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.yaml")
		m := NewYamlModule(path)

		configs := map[string]Entry{
			"k": {Value: nil, DefaultValue: "fallback"},
		}
		if err := m.Save(configs); err != nil {
			t.Fatalf("Save returned error: %v", err)
		}

		loaded := NewYamlModule(path)
		if err := loaded.Load(); err != nil {
			t.Fatalf("Load returned error: %v", err)
		}
		v, ok := loaded.Get("k")
		if !ok || v != "fallback" {
			t.Errorf("Get(k) = %v, %v; want 'fallback', true", v, ok)
		}
	})

	t.Run("saída é determinística (ordem alfabética)", func(t *testing.T) {
		path1 := filepath.Join(t.TempDir(), "config.yaml")
		path2 := filepath.Join(t.TempDir(), "config.yaml")

		configs := map[string]Entry{
			"z": {Value: 1},
			"a": {Value: 2},
			"m": {Value: 3},
		}

		// Save twice in independent modules; bytes must match.
		if err := NewYamlModule(path1).Save(configs); err != nil {
			t.Fatalf("Save 1 returned error: %v", err)
		}
		if err := NewYamlModule(path2).Save(configs); err != nil {
			t.Fatalf("Save 2 returned error: %v", err)
		}
		if readFile(t, path1) != readFile(t, path2) {
			t.Error("Save output is not deterministic between calls")
		}

		content := readFile(t, path1)
		// Verify alphabetical: 'a' before 'm' before 'z'.
		iA := strings.Index(content, "a:")
		iM := strings.Index(content, "m:")
		iZ := strings.Index(content, "z:")
		if !(iA >= 0 && iM > iA && iZ > iM) {
			t.Errorf("keys not in alphabetical order:\n%s", content)
		}
	})

	t.Run("erro ao gravar em diretório inexistente", func(t *testing.T) {
		m := NewYamlModule("/nonexistent-dir-xyz/config.yaml")
		err := m.Save(map[string]Entry{"k": {Value: "v"}})
		if err == nil {
			t.Fatal("expected error writing to invalid path")
		}
	})

	t.Run("mapa de configs vazio gera arquivo carregável", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "empty.yaml")
		m := NewYamlModule(path)

		if err := m.Save(map[string]Entry{}); err != nil {
			t.Fatalf("Save com mapa vazio retornou erro: %v", err)
		}

		loaded := NewYamlModule(path)
		if err := loaded.Load(); err != nil {
			t.Fatalf("Load do arquivo vazio retornou erro: %v", err)
		}
		if _, ok := loaded.Get("anything"); ok {
			t.Error("Get em arquivo vazio não deveria retornar nada")
		}
	})
}

func TestYamlModule_Load(t *testing.T) {
	t.Run("arquivo inexistente não é erro", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "does-not-exist.yaml")
		m := NewYamlModule(path)
		if err := m.Load(); err != nil {
			t.Fatalf("Load on missing file returned error: %v", err)
		}
		if len(m.data) != 0 {
			t.Errorf("data should be empty for missing file; got %v", m.data)
		}
	})

	t.Run("carrega YAML aninhado", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.yaml")
		yamlContent := `general:
  lang: en
git:
  branch: develop
  signing:
    enabled: true
`
		if err := os.WriteFile(path, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("failed to write fixture: %v", err)
		}

		m := NewYamlModule(path)
		if err := m.Load(); err != nil {
			t.Fatalf("Load returned error: %v", err)
		}

		if v, ok := m.Get("general.lang"); !ok || v != "en" {
			t.Errorf("Get(general.lang) = %v, %v; want 'en', true", v, ok)
		}
		if v, ok := m.Get("git.branch"); !ok || v != "develop" {
			t.Errorf("Get(git.branch) = %v, %v; want 'develop', true", v, ok)
		}
		if v, ok := m.Get("git.signing.enabled"); !ok || v != true {
			t.Errorf("Get(git.signing.enabled) = %v, %v; want true, true", v, ok)
		}
	})

	t.Run("erro em YAML inválido", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "broken.yaml")
		if err := os.WriteFile(path, []byte("not: valid: yaml: ::: ["), 0644); err != nil {
			t.Fatalf("failed to write fixture: %v", err)
		}

		m := NewYamlModule(path)
		if err := m.Load(); err == nil {
			t.Fatal("expected error for invalid YAML, got nil")
		}
	})

	t.Run("erro em outras falhas de leitura", func(t *testing.T) {
		// Pass a directory as the file path to trigger a non-IsNotExist read error.
		dir := t.TempDir()
		m := NewYamlModule(dir)
		if err := m.Load(); err == nil {
			t.Fatal("expected error reading a directory as a file, got nil")
		}
	})

	t.Run("arquivo vazio (0 bytes) não é erro", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "blank.yaml")
		if err := os.WriteFile(path, []byte(""), 0644); err != nil {
			t.Fatalf("failed to write fixture: %v", err)
		}

		m := NewYamlModule(path)
		if err := m.Load(); err != nil {
			t.Fatalf("Load de arquivo vazio retornou erro: %v", err)
		}
		// Get em qualquer chave deve retornar (nil, false) sem panic.
		if v, ok := m.Get("foo"); ok || v != nil {
			t.Errorf("Get em arquivo vazio = (%v, %v); want (nil, false)", v, ok)
		}
		if v, ok := m.Get("a.b.c"); ok || v != nil {
			t.Errorf("Get aninhado em arquivo vazio = (%v, %v); want (nil, false)", v, ok)
		}
	})
}

func TestYamlModule_Get(t *testing.T) {
	t.Run("chave inexistente", func(t *testing.T) {
		m := NewYamlModule("/tmp/ignored.yaml")
		v, ok := m.Get("missing")
		if ok {
			t.Error("ok should be false for missing key")
		}
		if v != nil {
			t.Errorf("value = %v, want nil", v)
		}
	})

	t.Run("caminho aninhado interrompido por escalar", func(t *testing.T) {
		m := NewYamlModule("/tmp/ignored.yaml")
		m.data = map[string]any{
			"git": "not-a-map",
		}
		v, ok := m.Get("git.branch")
		if ok {
			t.Error("ok should be false when traversing through a scalar")
		}
		if v != nil {
			t.Errorf("value = %v, want nil", v)
		}
	})

	t.Run("caminho aninhado completo", func(t *testing.T) {
		m := NewYamlModule("/tmp/ignored.yaml")
		m.data = map[string]any{
			"git": map[string]any{
				"signing": map[string]any{
					"enabled": true,
				},
			},
		}
		v, ok := m.Get("git.signing.enabled")
		if !ok {
			t.Fatal("expected key to be found")
		}
		if v != true {
			t.Errorf("value = %v, want true", v)
		}
	})

	t.Run("chave de nível intermediário retorna sub-mapa", func(t *testing.T) {
		m := NewYamlModule("/tmp/ignored.yaml")
		m.data = map[string]any{
			"git": map[string]any{
				"branch": "main",
			},
		}
		v, ok := m.Get("git")
		if !ok {
			t.Fatal("expected key 'git' to be found")
		}
		if _, isMap := v.(map[string]any); !isMap {
			t.Errorf("value = %v, want a map[string]any", v)
		}
	})
}

func TestYamlModule_SaveAndLoadRoundtrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")

	original := map[string]Entry{
		"app_name":            {Description: "App name", Value: "MyApp", DefaultValue: "Default"},
		"debug":               {Description: "Debug mode", Value: true, DefaultValue: false},
		"general.lang":        {Description: "Idioma", Value: "en", DefaultValue: "pt-BR"},
		"git.branch":          {Value: "main", DefaultValue: "main"},
		"git.signing.enabled": {Value: false, DefaultValue: false},
		"max_retries":         {Value: 3, DefaultValue: 1},
	}

	persister := NewYamlModule(path)
	if err := persister.Save(original); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	loaded := NewYamlModule(path)
	if err := loaded.Load(); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	// Spot-check a few values via Get; types should roundtrip via YAML.
	cases := []struct {
		key  string
		want any
	}{
		{"app_name", "MyApp"},
		{"debug", true},
		{"general.lang", "en"},
		{"git.branch", "main"},
		{"git.signing.enabled", false},
		{"max_retries", 3},
	}
	for _, c := range cases {
		v, ok := loaded.Get(c.key)
		if !ok {
			t.Errorf("Get(%q) not found", c.key)
			continue
		}
		if v != c.want {
			t.Errorf("Get(%q) = %v (%T), want %v (%T)", c.key, v, v, c.want, c.want)
		}
	}
}

func TestYamlModule_RoundtripVariedTypes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "types.yaml")

	original := map[string]Entry{
		"float":    {Value: 3.14, DefaultValue: 0.0},
		"strings":  {Value: []string{"a", "b", "c"}, DefaultValue: []string{}},
		"ints":     {Value: []int{1, 2, 3}, DefaultValue: []int{}},
		"nullable": {Value: nil, DefaultValue: "default-used"},
	}

	if err := NewYamlModule(path).Save(original); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	loaded := NewYamlModule(path)
	if err := loaded.Load(); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	// float64 roundtrip
	if v, ok := loaded.Get("float"); !ok {
		t.Error("Get(float) not found")
	} else if f, isFloat := v.(float64); !isFloat || f != 3.14 {
		t.Errorf("Get(float) = %v (%T), want 3.14 (float64)", v, v)
	}

	// slices voltam como []any (YAML genérico)
	if v, ok := loaded.Get("strings"); !ok {
		t.Error("Get(strings) not found")
	} else if s, isSlice := v.([]any); !isSlice {
		t.Errorf("Get(strings) = %v (%T), want []any", v, v)
	} else if len(s) != 3 || s[0] != "a" || s[1] != "b" || s[2] != "c" {
		t.Errorf("Get(strings) = %v, want [a b c]", s)
	}

	if v, ok := loaded.Get("ints"); !ok {
		t.Error("Get(ints) not found")
	} else if s, isSlice := v.([]any); !isSlice {
		t.Errorf("Get(ints) = %v (%T), want []any", v, v)
	} else if len(s) != 3 {
		t.Errorf("Get(ints) len = %d, want 3", len(s))
	}

	// nullable: Value nil + DefaultValue "default-used" → o arquivo grava "default-used"
	if v, ok := loaded.Get("nullable"); !ok || v != "default-used" {
		t.Errorf("Get(nullable) = (%v, %v); want ('default-used', true)", v, ok)
	}
}

func TestYamlModule_ImplementsPersistModule(t *testing.T) {
	var _ PersistModule = (*YamlModule)(nil)
}

// helpers

func readFile(t *testing.T, path string) string {
	t.Helper()
	bs, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %q: %v", path, err)
	}
	return string(bs)
}
