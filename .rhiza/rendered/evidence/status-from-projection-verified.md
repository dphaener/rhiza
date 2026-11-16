# status-from-projection-verified

evidence · `01M7H5R7SSSXKEW3YXFAJFJAY7` · actor [claude-code-091a07ae](../actor/claude-code-091a07ae.md)

## result

PASS

## statement

Verifier re-ran at 05746877d, each in a capped systemd-run scope with RHIZA_/SHIP_/REX_ unset: go test -trimpath -count=1 -p 2 ./internal/api -run 'Status|TestProjection' ok; -run Status ok; -run TestPropertyProjectionParity ok; ./cmd/rhiza -run TestScripts/status ok (status, status-terminal-current and the other status-named scripts pass); go vet ./internal/api exit 0. Diff read: status uses the projection's new changeSetInfos (index change_sets rows, with pending change sets folded again and re-sorted) when IndexPath is set, and FoldLog otherwise; TestStatusIndexedReadsNoLog shows status and status --hook read the log 0 times over the index (1 over the log) with equal counts, change set and text; TestProjectionDryRunSeesPending shows a pending change set appears as a draft with its kind and intent, sorted, count equal to changeSetIDs; TestStatusHookBadIndexSilent shows an unopenable index errors for status and the hook is silent with exit 0. Only internal/api status/projection/write_indexed and their tests changed, no TODOs. Measured: scratch build in /home/rhiza/code/rhiza, one warming run, then /usr/bin/time status 0.06 s (hook 0.10 s; parent commit's build 1.40 s), target 200 ms met. status, status --json and status --hook output byte-identical (cmp) to the parent commit's build; status --hook outside a repo silent, exit 0.

## Links

### Out

- verifies [status-from-projection](../capability/rhiza/views/fast-reads/status-from-projection.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
