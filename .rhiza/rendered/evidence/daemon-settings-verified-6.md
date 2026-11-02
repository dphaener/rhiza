# daemon-settings-verified-6

evidence · `01M6CZX4KQEEEWTADYXGV7K8YF` · actor [claude-code-f3436fe2](../actor/claude-code-f3436fe2.md)

## result

PASS

## statement

At d8ca4bc4 (an empty commit on fix 15d09912), go test ./... passes with gofmt and go vet clean, including daemon-settings.txtar, stage-autonomy.txtar, daemon-pause and daemon-lock. Regenerating describe-golden.txtar with -update leaves no diff. daemon-settings.txtar checks: every key listed with its default and source default; daemon set rhiza.daemon.maxRuns 1 then settings shows 1 (git config); a bad key or value is refused and writes nothing; with maxRuns 1, a --once poll starts exactly one of two ready change sets; --max-runs 5 shows source flag in daemon status; the next poll without the flag reads config again. By hand with this commit's build in a scratch repo, daemon set rhiza.daemon.nope 1, maxRuns 0, autonomy.accept maybe and daemon.paused maybe each exited 2 and left no rhiza.daemon or rhiza.autonomy keys. In the diff, autonomyStage now takes each stage's default from settingDefs. Its three callers (run gate accept, daemon ship and signoff) pass no default of their own, and accept's default is auto, so daemon settings shows accept = auto (default), which is what the gate does. This fixes the confirmed review f1, and stage-autonomy.txtar checks it with the key unset. make self-check fails only because .rhiza/rendered actor files and index.md are stale from this run's own model writes. No code is involved.

## Links

### Out

- verifies [daemon-settings](../capability/rhiza/orchestration/daemon/daemon-settings.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
