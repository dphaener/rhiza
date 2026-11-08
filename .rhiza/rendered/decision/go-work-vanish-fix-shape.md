# go-work-vanish-fix-shape

decision · `01M6WGVSGYDDRYDWFFPX7JKVPC`

## reasoning

Intent: milestone-29's statement is one root daemon building every project under one run cap on one machine, and this failure came from ~5 concurrent agents on the 4-core VM under that daemon; the previous fix of this kind served milestone-29 directly, and milestone-30 is about README, docs and MIT. Linking a milestone also keeps the change set schedulable while SkipNoMilestone (internal/api/daemon.go:1181) still skips milestone-less work. Stopping it: the asker searched gosweep (TempMaxAge 2h on top-level mtime, internal/gosweep/gosweep.go:19), cron, timers, tmpfiles and transcripts and found nothing; I also checked that the daemon's sweep under the test clock (2026-10-24, before real time) cannot judge a fresh folder stale, and that daemon-trims-go-cache.txtar points TMPDIR at $WORK. Without a culprit, a further hunt has no clear end, while a per-run GOTMPDIR outside TMPDIR removes the whole class of shared-TMPDIR deleters at once; the runs directory (run.Dir, internal/run/store.go:38-40) is shared by every run and holds one <run id>.jsonl journal each, with no per-run folder, so the leaf adds a <run id>/ folder beside the journal; Store's listing (internal/run/store.go:178-191) skips directories and non-.jsonl names, so that folder is never read as a run, and it sits on the repo's filesystem where nothing sweeps it. Retry: step commands are arbitrary shell from workflows (stepRunner.Command, internal/api/run.go:1690) and the engine is generic for any repo, so a Go-specific signature belongs beside the runner that already captures the output tail, not in internal/run; retrying exactly once matches the change set's intent and avoids masking a real, repeating failure.

## statement

(1) The change set serves milestone-29, like the last bug fix (Fix Review question handles): the failure is the root daemon building several runs at once on one machine, which is milestone 29's promise; milestone 30 (sharing Rhiza publicly) has nothing to do with it. (2) With no deleter found, stop it by isolation rather than by hunting: every process a run starts (command, wait and agent steps) gets GOTMPDIR set to a folder of that run's own outside the shared TMPDIR, at <git common dir>/rhiza/runs/<run id>/gotmp: a new <run id>/ folder the leaf creates beside the run's <run id>.jsonl journal in the runs directory every run shares (run.Dir), created when a step starts and removed when the run ends; nothing that sweeps ~/tmp by prefix or age (gosweep, tmpsweep, a test's TestMain, an agent's cleanup) can reach a live $WORK there. gosweep keeps sweeping TMPDIR unchanged. (3) The retry lives in the step runner's Command (internal/api/run.go stepRunner.Command), not in the run engine: when a command step exits non-zero and its output matches go's vanished-$WORK signature (open $WORK/<bNNN>/... : no such file or directory), it reruns the command once, writing one stderr line that says so, and the second result stands. The signature match is a small pure function with unit tests; the engine stays language-agnostic.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-457efc88-1](../question/question-457efc88-1.md)
- derived-from [go-work-vanish-plan](../source/go-work-vanish-plan.md)

### In

- decided-by [run-private-gotmpdir](../capability/rhiza/workflows/workflow-runs/run-private-gotmpdir.md)
- decided-by [step-retry-vanished-go-work](../capability/rhiza/workflows/workflow-runs/step-retry-vanished-go-work.md)

## History

- `01M6WGEPSYWV8R14VGKJX9GGJM` Fix A Go build inside a run step fails when its $WORK temp folder vanishes mid-build, and the run ends: plan run 01M6W7G5J0M6 (change set 01M6W7G8YBXB) failed at step rebase on 2026-11-08 ~00:15Z with 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' for many std packages while ~5 agents built and tested on the 4-core VM (load 28). The daemon's go build sweep (gosweep, TempMaxAge 2h) logged nothing in that window, so something else removed a live ~/tmp/go-build* folder: find what (candidates: a test's TestMain or tmpsweep removing TMPDIR entries by prefix, an agent's cleanup, go's own cleanup under a shared TMPDIR) and stop it. Also: a step whose command fails with this signature should retry once instead of ending the run. Unstuck by rhiza run resume --retry 01M6W7G5J0M6.
