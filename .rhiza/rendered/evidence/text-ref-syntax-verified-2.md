# text-ref-syntax-verified-2

evidence · `01M6KS126YR7T580JWVR71R5EZ` · actor [claude-code-6c11e4d5](../actor/claude-code-6c11e4d5.md)

## result

PASS

## statement

Verifier, attempt 2 at 22bd296a. Read internal/textref and the fix diff. go test -count=1 -v ./internal/textref passes. TestFind covers: a reference starts at text start or after whitespace, ( or [; handles, IDs and 4+ char prefixes are found; darin@hey.com, checkout@v7, '@key' and @@@alpha hold none. TestResolve: references, @@{<ID>} included, become @@{<full ID>}; @@{BOGUS} and @@nope come back unresolved with their offsets. New TestResolveLiteralFullID: an @@@{<full ID>} escape and an unmapped @@{<full ID>} are returned (Escaped set) and stored without braces, so they never read back as references. TestFormat: a live node shows @@handle; a retired or missing node shows @@{<full ID>}; stored @@{ID}, @@{BOGUS}, @@{<full ID>} and @@{...} show as @@@{...}; stored @@word shows @@@word. Fuzz runs: FuzzRoundTrip for 45s, about 2.27M execs, which now also checks that every stored reference came from lookup; FuzzStoredRoundTrip for 45s, about 1.92M execs. Neither failed. go test ./cmd/rhiza -run TestScripts/text-refs passes, including the new text-refs-write.txtar cases: a strict @@@{<full ID>} is refused with a usage error, and a lenient note stores the brace-less form with fold-ok. make fmt-check vet test passes. make self-check fails only on stale .rhiza/rendered (missing critic actor pages, index.md). Those come from the earlier critic chore commits; this commit touches no .rhiza file.

## Links

### Out

- verifies [text-ref-syntax](../capability/rhiza/model/text-references/text-ref-syntax.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
