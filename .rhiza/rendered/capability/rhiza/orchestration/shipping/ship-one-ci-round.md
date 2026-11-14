# ship-one-ci-round

capability · `01M79KPE2V3AZWKX2HR1X1PZ1M` · parent [shipping](../shipping.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScripts/ship-change-one-ci shows, with a fake gh and a fake conductor that lands and renders on the branch, exactly one CI wait and one push on the happy path, two CI waits when main moves before the merge, and the change set landed on main after the merge with no ship land step. Touches cmd/rhiza/testdata/script/ship-change-one-ci.txtar, and internal/api/ship*.go only if a ship command needs the change set unlanded.

## statement

On the conductor default the change set lands on the branch before the first push, so a run whose PR is green and whose main has not moved waits on CI once; a second CI round happens only when the merge finds main has moved.

## Links

### Out

- depends-on [builtin-conductor-default](../../workflows/built-in-workflows/builtin-conductor-default.md)
- decided-by [conductor-milestone-direction](../../../../decision/conductor-milestone-direction.md)

### In

- verifies [ship-one-ci-round-verified](../../../../evidence/ship-one-ci-round-verified.md)
- verifies [ship-one-ci-round-qa-2](../../../../evidence/ship-one-ci-round-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
