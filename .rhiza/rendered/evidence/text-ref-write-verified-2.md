# text-ref-write-verified-2

evidence · `01M6KR94CGPZKMCXT15H9PYP3Y` · actor [claude-code-477e958c](../actor/claude-code-477e958c.md)

## result

PASS

## statement

At e4375da1 (attempt 2, review cycle 2): go test -count=1 ./internal/api -run TestAddQuestionLenientRefs passes (a run session's question 'Use @@@store, not @@@ghost?' is written, stored as @@@{<full ID>} for store and literal @@@ghost, shown as 'Use @@@store, not @@@@@ghost?'). go test -count=1 ./cmd/rhiza -run TestScripts/(text-refs-write|text-refs-engine-writes|run-ship-change-rebase|describe-golden) passes: text-refs-write.txtar now writes a 20-char ID prefix and greps the full ID stored, and shows a note's @@@nowhere as @@@@@nowhere; text-refs-engine-writes.txtar shows question answer's Decision as 'Use @@@other, not @@@@@ghost' with the quote verbatim; run-ship-change-rebase.txtar conflicts on a file named @@@shipping, stores the question statement with @@@{<ID>}, and re-raising writes nothing (log bytes cmp equal). Mutation check in a scratch worktree: dropping the resolved-statement comparison in ship_rebase.go makes run-ship-change-rebase.txtar fail at line 203, so the dedupe test is sensitive. Scratch-build probe (clean env): rhiza add capability bad --field statement=@@nosuch exits 4 naming '@@nosuch in field statement' with next 'fix the handle, or write @@@@@ for a literal @@@', also with --dry-run, and writes nothing. guide text internal/guide/text/rhiza.md:82 tells agents to write @@@handle. make check: fmt-check, vet and test pass; self-check fails only on text-ref-check-verified (evidence for leaf text-ref-check, commit efb07bd4, not this leaf) holding a literal @@@{01ZZZZZZZZZZZZZZZZZZZZZZZZ}.

## Links

### Out

- verifies [text-ref-write](../capability/rhiza/model/text-references/text-ref-write.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
