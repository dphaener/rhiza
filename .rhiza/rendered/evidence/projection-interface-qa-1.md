# projection-interface-qa-1

evidence · `01M7H9P7TZ486T22AZ72Y0JXJH` · actor [claude-code-2c502734](../actor/claude-code-2c502734.md)

## result

PASS

## statement

QA: scratch builds of HEAD 9f005fded and of base 03b972942 (go build -trimpath, capped scope). In a scratch repo, apply --dry-run of four chained writes (add sso, add question, link to sso, set sso) gave the same output as base (IDs aside) and wrote nothing; the real apply then showed sso in slice and show. A garbage index file was rebuilt and status, change ls and render answered as before; an unopenable cache dir gives 'open index: ... permission denied' instead of a silent fold, as the leaf requires.

## Links

### Out

- verifies [projection-interface](../capability/rhiza/views/fast-reads/projection-interface.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
