# change-set-ref-qa-1

evidence · `01M6W122X39SKC29Z826A3CCQP` · actor [claude-code-2590959f](../actor/claude-code-2590959f.md)

## result

PASS

## statement

QA at 66e50695 with a built binary in a scratch repo: change new 'Link it' --ref with ESC, newline and backticks; --json and change show --json returned it verbatim (\u001b, \n); human change show printed '- Ref: "https://ex.com/i/1\x1b[31m\n`x`"' with no raw ESC (cat -v); the open entry is stored at v 4 with a ref key, and a change set without --ref has none. A binary built from the merge-base refused the entry: 'newer format: field "v": 4, this Rhiza reads 3'. go test ./internal/changeset/... and TestScripts/change-set-ref pass.

## Links

### Out

- verifies [change-set-ref](../capability/rhiza/change-sets/change-set-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
