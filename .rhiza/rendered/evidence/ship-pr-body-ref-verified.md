# ship-pr-body-ref-verified

evidence · `01M6W07K71E88TEVY4BT8YWAFE` · actor [claude-code-753ca50a](../actor/claude-code-753ca50a.md)

## result

PASS

## statement

At 174b2970 ran go test -trimpath ./internal/api/... -run 'OpenPR' -count=1 -v: all five OpenPR tests pass, incl. TestShipOpenPRBodyRef. It opens a PR for a change set whose ref holds backticks, a newline and an @@mention and checks the body sent to the fake gh has exactly one line Ref: `https://x.test/1 \x60a\x60 @@mention\n# heading` (one inline code span, backticks as \x60, newline as \n, no line escapes the span), and that a change set without a ref gets no Ref: line. Read the diff: findShipTarget carries c.Open.Ref; prBodyRef returns empty for no ref and escapes backticks, CR/LF and other control runes before wrapping in one code span. gofmt and go vet clean.

## Links

### Out

- verifies [ship-pr-body-ref](../capability/ship-through-rhiza/ship-pr-body-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
