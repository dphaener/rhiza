# text-ref-check-verified

evidence · `01M6KM46W2JFPW4HGWY9T4EBW1` · actor [claude-code-e4cb86aa](../actor/claude-code-e4cb86aa.md)

## result

PASS

## statement

At 3baa77b7 (leaf commits 49a4311e + 3baa77b7), with RHIZA_ROOT/CHANGE/RUN unset: go test -count=1 -run 'TestTextRefs|TestRunOrder' -v ./internal/check passes (dangling @@{ID} gives one fail text-ref finding naming holder, statement and the ID with OK false; retired gives a warn with OK true; live, plain text and look-alikes '@{ID}' '@{BOGUS}' '@{<full ID>}' @@@{ID} give nothing). go test -count=1 ./cmd/rhiza -run 'TestScripts/(check-text-refs|describe-golden|text-refs-write)$' passes: check-text-refs forges a dangling @@{ID} into a log entry by hand and rhiza check exits non-zero with that fail text-ref finding (also --json rule text-ref, ok false). text-ref is in check.Rules() (between anchor and render, asserted by TestRunOrder) and appears twice in describe --json. make check: fmt-check, vet and every package's tests pass, and self-check of this repo prints ok with no text-ref finding (only two in-flight warns from other change sets); the attempt 1 render finding is fixed: rendered actor pages claude-code-39e4e367 and claude-code-83858632 both have tracked Record agent actor entries. make check exited 2 only on the timing budgets TestReadBudget, TestBudget (slice) and TestGeneratedModel (modelapp) at load average 11-14 from parallel suites; slice and modelapp do not import internal/check, and TestBudget fails the same way at the pre-leaf commit 6916c7c7 under the same load, so this is contention, not this leaf. Manual probe with a scratch build in a scratch repo: stored @@{<full ID>}, @@{ID}, @@{BOGUS} give no finding; a forged reference to the unassigned ID 01ZZZZZZZZZZZZZZZZZZZZZZZZ gives fail text-ref naming lit, statement and the ID. Noted, not this leaf's acceptance: a write of literal @@@{<a 26-character ID>} stores @@{<ID>}, which check then reads as a reference (the stored form keeps no escape), per text-ref-syntax's design.

## Links

### Out

- verifies [text-ref-check](../capability/rhiza/model/text-references/text-ref-check.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
