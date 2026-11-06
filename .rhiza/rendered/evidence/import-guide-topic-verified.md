# import-guide-topic-verified

evidence · `01M6RMJ5H1KGHGKTNTMVAK4ZS4` · actor [claude-code-c01c241e](../actor/claude-code-c01c241e.md)

## result

PASS

## statement

At commit 45055620: a build of ./cmd/rhiza prints '- import  topic  502 tokens  Import: bringing an existing repo...' in `rhiza guide ls`, and `rhiza guide import` prints headings 1. Init, 2. Survey, 3. Record built capabilities, 4. Bring in load-bearing decisions, 5. Retire planning documents into Sources, 6. Write unbuilt work as draft change sets, 7. Verify, in that order. `go test -trimpath ./cmd/rhiza -run TestScripts/^guide$` passes; guide.txtar checks the import topic in guide ls and the seven steps in order with a (?s) regexp. `go test -trimpath ./internal/guide/...` passes, including TestCaps (size cap) and TestListOrder with import. Read the diff: commands and link types the text names (apply --record, decided-by, derived-from, change new, check, render, why) exist.

## Links

### Out

- verifies [import-guide-topic](../capability/rhiza/guide/import-guide-topic.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows
