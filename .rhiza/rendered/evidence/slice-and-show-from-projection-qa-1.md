# slice-and-show-from-projection-qa-1

evidence · `01M7H9P8K9H0N3WWBD6D14MHHP` · actor [claude-code-2c502734](../actor/claude-code-2c502734.md)

## result

PASS

## statement

QA: scratch builds of HEAD 9f005fded and of base 03b972942 (go build -trimpath, capped scope). slice, slice --json, change show and --json byte-identical to base in /home/rhiza/code/rhiza (780 / 1266 lines for the current change set) and in a scratch repo (what-if with an in-flight overlap, after a vanished folder: 'no change set' as base). slice outside a repo: not_found, next rhiza init; change show with none current: the same usage error, exit 2.

## Links

### Out

- verifies [slice-and-show-from-projection](../capability/rhiza/views/fast-reads/slice-and-show-from-projection.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
