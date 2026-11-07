# ship-pr-body-ref

capability · `01M6VYMKG6TVD7M4Y8QM04GP38` · parent [ship-through-rhiza](../ship-through-rhiza.md)

## acceptance

go test -trimpath ./internal/api/... -run 'OpenPR' passes with tests that the body written for a change set with a ref holding backticks, a newline and an @@mention shows it in one escaped inline code span, and the body for a change set without a ref has no ref line

## statement

ship open-pr lists the change set's ref in the PR body, inside an inline code span with backticks and newlines escaped so the untrusted text cannot inject markdown or mentions, and lists nothing when the change set has no ref

## Links

### Out

- decided-by [request-ref-shape](../../decision/request-ref-shape.md)

### In

- verifies [ship-pr-body-ref-verified](../../evidence/ship-pr-body-ref-verified.md)
- verifies [ship-pr-body-ref-qa-1](../../evidence/ship-pr-body-ref-qa-1.md)
- verifies [pr-357-requests-carry-an-opaque](../../evidence/pr-357-requests-carry-an-opaque.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
