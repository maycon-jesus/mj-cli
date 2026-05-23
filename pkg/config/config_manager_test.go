package config

import (
	"errors"
	"reflect"
	"testing"
)

// fakePersister is a stub PersistModule used to drive ConfigManager tests
// without touching the filesystem.
type fakePersister struct {
	stored     map[string]any
	saveCalls  int
	loadCalls  int
	lastSaved  map[string]Entry
	saveErr    error
	loadErr    error
}

func newFakePersister() *fakePersister {
	return &fakePersister{stored: map[string]any{}}
}

func (f *fakePersister) Save(configs map[string]Entry) error {
	f.saveCalls++
	f.lastSaved = configs
	return f.saveErr
}

func (f *fakePersister) Load() error {
	f.loadCalls++
	return f.loadErr
}

func (f *fakePersister) Get(key string) (any, bool) {
	v, ok := f.stored[key]
	return v, ok
}

func TestEntry_Effective(t *testing.T) {
	tests := []struct {
		name string
		e    Entry
		want any
	}{
		{
			name: "value não-nulo tem prioridade",
			e:    Entry{Value: "custom", DefaultValue: "fallback"},
			want: "custom",
		},
		{
			name: "value nulo cai no default",
			e:    Entry{Value: nil, DefaultValue: "fallback"},
			want: "fallback",
		},
		{
			name: "value false não é confundido com nil",
			e:    Entry{Value: false, DefaultValue: true},
			want: false,
		},
		{
			name: "zero int (0) não é confundido com nil",
			e:    Entry{Value: 0, DefaultValue: 42},
			want: 0,
		},
		{
			name: "default e value ambos nil",
			e:    Entry{Value: nil, DefaultValue: nil},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.effective(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("effective() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConfigManager(t *testing.T) {
	cm := NewConfigManager()
	if cm == nil {
		t.Fatal("NewConfigManager() returned nil")
	}
	if cm.entries == nil {
		t.Error("entries map should be initialized")
	}
	if len(cm.entries) != 0 {
		t.Errorf("entries should be empty initially; got %d", len(cm.entries))
	}
	if cm.persister != nil {
		t.Error("persister should be nil initially")
	}
}

func TestConfigManager_WithPersister(t *testing.T) {
	cm := NewConfigManager()
	p := newFakePersister()

	got := cm.WithPersister(p)
	if got != cm {
		t.Error("WithPersister should return the same manager for chaining")
	}
	if cm.persister != p {
		t.Error("persister was not assigned")
	}
}

func TestConfigManager_AddEntry(t *testing.T) {
	t.Run("registra entrada nova", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("foo", "a foo", "bar"); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}
		got, ok := cm.Get("foo")
		if !ok {
			t.Fatal("entry not found after AddEntry")
		}
		if got != "bar" {
			t.Errorf("Get = %v, want 'bar'", got)
		}
	})

	t.Run("erro ao registrar chave duplicada", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("foo", "", "x"); err != nil {
			t.Fatalf("first AddEntry returned error: %v", err)
		}
		err := cm.AddEntry("foo", "", "y")
		if err == nil {
			t.Fatal("expected error for duplicate key, got nil")
		}
	})

	t.Run("erro quando ancestral é folha", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("git", "", "leaf"); err != nil {
			t.Fatalf("AddEntry git returned error: %v", err)
		}
		err := cm.AddEntry("git.branch", "", "main")
		if err == nil {
			t.Fatal("expected error when ancestor is a leaf, got nil")
		}
	})

	t.Run("erro quando descendente já existe", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("git.branch", "", "main"); err != nil {
			t.Fatalf("AddEntry git.branch returned error: %v", err)
		}
		err := cm.AddEntry("git", "", "x")
		if err == nil {
			t.Fatal("expected error when descendant exists, got nil")
		}
	})

	t.Run("erro com folha em nível intermediário (descendente em sub-árvore)", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("git.signing.enabled", "", true); err != nil {
			t.Fatalf("AddEntry git.signing.enabled returned error: %v", err)
		}
		// git.signing tem git.signing.enabled como descendente; não pode virar folha.
		err := cm.AddEntry("git.signing", "", "x")
		if err == nil {
			t.Fatal("expected error when intermediate key has descendants, got nil")
		}
	})

	t.Run("erro com ancestral em nível intermediário", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("git.signing", "", "leaf"); err != nil {
			t.Fatalf("AddEntry git.signing returned error: %v", err)
		}
		// git.signing já é folha; git.signing.enabled deve falhar.
		err := cm.AddEntry("git.signing.enabled", "", true)
		if err == nil {
			t.Fatal("expected error when intermediate ancestor is a leaf, got nil")
		}
	})

	t.Run("chaves irmãs aninhadas convivem", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("git.branch", "", "main"); err != nil {
			t.Fatalf("AddEntry git.branch returned error: %v", err)
		}
		if err := cm.AddEntry("git.remote", "", "origin"); err != nil {
			t.Fatalf("AddEntry git.remote returned error: %v", err)
		}
		if err := cm.AddEntry("git.signing.enabled", "", false); err != nil {
			t.Fatalf("AddEntry git.signing.enabled returned error: %v", err)
		}
	})

	t.Run("carrega value do persister quando disponível", func(t *testing.T) {
		p := newFakePersister()
		p.stored["general.lang"] = "en"
		cm := NewConfigManager().WithPersister(p)

		if err := cm.AddEntry("general.lang", "Idioma", "pt-BR"); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}

		got, ok := cm.Get("general.lang")
		if !ok {
			t.Fatal("entry not found")
		}
		if got != "en" {
			t.Errorf("Get = %v, want 'en' (loaded from persister)", got)
		}
	})

	t.Run("usa default quando persister não tem a chave", func(t *testing.T) {
		p := newFakePersister()
		cm := NewConfigManager().WithPersister(p)

		if err := cm.AddEntry("general.lang", "Idioma", "pt-BR"); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}

		got, ok := cm.Get("general.lang")
		if !ok {
			t.Fatal("entry not found")
		}
		if got != "pt-BR" {
			t.Errorf("Get = %v, want 'pt-BR' (default)", got)
		}
	})

	t.Run("persister retorna (nil, true): Value=nil cai no DefaultValue", func(t *testing.T) {
		p := newFakePersister()
		p.stored["k"] = nil // chave presente no backend, mas com valor nulo
		cm := NewConfigManager().WithPersister(p)

		if err := cm.AddEntry("k", "", "fallback"); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}

		// O Value armazenado é nil; effective() devolve o DefaultValue.
		entry := cm.entries["k"]
		if entry.Value != nil {
			t.Errorf("entry.Value = %v, want nil (persister forneceu nil)", entry.Value)
		}
		got, _ := cm.Get("k")
		if got != "fallback" {
			t.Errorf("Get = %v, want 'fallback' (default usado quando Value é nil)", got)
		}
	})
}

func TestConfigManager_Get(t *testing.T) {
	t.Run("chave inexistente retorna false", func(t *testing.T) {
		cm := NewConfigManager()
		v, ok := cm.Get("missing")
		if ok {
			t.Error("ok should be false for missing key")
		}
		if v != nil {
			t.Errorf("value = %v, want nil", v)
		}
	})

	t.Run("retorna default quando value é nil", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("k", "", "default-val"); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}
		v, ok := cm.Get("k")
		if !ok {
			t.Fatal("entry should exist")
		}
		if v != "default-val" {
			t.Errorf("value = %v, want 'default-val'", v)
		}
	})

	t.Run("retorna value após Set", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("k", "", "default-val"); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}
		cm.Set("k", "updated")
		v, _ := cm.Get("k")
		if v != "updated" {
			t.Errorf("value = %v, want 'updated'", v)
		}
	})
}

func TestConfigManager_Set(t *testing.T) {
	t.Run("retorna false para chave inexistente", func(t *testing.T) {
		cm := NewConfigManager()
		if ok := cm.Set("missing", 123); ok {
			t.Error("Set should return false for missing key")
		}
	})

	t.Run("atualiza chave existente", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("k", "", "old"); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}
		if ok := cm.Set("k", "new"); !ok {
			t.Error("Set should return true for existing key")
		}
		v, _ := cm.Get("k")
		if v != "new" {
			t.Errorf("value = %v, want 'new'", v)
		}
	})

	t.Run("preserva DefaultValue ao atualizar Value", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("k", "desc", "default"); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}
		cm.Set("k", "current")

		entry := cm.entries["k"]
		if entry.DefaultValue != "default" {
			t.Errorf("DefaultValue = %v, want 'default'", entry.DefaultValue)
		}
		if entry.Description != "desc" {
			t.Errorf("Description = %v, want 'desc'", entry.Description)
		}
	})
}

func TestConfigManager_Has(t *testing.T) {
	t.Run("retorna false para chave inexistente", func(t *testing.T) {
		cm := NewConfigManager()
		if cm.Has("missing") {
			t.Error("Has should return false for missing key")
		}
	})

	t.Run("retorna true para chave registrada via AddEntry", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("k", "", "default"); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}
		if !cm.Has("k") {
			t.Error("Has should return true for registered key")
		}
	})

	t.Run("retorna true mesmo quando Value é nil (apenas DefaultValue definido)", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("k", "", "default"); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}
		// Value permanece nil; Has checa existência, não conteúdo.
		if cm.entries["k"].Value != nil {
			t.Fatalf("precondição: Value deveria ser nil, got %v", cm.entries["k"].Value)
		}
		if !cm.Has("k") {
			t.Error("Has should return true even when Value is nil")
		}
	})

	t.Run("retorna true após Set", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("k", "", "default"); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}
		cm.Set("k", "updated")
		if !cm.Has("k") {
			t.Error("Has should return true after Set on existing key")
		}
	})

	t.Run("não confunde chaves aninhadas com ancestrais", func(t *testing.T) {
		cm := NewConfigManager()
		if err := cm.AddEntry("git.branch", "", "main"); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}
		if cm.Has("git") {
			t.Error("Has should return false for ancestor of a registered key")
		}
		if !cm.Has("git.branch") {
			t.Error("Has should return true for the exact registered key")
		}
	})
}

func TestConfigManager_Save(t *testing.T) {
	t.Run("erro sem persister", func(t *testing.T) {
		cm := NewConfigManager()
		err := cm.Save()
		if err == nil {
			t.Fatal("expected error without persister")
		}
	})

	t.Run("delega para persister", func(t *testing.T) {
		p := newFakePersister()
		cm := NewConfigManager().WithPersister(p)
		if err := cm.AddEntry("a", "", 1); err != nil {
			t.Fatalf("AddEntry returned error: %v", err)
		}
		cm.Set("a", 2)

		if err := cm.Save(); err != nil {
			t.Fatalf("Save returned error: %v", err)
		}

		if p.saveCalls != 1 {
			t.Errorf("saveCalls = %d, want 1", p.saveCalls)
		}
		if got := p.lastSaved["a"].Value; got != 2 {
			t.Errorf("lastSaved[a].Value = %v, want 2", got)
		}
	})

	t.Run("propaga erro do persister", func(t *testing.T) {
		p := newFakePersister()
		p.saveErr = errors.New("disk full")
		cm := NewConfigManager().WithPersister(p)

		err := cm.Save()
		if err == nil || !errors.Is(err, p.saveErr) {
			t.Fatalf("Save error = %v, want %v", err, p.saveErr)
		}
	})
}

func TestConfigManager_Load(t *testing.T) {
	t.Run("erro sem persister", func(t *testing.T) {
		cm := NewConfigManager()
		err := cm.Load()
		if err == nil {
			t.Fatal("expected error without persister")
		}
	})

	t.Run("delega para persister", func(t *testing.T) {
		p := newFakePersister()
		cm := NewConfigManager().WithPersister(p)
		if err := cm.Load(); err != nil {
			t.Fatalf("Load returned error: %v", err)
		}
		if p.loadCalls != 1 {
			t.Errorf("loadCalls = %d, want 1", p.loadCalls)
		}
	})

	t.Run("propaga erro do persister", func(t *testing.T) {
		p := newFakePersister()
		p.loadErr = errors.New("parse error")
		cm := NewConfigManager().WithPersister(p)

		err := cm.Load()
		if err == nil || !errors.Is(err, p.loadErr) {
			t.Fatalf("Load error = %v, want %v", err, p.loadErr)
		}
	})
}
