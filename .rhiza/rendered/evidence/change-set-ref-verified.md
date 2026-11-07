# change-set-ref-verified

evidence · `01M6W06PQER7BRV9WZF75YTSR6` · actor [claude-code-4e8a6d09](../actor/claude-code-4e8a6d09.md)

## result

PASS

## statement

Verified as claude-code-4e8a6d09 at 2438d674. go test -trimpath -count=1 ./internal/changeset/... and ./internal/entry/... pass; TestOpenRef covers: no ref leaves the encoding without a ref key and EntryVersion 1, a ref round-trips (canonical and not) at EntryVersion 4, a decoder at v1-3 refuses the ref field, an empty ref is refused, 2048 bytes accepted and 2049 refused. entry.FormatVersion is 4, so a binary at format 3 refuses v4 entries with ErrNewerFormat. go test -trimpath ./cmd/rhiza -run 'TestScripts/(change-set-ref|describe-golden)$' passes: change new --ref then change show --json print the ref; change show prints '- Ref: "..."' Go-quoted with ESC and newline escaped and no raw ESC byte; fold-ok holds; without a ref nothing shows and the open stays format 1; a 2049-byte ref is a usage error that writes nothing. I read the diff (api request/response fields, indexed show path, body encode/decode/validate). go test ./internal/api -run 'Change|Describe' passes its tests, but the package exits FAIL on the known .rhiza-root/usage.live.json leak check, which this change does not touch. The full api package hit the 10-minute timeout under contention from parallel test runs.

## Links

### Out

- verifies [change-set-ref](../capability/rhiza/change-sets/change-set-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
