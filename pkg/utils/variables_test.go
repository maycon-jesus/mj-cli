package utils

import "testing"

func TestReplaceVariables(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		variables map[string]string
		want      string
	}{
		{
			name:      "uma variável",
			input:     "olá {{nome}}",
			variables: map[string]string{"nome": "Ana"},
			want:      "olá Ana",
		},
		{
			name:      "múltiplas variáveis",
			input:     "{{nome}} tem {{count}} mensagens",
			variables: map[string]string{"nome": "Ana", "count": "3"},
			want:      "Ana tem 3 mensagens",
		},
		{
			name:      "mesma variável repetida",
			input:     "{{x}} e {{x}}",
			variables: map[string]string{"x": "ok"},
			want:      "ok e ok",
		},
		{
			name:      "variável não presente no mapa",
			input:     "olá {{nome}}, porta {{port}}",
			variables: map[string]string{"nome": "Ana"},
			want:      "olá Ana, porta {{port}}",
		},
		{
			name:      "mapa nil",
			input:     "olá {{nome}}",
			variables: nil,
			want:      "olá {{nome}}",
		},
		{
			name:      "mapa vazio",
			input:     "olá {{nome}}",
			variables: map[string]string{},
			want:      "olá {{nome}}",
		},
		{
			name:      "string sem placeholders",
			input:     "sem variáveis aqui",
			variables: map[string]string{"nome": "Ana"},
			want:      "sem variáveis aqui",
		},
		{
			name:      "string vazia",
			input:     "",
			variables: map[string]string{"nome": "Ana"},
			want:      "",
		},
		{
			name:      "placeholder adjacentes",
			input:     "{{a}}{{b}}",
			variables: map[string]string{"a": "X", "b": "Y"},
			want:      "XY",
		},
		{
			name:      "variável com underscore e dígitos",
			input:     "{{var_1}} e {{_ok2}}",
			variables: map[string]string{"var_1": "A", "_ok2": "B"},
			want:      "A e B",
		},
		{
			name:      "placeholder com espaços interno não é substituído",
			input:     "{{ nome }}",
			variables: map[string]string{"nome": "Ana"},
			want:      "{{ nome }}",
		},
		{
			name:      "placeholder vazio não é substituído",
			input:     "{{}}",
			variables: map[string]string{"": "valor"},
			want:      "{{}}",
		},
		{
			name:      "mustache simples não é substituído",
			input:     "{nome}",
			variables: map[string]string{"nome": "Ana"},
			want:      "{nome}",
		},
		{
			name:      "valor da variável é string vazia",
			input:     "olá {{nome}}!",
			variables: map[string]string{"nome": ""},
			want:      "olá !",
		},
		{
			name:      "valor contém mustache literal",
			input:     "{{a}}",
			variables: map[string]string{"a": "{{b}}"},
			want:      "{{b}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReplaceVariables(tt.input, tt.variables)
			if got != tt.want {
				t.Errorf("ReplaceVariables() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReplaceVariablesWithFn(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		replaceFn ReplaceFn
		want      string
	}{
		{
			name:  "substituição básica",
			input: "olá {{nome}}",
			replaceFn: func(v string) string {
				return "Ana"
			},
			want: "olá Ana",
		},
		{
			name:  "múltiplas variáveis resolvidas pela fn",
			input: "{{a}} e {{b}}",
			replaceFn: func(v string) string {
				m := map[string]string{"a": "X", "b": "Y"}
				return m[v]
			},
			want: "X e Y",
		},
		{
			name:  "fn retorna vazio — placeholder mantido",
			input: "olá {{nome}}",
			replaceFn: func(v string) string {
				return ""
			},
			want: "olá {{nome}}",
		},
		{
			name:  "fn retorna vazio para alguma variável apenas",
			input: "{{a}} e {{b}}",
			replaceFn: func(v string) string {
				if v == "a" {
					return "X"
				}
				return ""
			},
			want: "X e {{b}}",
		},
		{
			name:      "fn nil — string sem alteração",
			input:     "olá {{nome}}",
			replaceFn: nil,
			want:      "olá {{nome}}",
		},
		{
			name:      "sem placeholders",
			input:     "sem variáveis",
			replaceFn: func(v string) string { return "ignorado" },
			want:      "sem variáveis",
		},
		{
			name:      "string vazia",
			input:     "",
			replaceFn: func(v string) string { return "ignorado" },
			want:      "",
		},
		{
			name:  "fn recebe o nome correto da variável",
			input: "{{porta}}",
			replaceFn: func(v string) string {
				if v != "porta" {
					t.Errorf("ReplaceFn recebeu variável %q, esperava %q", v, "porta")
				}
				return "8080"
			},
			want: "8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReplaceVariablesWithFn(tt.input, tt.replaceFn)
			if got != tt.want {
				t.Errorf("ReplaceVariablesWithFn() = %q, want %q", got, tt.want)
			}
		})
	}
}
