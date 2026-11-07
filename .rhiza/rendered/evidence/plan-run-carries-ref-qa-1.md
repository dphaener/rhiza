# plan-run-carries-ref-qa-1

evidence · `01M6W12327KVMDWV6F5PN8662B` · actor [claude-code-2590959f](../actor/claude-code-2590959f.md)

## result

PASS

## statement

QA at 66e50695 with a built binary: run start --workflow plan --new 'Planned' --ref 'https://ex.com/9' --permission-mode manual queued the plan run and opened a change set whose change show --json has "ref":"https://ex.com/9"; run start --workflow plan --ref x without --new is a usage error '--ref goes with --new'. The daemon's detached start (needs Rex and an agent) was not run by hand; TestRequestPlanRunArgsPassRef covers that the plan run arguments carry --ref for a request with one and none without.

## Links

### Out

- verifies [plan-run-carries-ref](../capability/rhiza/orchestration/daemon/plan-requests/plan-run-carries-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
