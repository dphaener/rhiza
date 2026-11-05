# text-ref-syntax-verified

evidence · `01M6K6N4N191EM69Z3NH6PRNQH` · actor [claude-code-9aae00ea](../actor/claude-code-9aae00ea.md)

## result

PASS

## statement

Verified internal/textref at 6d9b05e3 (attempt 2). go test -count=1 ./internal/textref passes; FuzzRoundTrip and FuzzStoredRoundTrip each ran 30s clean. TestFind/TestResolve/TestFormat cover every acceptance point, and a scratch probe (removed) confirmed: references start at text start, after whitespace, ( or [; handles, IDs and 4+ char prefixes (@01M is none); darin@hey.com, checkout@v7, '@key' hold none; @@@ resolves to a literal @; Resolve rewrites @@alpha, @@01M6HP0C and @@{01M6HP0C} to @@{<full ID>} and returns @@{BOGUS} and @@01M6 unresolved with offsets; Format shows a live stored ref as @@alpha, a retired or unknown full ID as @@{<full ID>}, and stored @@{ID}, @@{BOGUS} and @@{<full ID>} as @@@{ID}, @@@{BOGUS}, @@@{<full ID>} (attempt 1's finding f1 is fixed and covered in TestFormat), and @@word as @@@word. fmt-check and vet pass. make check's other failures do not come from this leaf: a stale render of the builder's committed actor record, and the TestReadBudget/TestBudget timing budgets in internal/api and internal/slice at load average 8-11. Nothing imports textref, and the leaf only adds that package.

## Links

### Out

- verifies [text-ref-syntax](../capability/rhiza/model/text-references/text-ref-syntax.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
