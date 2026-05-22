---
name: package-readme
description: Cria ou atualiza o README.md de um pacote Go do projeto a partir da análise do código-fonte daquele pacote. Use quando o usuário pedir para gerar/criar/atualizar/documentar o README de um pacote específico (ex.: "crie o README do pkg/cube", "atualize a documentação do pacote config").
---

# package-readme

Cria ou atualiza o `README.md` de um pacote Go deste repositório, analisando o código-fonte do pacote e seguindo o estilo de documentação já estabelecido no projeto.

## Quando usar

Invoque esta skill quando o usuário pedir para:
- "criar README do pacote X"
- "atualizar a documentação do pkg/X"
- "documentar o pacote X"
- "gerar README para internal/X"

Sempre exija que o usuário indique **qual pacote** documentar. Se ele não disser, pergunte antes de prosseguir — não escolha um pacote por conta própria.

## Passos

1. **Resolver o pacote** — identifique o diretório do pacote a partir do que o usuário pediu.
   - Pacotes públicos ficam em `pkg/` (ex.: `pkg/cube`, `pkg/tint`).
   - Código interno fica em `internal/` (ex.: `internal/config`, `internal/services`).
   - Se houver ambiguidade (nome existe em mais de um lugar, ou não existe), liste as opções e peça confirmação.

2. **Verificar se já existe README** — execute `ls <dir>/README.md`.
   - Se existir: é uma **atualização**. Leia o README atual primeiro para preservar seções ainda válidas e o tom.
   - Se não existir: é uma **criação** do zero.

3. **Analisar o código-fonte** — leia os arquivos `.go` do pacote (ignore `*_test.go` para a estrutura, mas use-os para entender comportamento e exemplos de uso):
   - Identifique os **tipos exportados** (structs, interfaces) e seus campos relevantes.
   - Identifique as **funções/métodos exportados** — assinatura, propósito, parâmetros e retornos.
   - Identifique **constantes, variáveis de pacote e enums** exportados.
   - Note **comportamentos não óbvios**: thread safety, efeitos colaterais, casos de borda, variáveis de ambiente lidas.
   - Veja como o pacote é **usado pelo resto do projeto** (`grep` pelo nome do pacote em `cmd/`, `internal/`, `main.go`) para extrair exemplos reais.
   - Confira a **doc do pacote** (comentário acima de `package X`) se existir.

4. **Escrever o README** seguindo a estrutura abaixo.

5. **Confirmar antes de criar/sobrescrever** — em atualização, mostre o que mudou ou descreva as alterações. Não invente APIs: documente apenas o que existe no código.

## Estrutura do README

Espelhe o estilo dos READMEs existentes do projeto (`pkg/logger/README.md`, `pkg/utils/README.md`, `pkg/cmd/README.md`). Adapte as seções ao tamanho do pacote — pacotes pequenos não precisam de todas.

```markdown
# <nome-do-pacote>

<Parágrafo curto: o que o pacote faz e qual problema resolve. Se for
construído sobre algo da stdlib ou dependência, mencione com link.>

## Uso rápido

<Bloco de código Go mostrando o caso de uso mais comum, ponta a ponta.>

## Componentes

### <Tipo ou grupo de funções>

<Descrição. Bloco de código com a definição do tipo quando ajudar.
Tabela de métodos quando houver muitos:>

| Método | Descrição |
|--------|-----------|
| `Foo()` | ... |

## <Seções específicas do pacote>

<Comportamentos, variáveis de ambiente, casos de borda — use tabelas
para mapear situação → resultado quando couber.>

## Thread Safety

<Só inclua se relevante: o que é protegido e como.>

## Testes

​```bash
go test ./<caminho-do-pacote>/ -v
​```

<O que os testes cobrem, em bullets.>
```

## Regras de estilo

- **Idioma**: português do Brasil (pt-BR). Termos técnicos consagrados (handler, thread-safe, struct, fan-out) podem ficar em inglês.
- Documente **apenas a API exportada** e o comportamento real do código. Não documente funções não exportadas, a menos que sejam essenciais para entender o pacote.
- Exemplos de código devem ser **válidos e realistas** — prefira extrair de usos reais no projeto a inventar.
- Use **tabelas** para mapear método→descrição e situação→resultado; use **blocos de código** para definições de tipo e exemplos.
- Foque no **porquê e no como usar**, não em repetir a assinatura linha a linha.
- Mantenha consistência com os READMEs já existentes (cabeçalhos, tom, profundidade).

## Anti-padrões

- ❌ Documentar APIs que não existem no código ou inventar comportamento.
- ❌ Copiar a assinatura de cada função sem explicar uso ou propósito.
- ❌ README genérico que serviria para qualquer pacote.
- ❌ Sobrescrever um README existente sem antes lê-lo e preservar o que ainda é válido.
- ❌ Escrever em inglês quando o repositório usa pt-BR.
