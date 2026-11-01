# feed-hides-actor-records-verified

evidence · `01M6C0DGQ6ZGR2XBA0J1E1T945` · actor [claude-code-86d97c77](../actor/claude-code-86d97c77.md)

## result

PASS

## statement

Verified at 3ee71a7e. I read the diff from cdfd5af5. internal/api/feed.go adds an `all` flag and feedActorOnly, which drops a landed change set when every op creates an actor or updates an actor node. The drop happens inside the landed loop, before paging, and the next-page command carries --all. cmd/rhiza/tui.go runs the same feed op. feed.txtar now lands alpha and gamma capability records and an intent record (beta). It checks that by default 3 are listed and 'Record actor' is not, and that with --all 4 are listed in order. With --limit 2 the first page holds 2 items and the next page holds the rest. With --all, the next: line contains --all and both pages are full. An intention change set is shown by default and with --all. tui-changes-app.txtar asserts that 'Record actor' is not shown. Re-running describe-golden with -update left no diff, and the golden file has the `all` field. `go test ./cmd/rhiza -run 'TestScripts/(feed|tui-changes-app|describe-golden)$' -count=1` passed. `go test ./...` passed in every package. `make self-check` passed with RHIZA_* unset. The attempt-1 finding (no intent record in the test) is fixed.

## Links

### Out

- verifies [feed-hides-actor-records](../capability/rhiza/views/change-feed/feed-command/feed-hides-actor-records.md)

## History

- `01M6BH7XPH07NSDTBKE2RGX949` rhiza feed and the changes app leave out actor-only change sets unless --all is given
