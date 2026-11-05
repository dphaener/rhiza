# text-ref-write-verified

evidence · `01M6K80PZDBN9YKGVKMQJG9AB7` · actor [claude-code-09d03f6f](../actor/claude-code-09d03f6f.md)

## result

PASS

## statement

Verifier claude-code-09d03f6f at f65c3f54: go test ./cmd/rhiza -run 'TestScripts/(text-refs|describe-golden)' passes (text-refs-write.txtar, text-refs-engine-writes.txtar, describe-golden unchanged since lenient is unexported). make check passed fmt-check, vet and test, then self-check failed only on stale .rhiza/rendered files (no render before landing, as expected). Ran a scratch repo by hand with a fresh build under env -i. @<16-char ID prefix> stores @@{full ID}. An add with @@nosuch in statement and @@nope2 in acceptance exits 4 naming both refs with their fields and next 'fix the handle, or write @@@ for a literal @', and nothing is written. --json gives code not_found. A source's summary/title keep @@nosuch/@other verbatim. After gone is retired: set holder statement=@{gone ID} where statement already held it succeeds and stores it; the same ref in holder's acceptance field and in a new node exits 4. Read the diff: resolution is in buildAdd/buildSet with strict/lenient modes. Engine writes (notes, question answer, run/check/leaf/review/interview/rebase questions, ship evidence) are lenient per text-ref-engine-writes. The rebase dedupe compares against the leniently resolved statement. The guide tells agents to write @@handle.

## Links

### Out

- verifies [text-ref-write](../capability/rhiza/model/text-references/text-ref-write.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
