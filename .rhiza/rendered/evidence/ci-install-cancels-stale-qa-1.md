# ci-install-cancels-stale-qa-1

evidence · `01M6M04BERDRJGC1Y1TRBYERZS` · actor [claude-code-ca676fb4](../actor/claude-code-ca676fb4.md)

## result

PASS

## statement

At 83445164 I built rhiza (go build ./cmd/rhiza) and ran rhiza init then rhiza ci install in a fresh scratch git repo: it wrote .github/workflows/rhiza.yml, which PyYAML parses to on {push: {branches: [main]}, pull_request}, concurrency group 'rhiza-${{ github.event.pull_request.number || github.ref }}' with cancel-in-progress "${{ github.event_name == 'pull_request' }}", and keeps the steps Check against the base branch, Check against the commit before the push and Check a first push. go test -trimpath ./cmd/rhiza -run TestScripts/ci-install -count=1 passed.

## Links

### Out

- verifies [ci-install-cancels-stale](../capability/rhiza/integrity/ci-install-cancels-stale.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
