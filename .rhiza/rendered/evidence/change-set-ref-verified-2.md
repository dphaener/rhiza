# change-set-ref-verified-2

evidence · `01M6W1QJAJG0G134FYJXZYSDNB` · actor [claude-code-29343599](../actor/claude-code-29343599.md)

## result

PASS

## statement

At ab36a226 (leaf commits 2438d674 and fix ab36a226), go test -trimpath -count=1 ./internal/changeset/... ./internal/entry/... ./internal/render/... passes. TestOpenRef round-trips an open with a ref (v4, decoded and canonical), checks that no ref leaves the encoding without a ref key at EntryVersion 1, rejects a 2049-byte ref while accepting 2048, and refuses the ref field at v1-v3. TestCorpusFolds, TestPropertyPermutation and render TestGolden pass on the new format-4/ref corpus case. go test -trimpath ./cmd/rhiza -run 'TestScripts/(change-set-ref|describe-golden)$' passes. The script checks that change new --ref, then change show --json and change show, print the ref. The human line is Go-quoted (\x1b and \n escaped, no raw ESC byte). The ref is trimmed and a blank one means none. An over-limit ref is a usage error that writes nothing, and fold-ok holds. Diff read: Ref is encoded only when set, decodeBody accepts it only at v>=4, and FormatVersion is now 4.

## Links

### Out

- verifies [change-set-ref](../capability/rhiza/change-sets/change-set-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
