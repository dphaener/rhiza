# root-continuations-first-verified

evidence · `01M6XXJA6JGZSXHR2B87BM1HDR` · actor [claude-code-bd671b3d](../actor/claude-code-bd671b3d.md)

## result

PASS

## statement

Verifier claude-code-bd671b3d at 2bde24d9: read the diff (internal/daemon/roundrobin.go fills slots in three tiers via ProjectPick.Continues/Queued; internal/api/daemon_projects.go daemonDemand reports the continuation and queued-start counts). Ran go test -trimpath -count=1 ./internal/daemon/...: ok, incl. TestRoundRobinContinuationsFirst (maxRuns 3, 2 running, A turn next with 2 ready, B 1 continuation -> [Bc]; 2 free slots, A queued start+ready, B continuation -> [Bc As]) and TestRoundRobinContinuationPerProjectMax, and the old cases TestRoundRobinAcrossProjects/PerProjectMax/RotatesAcrossPolls still pass unchanged. The repo test for 'else its start' only used a roadmap pick (Queued 0), so I checked it in a scratch copy of roundrobin.go: A at max with a continuation, B with a true queued start -> [B Bs]; with C also holding a continuation -> [C Cc]; the tier order holds whatever after is. go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-round-robin|daemon-continue-slot)$': both PASS (run with the run's RHIZA_* and GOTMPDIR unset, because GOTMPDIR under ~/code/rhiza/.git makes rhiza init in a script find the host repo; that is environmental, not caused by this change). go vet ./internal/api ./internal/daemon: clean.

## Links

### Out

- verifies [root-continuations-first](../capability/rhiza/orchestration/daemon/root-continuations-first.md)

## History

- `01M6XC4M3W7Y351WGXE8DTZACT` Fix A queued continuation can wait far longer than its slot should take: run 01M6WMX236BTT04TYQDACGKBZ4 (autonomy per workflow, PR #376 green) was queued with rhiza run resume at 2026-11-08T08:09Z and was still not started at 09:24Z; daemon status listed it as 'interrupted' among the live runs next to 3 running ones at rhiza.daemon.maxRuns 3, while two plain starts were 'waiting for a slot'. Check whether an interrupted run with a queued continuation counts against the cap itself (so it needs live <= cap-1 to start), and whether continuations take the next free slot ahead of new starts as #226/#243 intend. Unstuck by driving it with rhiza run resume --foreground.
