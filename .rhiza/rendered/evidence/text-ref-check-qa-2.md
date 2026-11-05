# text-ref-check-qa-2

evidence · `01M6KVTZTYWPGMXBHZTC8ZYP4A` · actor [claude-code-e64add76](../actor/claude-code-e64add76.md)

## result

PASS

## statement

Scratch build: a hand-edited log entry holding a braced reference to a forged 26-character ID made rhiza check exit 1 with 'fail text-ref ... victim field statement holds @@{01ZZ...}, which names no node'; a kept reference to a retired node gave a warn with check still ok; hand-written '@{ID}' gave nothing. On this repo at the branch head with a fresh cache, rhiza check reports no text-ref finding (cycle-1 f1 fixed); its only fails are stale rendered actor/index files from the verifiers' evidence commits, which the run's gate render regenerates

## Links

### Out

- verifies [text-ref-check](../capability/rhiza/model/text-references/text-ref-check.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
