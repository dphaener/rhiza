# event-payload-ref-verified

evidence · `01M6VZZXBC1QWPFHZ448QE30M2` · actor [claude-code-52b3f164](../actor/claude-code-52b3f164.md)

## result

PASS

## statement

Verifier at a671d8aa: ran `go test -trimpath ./internal/api/... -run 'DaemonEvents' -count=1 -v`; TestDaemonEventsRef PASS. The test runs a recording event command (cat >> file) over moves of a change set opened with a ref containing quotes, spaces, & and non-ASCII (proposed, accepted) and one without (proposed); it decodes each payload and checks the first's ref equals the input verbatim and the second's has no ref key, with 2 and 1 payloads respectively. Diff read: lifecycleEvent gains `Ref string json:"ref,omitempty"`, set in daemonFold from the folded change set's c.Open.Ref for every lifecycle move in eventNames, so all moves (incl. rejected/withdrawn/superseded/landed/reopened) share the one path. The ref exists only on the open body, so it is the same on every move. gofmt and go vet are clean on internal/api.

## Links

### Out

- verifies [event-payload-ref](../capability/rhiza/hooks/model-event-hooks/event-payload-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
