# plan-run-carries-ref-verified

evidence · `01M6W0BTKAYFYJSAEFBF23018N` · actor [claude-code-320da176](../actor/claude-code-320da176.md)

## result

PASS

## statement

Verifier claude-code-320da176 at 73898e53: read diffs of ba19dd35 and 73898e53 (RunStartRequest.Ref, --ref requires --new, newPlanChange opens changeset.Open with Ref; daemonStartPlan detaches planRunArgs, which appends --ref <ref> only when the request has one). go test -trimpath ./internal/api/... -run 'Request' -count=1 passed, incl. TestRequestPlanRunArgsPassRef asserting planRunArgs carries --ref with the request's ref and no --ref without one. go test -trimpath ./cmd/rhiza -run TestScripts/request-ref passed, its step running run start --workflow plan --new 'Planned' --ref https://example.com/issues/9 then change show --json showing "ref": "https://example.com/issues/9". TestScripts/describe-golden passed; gofmt and go vet on internal/api clean.

## Links

### Out

- verifies [plan-run-carries-ref](../capability/rhiza/orchestration/daemon/plan-requests/plan-run-carries-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
