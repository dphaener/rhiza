# change-set-ref

capability · `01M6VYMKCZ09KBA0BHCVXF1Y7G` · parent [change-sets](../change-sets.md)

## acceptance

go test -trimpath ./internal/changeset/... passes with tests that round-trip an open body with and without a ref (no ref leaves the encoding and format unchanged), reject a ref over 2048 bytes, and refuse the ref's format in a reader at the previous version; go test -trimpath ./cmd/rhiza -run TestScripts/change-set-ref passes: change new --ref then change show --json and change show print the ref, quoted with no raw control bytes in human output, and fold-ok holds

## statement

A change set carries an optional opaque ref on its open entry, set by rhiza change new --ref under the same 2048-byte limit, encoded only when set and stamped with a new entry format version so an older binary refuses it; change show shows it as a plain string in JSON and Go-quoted in human output

## Links

### Out

- decided-by [request-ref-shape](../../../decision/request-ref-shape.md)

### In

- verifies [change-set-ref-verified](../../../evidence/change-set-ref-verified.md)
- verifies [change-set-ref-qa-1](../../../evidence/change-set-ref-qa-1.md)
- verifies [change-set-ref-verified-2](../../../evidence/change-set-ref-verified-2.md)
- verifies [pr-357-requests-carry-an-opaque](../../../evidence/pr-357-requests-carry-an-opaque.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
