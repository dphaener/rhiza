# text-ref-syntax-qa-1

evidence · `01M6KNZR9PJA7T1KJYWYVD36DM` · actor [claude-code-5c694237](../actor/claude-code-5c694237.md)

## result

PASS

## statement

QA cycle 1 at 0e9f6a6e, scratch build in a scratch repo (RHIZA_* unset): 'Builds on @@other, me@x.org, checkout@v7, '@key', @@@lit, (@@local) [@<24-char ID prefix>]' stored @@{<other ID>}, @@{<local ID>} and @@{<other ID>} for the three references, kept the email, checkout@v7 and '@key' as written and stored @@@lit as @@lit; show printed '@renamed ... @@@lit ... (@@local) [@@renamed]' after a rename; an ambiguous 4-char prefix @@01M6 and @@{BOGUS} were unresolved. go test ./internal/textref passes. Not passing: a write of @@@{<26-char ID>} stores @@{<ID>}, which Format and check then treat as a reference (see the qa finding).

## Links

### Out

- verifies [text-ref-syntax](../capability/rhiza/model/text-references/text-ref-syntax.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
