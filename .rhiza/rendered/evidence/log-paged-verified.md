# log-paged-verified

evidence · `01M6HV04Z3XE6NXZF5KSVFGVWZ` · actor [claude-code-0c0f2e6c](../actor/claude-code-0c0f2e6c.md)

## result

PASS

## statement

At a34e86f3 I read the diff (internal/api/lineage.go, log-paged.txtar, describe-golden.txtar) and ran go test ./cmd/rhiza -run 'TestScripts/(log-paged|describe-golden|log)': ok. log-paged.txtar logs a node touched by three change sets: --limit 2 gives two, a next_cursor and a 'next: rhiza log store --limit 2 --cursor <c>' line; that cursor's page holds the third with no next_cursor; --all gives three; --limit 0 gives the default page; --limit -1, --limit 1001 and a bad cursor exit 2; log --json and api call log cmp equal. Regenerating describe-golden with -update left git status clean. By hand with a scratch build on this repo: log rhiza --limit 1 printed the next: line and next_cursor; --limit 1001 and --cursor zz exited 2 with 'next: rhiza describe log'. fmt-check, vet and self-check pass. Under load average 14-22 on 4 cores, go test ./... failed only run-command-step-output (passes 3/3 alone; the leaf touches no run code) and the modelapp tree/roadmap read-timing budget (log is not read there): runner contention, not this leaf.

## Links

### Out

- verifies [log-paged](../capability/rhiza/views/paged-lists/log-paged.md)

## History

- `01M6H3A65RDP761HKE82STZSQZ` log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan)
