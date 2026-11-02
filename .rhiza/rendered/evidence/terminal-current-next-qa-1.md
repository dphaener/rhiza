# terminal-current-next-qa-1

evidence · `01M6CXVSG9FW5PHHESG7WQ3NYF` · actor [claude-code-4d1d9a27](../actor/claude-code-4d1d9a27.md)

## result

PASS

## statement

At 423ba9c5, by hand with a scratch build in a scratch git repo, RHIZA_* unset: created change set A, propose/accept/land, wrote A into .rhiza/local/change. rhiza status showed '(landed): Daemon test cs' and 'next: rhiza change use --none'; rhiza status --hook printed 'The current change set is landed; clear it with `rhiza change use --none`.'; rhiza add capability other exited 3 with 'next: rhiza change use --none', and --json gave code conflict with the same next; the current file still named A afterwards. After rhiza change use --none, status showed 'Change set: none current' with no hint and the hook 'No current change set.' go test ./cmd/rhiza -run TestScripts/status-terminal-current passes.

## Links

### Out

- verifies [terminal-current-next](../capability/rhiza/change-sets/current-change/terminal-current-next.md)

## History

- `01M6CA9QMHN6N95AS3F935ZVN2` daemon status keeps the daemon's state through a re-exec, and a terminal current change set comes with the command that clears it
