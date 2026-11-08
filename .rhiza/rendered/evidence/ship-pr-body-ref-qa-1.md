# ship-pr-body-ref-qa-1

evidence · `01M6W12389SC0QE6KMY04T3SPP` · actor [claude-code-2590959f](../actor/claude-code-2590959f.md)

## result

PASS

## statement

QA at 66e50695: read prBodyRef and its use in ship open-pr (the ref from the change set's open, appended to the summary as one inline code span with backticks, CR/LF and control characters escaped). Opening a real PR needs gh and was not run by hand; go test ./internal/api -run OpenPR passes, including TestShipOpenPRBodyRef (backticks, newline and @@mention ref in one escaped span; no ref line without a ref).

## Links

### Out

- verifies [ship-pr-body-ref](../capability/ship-through-rhiza/ship-pr-body-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
