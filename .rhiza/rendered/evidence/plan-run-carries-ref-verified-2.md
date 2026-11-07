# plan-run-carries-ref-verified-2

evidence · `01M6W1TV6AV4YRRQAE1JHZQ2W4` · actor [claude-code-0bfd94a9](../actor/claude-code-0bfd94a9.md)

## result

PASS

## statement

Verifier at 44e2caf9 (attempt 2, fixing finding f1). `go test -trimpath ./internal/api/... -run Request` passes: TestRequestPlanRunArgsPassRef checks planRunArgs passes the ref as the single word --ref=<ref> for https://example.com/issues/7, --dry-run, --change=<id> and -h, and passes no --ref word without a ref. `go test -trimpath ./cmd/rhiza -run TestScripts/request-ref` passes, including the run start --workflow plan --new 'Planned' --ref step whose change show --json has "ref": "https://example.com/issues/9". Read cmd/rhiza/argv.go parseGlobals: --ref=... matches no global name and passes through whole. With a scratch build in a scratch repo, `change new A --ref=--dry-run --json` wrote a real change set whose change show --json has "ref": "--dry-run", so the ref is no longer read as a global flag.

## Links

### Out

- verifies [plan-run-carries-ref](../capability/rhiza/orchestration/daemon/plan-requests/plan-run-carries-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
