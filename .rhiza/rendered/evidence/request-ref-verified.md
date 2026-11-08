# request-ref-verified

evidence · `01M6W0SFZGT2FHCVNN2E790VEQ` · actor [claude-code-6ff7dea4](../actor/claude-code-6ff7dea4.md)

## result

PASS

## statement

Verifier claude-code-6ff7dea4 at 8aa0e068: ran go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/(request-ref|describe-golden) and every *request* script (request-ref, describe-golden, plan-requests, daemon-request-slot, daemon-milestone-request, tui-request-app): all PASS. request-ref.txtar queues a URL ref (JSON plain string, human [ref "..."]), an ESC+newline ref (JSON a\u001b[31mb\nc verbatim, human Go-quoted with no raw ESC byte), a space-padded ref (stored trimmed), and a 2049-byte ref (usage error naming 2048, nothing queued). The diff adds Ref to RequestRequest, RequestInfo and daemon.Request; it trims, checks len against changeset.MaxRef=2048 and UTF-8 validity, uses omitempty so an empty ref means none, and %q in human output. describe-golden is regenerated with the ref fields and passes. internal/api Request tests pass; go vet and gofmt are clean.

## Links

### Out

- verifies [request-ref](../capability/rhiza/orchestration/daemon/plan-requests/request-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
