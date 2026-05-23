---
name: commit-message
description: Analisa os arquivos no stage do git (git diff --staged) e gera duas versões de mensagem de commit (uma curta e uma longa completa) seguindo o padrão Conventional Commits em pt-BR. Use quando o usuário pedir para gerar/sugerir/criar uma mensagem de commit a partir do que está staged.
---

# commit-message

Gera **duas versões** de mensagem de commit no padrão **Conventional Commits** em **português do Brasil (pt-BR)** a partir das mudanças que estão no stage do git: uma **curta** (só o título) e uma **longa** (título + corpo + rodapés).

## Quando usar

Invoque esta skill quando o usuário pedir para:
- "gerar mensagem de commit"
- "sugerir commit message"
- "criar mensagem para o que está no stage"
- "fazer commit do que está staged" (gere a mensagem; só rode `git commit` se o usuário pedir explicitamente)

## Passos

1. **Inspecionar o stage** — execute em paralelo:
   - `git status --short` — lista resumida do que está staged/unstaged
   - `git diff --staged` — diff completo dos arquivos no stage
   - `git log -n 10 --oneline` — últimos commits para alinhar com o estilo do repositório

2. **Validar** — se não houver nada staged (`git diff --staged` vazio):
   - Avisar o usuário de que não há arquivos no stage
   - Não inventar uma mensagem; perguntar se deve incluir as mudanças unstaged ou abortar

3. **Analisar o diff** identificando:
   - **Tipo**: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `style`, `perf`, `build`, `ci`, `revert`
   - **Escopo** (opcional): área/módulo afetado (ex.: `cube`, `cmd`, `config`, `i18n`). Use diretórios reais do repositório como referência.
   - **Breaking change**: se houver remoção/alteração incompatível de API pública, marcar com `!` após o tipo/escopo e/ou adicionar rodapé `BREAKING CHANGE:`
   - **Descrição**: foco no **porquê**, não no o que. Imperativo, minúsculo, sem ponto final, ≤ 72 caracteres na primeira linha.

4. **Formato de saída** — apresente **sempre as duas versões** ao usuário antes de commitar:

   **Versão curta** — apenas a linha de título, sem corpo. Para commits rápidos.

   ```
   <tipo>(<escopo opcional>): <descrição curta em pt-BR>
   ```

   **Versão longa** — título + corpo explicando o porquê + rodapés. Para commits que merecem contexto.

   ```
   <tipo>(<escopo opcional>): <descrição curta em pt-BR>

   <corpo explicando o porquê — quebre em 72 colunas>

   <rodapé opcional: BREAKING CHANGE: ... / Refs: #123>
   ```

   Apresente as duas sob títulos claros (ex.: `### Curta` e `### Longa`), cada uma
   em seu próprio bloco de código para facilitar a cópia. A linha de título pode ser
   idêntica nas duas versões; a longa apenas acrescenta corpo/rodapés. Se houver
   `BREAKING CHANGE`, inclua o `!` no tipo/escopo em **ambas** as versões e o rodapé
   na versão longa.

5. **Múltiplas mudanças não relacionadas** — se o stage misturar mudanças que mereceriam commits separados, avise o usuário e sugira dividir (`git reset` + re-staging seletivo). Não force uma única mensagem genérica.

6. **Não commitar automaticamente.** Mostre as duas mensagens propostas e aguarde aprovação. Se o usuário pedir para commitar sem dizer qual versão, pergunte qual usar (curta ou longa). Só execute `git commit` se o usuário disser explicitamente para commitar.

## Regras de estilo (pt-BR)

- Descrição em **letras minúsculas**, verbo no **infinitivo** ou **imperativo curto** (ex.: `adicionar`, `corrigir`, `remover`, `atualizar`, `refatorar`).
- Sem ponto final na primeira linha.
- Evite termos vagos como "ajustes", "melhorias", "mudanças diversas" — seja específico sobre o que mudou e por quê.
- "add" implica funcionalidade nova; "update" implica enhancement; "fix" implica correção de bug. Não confundir.
- Não inclua nome do autor, arquivos alterados ou referências triviais que já aparecem no diff.

## Exemplos

Cada exemplo mostra as duas versões geradas para o mesmo conjunto de mudanças.

**Curta**

```
fix(config): corrigir leitura do diretório ~/.mj-cli quando HOME não está definido
```

**Longa**

```
fix(config): corrigir leitura do diretório ~/.mj-cli quando HOME não está definido

A leitura assumia que os.UserHomeDir() nunca falharia, o que quebrava
em ambientes containerizados sem HOME exportado.
```

---

**Curta**

```
refactor(cmd)!: trocar count de argumento posicional para flag --count
```

**Longa**

```
refactor(cmd)!: trocar count de argumento posicional para flag --count

Alinha a interface de `cube scramble` com as demais subcomandos, que
já usam flags nomeadas para parâmetros opcionais.

BREAKING CHANGE: o argumento posicional `count` em `cube scramble`
foi removido. Use `--count` (ou `-c`).
```

## Anti-padrões

- ❌ `update files` / `pequenos ajustes` / `wip`
- ❌ Mensagem em inglês quando o repositório usa pt-BR
- ❌ Descrição genérica que serviria para qualquer commit
- ❌ Commitar sem mostrar a mensagem para o usuário primeiro
