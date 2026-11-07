# request-ref

capability · `01M6VYMKBYZQBDVR648YSHRVEY` · parent [plan-requests](../plan-requests.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScripts/request-ref passes: a script queues a request with a ref holding a URL, one with an escape sequence and newline (request ls --json returns it verbatim, human request ls prints it quoted with no raw control bytes), one with surrounding spaces (trimmed), one over 2048 bytes (usage error naming 2048, nothing queued); describe-golden is regenerated

## statement

rhiza request '<text>' --ref <string> stores an opaque external reference on the queued request, trimmed and at most 2048 bytes (a longer one is a usage error naming the limit, an empty one means none), and request ls shows it: as a plain string in JSON, Go-quoted in human output

## Links

### Out

- decided-by [request-ref-shape](../../../../../decision/request-ref-shape.md)

### In

- verifies [request-ref-verified](../../../../../evidence/request-ref-verified.md)
- verifies [request-ref-qa-1](../../../../../evidence/request-ref-qa-1.md)
- verifies [pr-357-requests-carry-an-opaque](../../../../../evidence/pr-357-requests-carry-an-opaque.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
