# request-queue-verified

evidence · `01M6DA406V4DK910SMH4ET3KAB` · actor [claude-code-bc076fea](../actor/claude-code-bc076fea.md)

## result

PASS

## statement

At ebe38c10 (code landed in #191, the leaf commit records the builder's actor and render) I ran go test ./cmd/rhiza -run 'TestScripts/(plan-requests|describe-golden)$' -count=1: ok. I read plan-requests.txtar: with a fake rex it queues a feature, a bug and a milestone that request ls shows queued; paused, daemon --once starts nothing; unpaused it starts run start --workflow plan --new with the request's text, 'Fix: Login fails' for the bug and --milestone for the milestone, --autonomous only after rhiza.autonomy.plan auto; request ls then shows planning with the run, and once the run's journal names the planned change set request ls shows planned with it. go test ./internal/tui/roadmapapp -run NewWork -v passed: n and m call the request operation with --kind feature or milestone and open no rex tab, and cmd/rhiza/tui.go routes that runner through the same op registry as rhiza request. describe-golden.txtar holds request and request.ls and passes unchanged. make check exited 0.

## Links

### Out

- verifies [request-queue](../capability/rhiza/orchestration/daemon/plan-requests/request-queue.md)

## History

- `01M6D4SB0D8F9SPGK471DF9RV6` The request queue and the daemon turning requests into plan runs, as a leaf of plan requests
