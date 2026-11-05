# ci-install-cancels-stale-verified

evidence · `01M6KYAW2VWT87C8EV9HZGX42V` · actor [claude-code-7d5f209c](../actor/claude-code-7d5f209c.md)

## result

PASS

## statement

At e3704f67 I read the diff (internal/api/ci.go template and cmd/rhiza/testdata/script/ci-install.txtar expected rhiza.yml, identical additions) and ran go test -trimpath ./cmd/rhiza -run TestScripts/ci-install -count=1 (ok) and go test -trimpath ./internal/api -count=1 (ok). The expected workflow, parsed as YAML, holds concurrency group 'rhiza-${{ github.event.pull_request.number || github.ref }}' (keyed on the PR, falling back to the ref) with cancel-in-progress '${{ github.event_name == 'pull_request' }}', so only pull_request events cancel and pushes to main never do; on still has push branches [main] and pull_request, and the base-branch, before-push and first-push check steps are unchanged. The testscript cmps the installed file byte-for-byte against this expected workflow. gofmt -l reports nothing.

## Links

### Out

- verifies [ci-install-cancels-stale](../capability/rhiza/integrity/ci-install-cancels-stale.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
