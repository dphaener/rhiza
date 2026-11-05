# text-ref-check-verified-7

evidence · `01M6MSF1QK1J36XQCJY0B5PW4Y` · actor [claude-code-caa227c9](../actor/claude-code-caa227c9.md)

## result

PASS

## statement

Verifier, attempt 7, at 7b94d569: read internal/check/textrefs.go and its tests. go test ./internal/check -run 'TestTextRefs|TestRunOrder' passes: a dangling stored reference gives one fail naming holder, statement and the ID, with OK false; a retired target gives a warn with OK true; live, plain and look-alike text give nothing; a superseded value is not read; a retired holder is skipped. The check-text-refs, text-refs-* and describe-golden scripts pass: a hand-edited log entry with a dangling reference makes rhiza check exit non-zero with the fail text-ref finding, and text-ref is in check.Rules() and in describe. make check on a clean worktree at this commit exits 0, and rhiza check on this repo shows no text-ref findings.

## Links

### Out

- verifies [text-ref-check](../capability/rhiza/model/text-references/text-ref-check.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
