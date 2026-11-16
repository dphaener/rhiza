# inbox-from-projection-verified

evidence · `01M7H6FDB1MZV007KKZ8QNBJ9V` · actor [claude-code-2025b3d3](../actor/claude-code-2025b3d3.md)

## result

PASS

## statement

Verified fea2a6b29 at HEAD under the capped scope: go test -trimpath -count=1 -p 2 ./internal/api -run Inbox ok (includes TestInboxIndexedReadsNoLog: zero log reads with IndexPath set, overlaps/titles/gate questions equal to the log fold, a gate run whose worktree is gone keeps its item with the question handle, no statement, no error); go test ./internal/api -run 'TestPropertyProjectionParity|Proposal|Uncovered|Overlap|Daemon' ok; go test ./cmd/rhiza -run TestScripts/inbox ok over all 20 inbox*.txtar (RHIZA_/SHIP_/REX_ unset); go vet ./internal/api exit 0; grep -n 'index.Open\|openWriteIndex' internal/api/inbox.go prints nothing (exit 1). Diff read: inboxFold uses ws.useProjection/changeSetInfos/entries/model, overlaps name targets via overlapItems over the projection's as-is model (same handle/statement as before); inboxUncoveredProposals (callers inbox.go:841 and :1236 only) uses only len(leaves), so the nil rank changes no item or --json; inboxGateRead already uses openProjection and is unchanged; internal/run untouched; no TODOs. Discovered-from draft 01M7H10BV65XX3R09JMQDM49C2 'Inbox within 200 ms' exists per Decision inbox-rescoped-after-parity.

## Links

### Out

- verifies [inbox-from-projection](../capability/rhiza/views/fast-reads/inbox-from-projection.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
