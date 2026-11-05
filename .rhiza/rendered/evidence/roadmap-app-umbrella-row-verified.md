# roadmap-app-umbrella-row-verified

evidence · `01M6A2Q7MZV52435XMHNASG9X4` · actor [claude-code-19607fcf](../actor/claude-code-19607fcf.md)

## result

PASS

## statement

At 90d02941 ran go test -count=1 -run TestUmbrella -v ./internal/tui/roadmapapp: TestUmbrellaGroupIsOneRowUntilEnter passes. It uses a fake runner returning an umbrella group and checks that the collapsed view shows one row 'rhiza / areas (2 change sets)' and no member rows, that enter on it expands to both change set rows (j then lands on the first one), that enter collapses it again, and that Result() still returns the runner's roadmap bytes unchanged. That last check is the --once path, which prints result() in cmd/rhiza/tui.go, and the commit leaves cmd/ untouched. Read the diff: umbrella rows are selectable, enter toggles before opening, and the cursor is kept across reloads. Also ran go test ./... (all ok), gofmt -l (clean) and go vet. make self-check passes with the session's RHIZA_* env unset. With the env set it flags only this verifier session's own actor render, which the commit does not cause.

## Links

### Out

- verifies roadmap-app-umbrella-row (retired)

## History

- `01M686RTJM0TBSN6C1FXES93WN` Umbrella change sets, which serve an intent marked umbrella, are grouped apart in the roadmap, are never ready, and collapse to one row in the roadmap app
