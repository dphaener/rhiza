# daemon-fold-from-projection-qa-1

evidence · `01M7H9P8GEWEV8ZBRAQE7WTN4Q` · actor [claude-code-2c502734](../actor/claude-code-2c502734.md)

## result

PASS

## statement

QA: scratch builds of HEAD 9f005fded and of base 03b972942 (go build -trimpath, capped scope). rhiza daemon --once in two scratch roots (base and HEAD) over a log with a proposed change set with a leaf, one without and two overlapping drafts: same output, both picked 'Proposed with leaf' for plan-review and skipped the leafless one; daemon status --json identical.

## Links

### Out

- verifies [daemon-fold-from-projection](../capability/rhiza/views/fast-reads/daemon-fold-from-projection.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
