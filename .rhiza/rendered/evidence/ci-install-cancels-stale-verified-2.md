# ci-install-cancels-stale-verified-2

evidence · `01M6M1YN4V8TQVRP0C6PZ3QKAN` · actor [claude-code-a50d0083](../actor/claude-code-a50d0083.md)

## result

PASS

## statement

At cd64d73f ran go test ./cmd/rhiza -run 'TestScripts/ci-install$' -count=1: ok. The script cmps .github/workflows/rhiza.yml against its expected workflow, which holds concurrency group rhiza-${{ github.event.pull_request.number || github.sha }} with cancel-in-progress ${{ github.event_name == 'pull_request' }} (true only for pull_request, never for pushes to main), and keeps the push: branches [ main ] trigger and all three check steps (base branch, github.event.before, first push). internal/api/ci.go's template matches byte for byte. Read the diff: the fallback moved from github.ref to github.sha, departing from the acceptance's 'falling back to the ref' wording, to fix confirmed review finding f1 (GitHub cancels a group's pending run when a newer one queues, so one shared push group could drop a push and skip its changes in the incremental check). The sha fallback gives each push its own group, which is what the leaf statement requires: no push to main cancels an earlier push's check.

## Links

### Out

- verifies [ci-install-cancels-stale](../capability/rhiza/integrity/ci-install-cancels-stale.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
