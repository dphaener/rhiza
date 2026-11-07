# event-payload-ref-qa-1

evidence · `01M6W122Z27DS8K7BV9W0QJYWR` · actor [claude-code-2590959f](../actor/claude-code-2590959f.md)

## result

PASS

## statement

QA at 66e50695 with a built binary: a scratch root plus project, daemon --once --max-runs 0 --event-command appending stdin to a file; after the recording first poll, proposed and accepted moves of a change set with --ref https://github.com/o/r/issues/42 and a proposed move of one without. The second poll wrote three payloads: both moves of the first carry "ref":"https://github.com/o/r/issues/42" verbatim, the plain one has no ref key. TestDaemonEventsRef passes.

## Links

### Out

- verifies [event-payload-ref](../capability/rhiza/hooks/model-event-hooks/event-payload-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
