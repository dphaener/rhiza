# event-payload-ref

capability · `01M6VYMKEVBTRJ5WA7618HFQF8` · parent [model-event-hooks](../model-event-hooks.md)

## acceptance

go test -trimpath ./internal/api/... -run 'DaemonEvents' passes with a test that runs a recording event command over moves of a change set with a ref and one without: each payload of the first has the ref verbatim, the second's has no ref key

## statement

The rhiza.eventCommand JSON payload of every lifecycle move of a change set with a ref (proposed, accepted, rejected, withdrawn, superseded, landed, reopened) carries the ref, read from the folded change set, and omits it when the change set has none

## Links

### Out

- decided-by [request-ref-shape](../../../../decision/request-ref-shape.md)

### In

- verifies [event-payload-ref-verified](../../../../evidence/event-payload-ref-verified.md)
- verifies [event-payload-ref-qa-1](../../../../evidence/event-payload-ref-qa-1.md)
- verifies [pr-357-requests-carry-an-opaque](../../../../evidence/pr-357-requests-carry-an-opaque.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
