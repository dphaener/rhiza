# terminal-current-next-verified

evidence · `01M6CWXDQV9K0KWPX0Z0BWGW61` · actor [claude-code-0a2e2974](../actor/claude-code-0a2e2974.md)

## result

PASS

## statement

At 61d0963e: go test ./cmd/rhiza -run TestScripts/status-terminal-current passes; go test ./... all ok, gofmt clean, go vet clean (self-check failed only on stale rendered actor files from this run's sessions). By hand in a scratch git repo with a scratch build: created change set A, landed it with --change A (propose/accept/land), wrote A into .rhiza/local/change; rhiza status showed '(landed): Land me' and 'next: rhiza change use --none'; rhiza status --hook printed 'The current change set is landed; clear it with `rhiza change use --none`.'; rhiza add capability other exited 3 with 'next: rhiza change use --none', and --json gave code conflict with the same next; the file was not cleared implicitly. A withdrawn current change set got the same next and hook hint. After rhiza change use --none, status showed 'Change set: none current' with no hint and the hook printed 'No current change set.' Diff read: status adds the next only when the current file (not --change) names a terminal change set; writeTarget uses it only when the refused set is the current one.

## Links

### Out

- verifies [terminal-current-next](../capability/rhiza/change-sets/current-change/terminal-current-next.md)

## History

- `01M6CA9QMHN6N95AS3F935ZVN2` daemon status keeps the daemon's state through a re-exec, and a terminal current change set comes with the command that clears it
