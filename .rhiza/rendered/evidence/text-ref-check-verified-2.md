# text-ref-check-verified-2

evidence · `01M6KV9CHA7FXCDXX1A86VG0W4` · actor [claude-code-c9006c83](../actor/claude-code-c9006c83.md)

## result

PASS

## statement

Verifier, attempt 2 at 427029f8. Read the fix diff: textRefs now scans only the folded nodes' current fields in each change set's view, attributing each field to the change set that last wrote it (lastWriter over node history), so superseded values are never read; the old Evidence's statement is revised through a repeated op ID (a documented fold revision) and no longer stores the forged ID as a reference. With RHIZA_ROOT/CHANGE/RUN unset: go test -count=1 -run 'TestTextRefs|TestRunOrder' -v ./internal/check passes, covering dangling (one fail naming holder, statement and the ID, OK false), retired (warn, OK true), live, plain text and the quoted look-alikes (nothing), and the new TestTextRefsCurrentValue (a reference a later write dropped gives nothing; one a later write added fails). go test ./cmd/rhiza -run 'TestScripts/(check-text-refs|describe-golden)$' passes: check-text-refs forges a dangling reference into a log entry by hand and rhiza check exits non-zero with that fail text-ref finding (also --json rule text-ref, ok false), and describe lists anchor, text-ref, render; describe-golden holds text-ref twice. A scratch build's rhiza check on this repo prints ok with no text-ref finding (only two in-flight warns from other change sets), so the attempt 1 blocker is fixed. make check: fmt-check, vet and every package's tests pass; it exited non-zero only on timing budgets TestReadBudget (internal/api) and TestGeneratedModel (modelapp, which does not import internal/check) at load average 6-8, and both pass when rerun alone with -p 1; make fmt-check vet self-check then passes.

## Links

### Out

- verifies [text-ref-check](../capability/rhiza/model/text-references/text-ref-check.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
