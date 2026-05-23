package utils

import (
	"reflect"
	"testing"
)

func TestCreateDeepMap(t *testing.T) {
	tests := []struct {
		name  string
		keys  []string
		value interface{}
		want  map[string]interface{}
	}{
		{
			name:  "lista de chaves vazia",
			keys:  []string{},
			value: "x",
			want:  nil,
		},
		{
			name:  "chaves nil",
			keys:  nil,
			value: "x",
			want:  nil,
		},
		{
			name:  "uma chave",
			keys:  []string{"a"},
			value: "x",
			want:  map[string]interface{}{"a": "x"},
		},
		{
			name:  "duas chaves",
			keys:  []string{"a", "b"},
			value: 42,
			want: map[string]interface{}{
				"a": map[string]interface{}{"b": 42},
			},
		},
		{
			name:  "três chaves aninhadas",
			keys:  []string{"a", "b", "c"},
			value: true,
			want: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{"c": true},
				},
			},
		},
		{
			name:  "valor nil",
			keys:  []string{"a"},
			value: nil,
			want:  map[string]interface{}{"a": nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateDeepMap(tt.keys, tt.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDeepMap(%v, %v) = %v, want %v", tt.keys, tt.value, got, tt.want)
			}
		})
	}
}

func TestMergeDeepMaps(t *testing.T) {
	tests := []struct {
		name string
		dest map[string]interface{}
		src  map[string]interface{}
		want map[string]interface{}
	}{
		{
			name: "chaves disjuntas",
			dest: map[string]interface{}{"a": 1},
			src:  map[string]interface{}{"b": 2},
			want: map[string]interface{}{"a": 1, "b": 2},
		},
		{
			name: "src sobrescreve valor escalar",
			dest: map[string]interface{}{"a": 1},
			src:  map[string]interface{}{"a": 2},
			want: map[string]interface{}{"a": 2},
		},
		{
			name: "merge recursivo de submapas",
			dest: map[string]interface{}{
				"a": map[string]interface{}{"x": 1},
			},
			src: map[string]interface{}{
				"a": map[string]interface{}{"y": 2},
			},
			want: map[string]interface{}{
				"a": map[string]interface{}{"x": 1, "y": 2},
			},
		},
		{
			name: "src mapa sobrescreve dest escalar",
			dest: map[string]interface{}{"a": 1},
			src: map[string]interface{}{
				"a": map[string]interface{}{"x": 1},
			},
			want: map[string]interface{}{
				"a": map[string]interface{}{"x": 1},
			},
		},
		{
			name: "src escalar sobrescreve dest mapa",
			dest: map[string]interface{}{
				"a": map[string]interface{}{"x": 1},
			},
			src:  map[string]interface{}{"a": 5},
			want: map[string]interface{}{"a": 5},
		},
		{
			name: "src vazio mantém dest",
			dest: map[string]interface{}{"a": 1},
			src:  map[string]interface{}{},
			want: map[string]interface{}{"a": 1},
		},
		{
			name: "dest vazio recebe src",
			dest: map[string]interface{}{},
			src:  map[string]interface{}{"a": 1},
			want: map[string]interface{}{"a": 1},
		},
		{
			name: "merge profundo de múltiplos níveis",
			dest: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{"c": 1},
				},
			},
			src: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{"d": 2},
				},
			},
			want: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{"c": 1, "d": 2},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeDeepMaps(tt.dest, tt.src)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeDeepMaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeDeepMapsModificaDest(t *testing.T) {
	dest := map[string]interface{}{"a": 1}
	src := map[string]interface{}{"b": 2}

	got := MergeDeepMaps(dest, src)
	if !reflect.DeepEqual(got, dest) {
		t.Errorf("MergeDeepMaps deve retornar o próprio dest; got %v, dest %v", got, dest)
	}
	if _, ok := dest["b"]; !ok {
		t.Error("MergeDeepMaps deve mutar dest adicionando chaves de src")
	}
}

func TestCreateDeepMapComMergeDeepMaps(t *testing.T) {
	a := CreateDeepMap([]string{"server", "host"}, "localhost")
	b := CreateDeepMap([]string{"server", "port"}, 8080)

	got := MergeDeepMaps(a, b)
	want := map[string]interface{}{
		"server": map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergeDeepMaps(CreateDeepMap...) = %v, want %v", got, want)
	}
}
