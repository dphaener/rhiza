# plan-run-carries-ref

capability · `01M6VYMKDW2EP33TZA6ABET0V8` · parent [plan-requests](../plan-requests.md)

## acceptance

go test -trimpath ./internal/api/... -run 'Request' passes with a test that the daemon's run start arguments for a request with a ref pass it on, and one without a ref pass none; go test -trimpath ./cmd/rhiza -run TestScripts/request-ref passes its step showing run start --new with the request's ref opens a change set whose change show --json has the ref

## statement

When the daemon starts the plan run for a request with a ref, the run opens its change set with that ref, so the planned change set's open entry carries the request's ref

## Links

### Out

- decided-by [request-ref-shape](../../../../../decision/request-ref-shape.md)

### In

- verifies [plan-run-carries-ref-verified](../../../../../evidence/plan-run-carries-ref-verified.md)
- verifies [plan-run-carries-ref-qa-1](../../../../../evidence/plan-run-carries-ref-qa-1.md)
- verifies [plan-run-carries-ref-verified-2](../../../../../evidence/plan-run-carries-ref-verified-2.md)
- verifies [pr-357-requests-carry-an-opaque](../../../../../evidence/pr-357-requests-carry-an-opaque.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
