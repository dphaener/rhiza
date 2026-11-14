# classic-tests-pinned

capability · `01M7AHXFP07DFAG8GEEEHRXJ5S` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-headless-runs|run-log|run-ship-change|run-ship-change-rebase|run-ship-ci-fix|run-ship-ci-lint-fix|run-ship-conflict-fix|run-ship-deploy-check|run-ship-deploy-timeout|run-ship-land-conflict|run-ship-main-dispatch|run-ship-main-dispatch-fail|run-ship-main-dispatch-resume|run-ship-model-conflict|run-ship-rerun-fail|run-ship-rerun-pass|daemon-max-runs-pending|daemon-milestone-request|daemon-notify|daemon-notify-hook|daemon-rex-session|inbox-blocker-advice|inbox-blocking-questions|inbox-landed-resumable|inbox-list|inbox-open|inbox-permission|inbox-projects|inbox-signoff|inbox-stale)$' go test -trimpath ./internal/api -run 'TestPrePushJournal|TestInboxReadsChangeSetsOnce' and go test -trimpath ./internal/workflow -run 'TestShipChangeBuiltin|TestPlanRerunsCIFlakes|TestShipChangeLocalGate|TestShipChangeModelConflictFix|TestShipChangeLandedInfraWaitScript|TestShipFixPromptsOneLine|TestShipChangeMergeRefusedStops|TestMaxIntervalShipChangeCIWaits' pass with ship-change still the classic flow, and grep -L ship-change-classic run on those 30 scripts and on internal/api/ship_prepush_test.go, internal/api/inbox_fold_test.go and internal/workflow/workflow_test.go prints nothing, and inside the bodies of those eight workflow tests no Builtin("ship-change") or load("ship-change") remains (a grep -n of that pattern on internal/workflow/workflow_test.go prints only lines within TestBuiltinShipChange, which builtin-conductor-default adds and which reads the new builtin). The first 16 scripts start a run on the default ship-change (rhiza run start with no --workflow, or the daemon's ready pick in daemon-headless-runs) and assert its steps; each passes --workflow ship-change-classic, or a .rhiza/dispatch.toml with [ready] feature = "ship-change-classic" for the daemon. The other 14 seed a run log whose started entry names workflow ship-change at a classic step; each names ship-change-classic there. TestPrePushJournal and TestInboxReadsChangeSetsOnce seed run.Started{Workflow: "ship-change"}; each names ship-change-classic. The eight workflow tests TestShipChangeBuiltin, TestPlanRerunsCIFlakes, TestShipChangeLocalGate, TestShipChangeModelConflictFix, TestShipChangeLandedInfraWaitScript, TestShipFixPromptsOneLine, TestShipChangeMergeRefusedStops and TestMaxIntervalShipChangeCIWaits parse the builtin ship-change and assert classic steps (infra waits, ship check steps, model-conflict fix, merge refused); each loads ship-change-classic (Builtin("ship-change-classic") or load("ship-change-classic")) instead. Scripts that copy their own .rhiza/workflows/ship-change.toml (run-ship-infra-wait, run-ship-landed-infra-wait, run-ship-retry-ci-triage, run-ship-retry-pending, run-ship-merge-draft, run-ship-merge-main-moved, run-ship-change-fail, headless-run-agent-tabs, plan-workflow-ci-infra-timeout, plan-workflow-pr, run-records-dispatch, run-ship-needs-leaves) shadow the builtin and are left as they are; leafless-refusal-by-leaves-step owns run-ship-needs-leaves. Touches exactly those 30 scripts, internal/api/ship_prepush_test.go, internal/api/inbox_fold_test.go and internal/workflow/workflow_test.go (only those eight tests; builtin-conductor-default adds TestBuiltinShipChange to the same file).

## statement

Every existing test that drives the classic ship-change's own steps (land, rebase, conflict, CI, gate and infra loops, sign-off, merge, deploy) names ship-change-classic, by --workflow or a dispatch.toml in the script, so it keeps covering the ship machinery after the default changes; no such test is deleted.

## Links

### Out

- depends-on [builtin-ship-change-classic](builtin-ship-change-classic.md)
- depends-on [leafless-refusal-by-leaves-step](leafless-refusal-by-leaves-step.md)
- decided-by [builtin-ship-change-open-points](../../../../decision/builtin-ship-change-open-points.md)

### In

- depends-on [builtin-conductor-default](builtin-conductor-default.md)
- verifies [classic-tests-pinned-verified](../../../../evidence/classic-tests-pinned-verified.md)
- verifies [classic-tests-pinned-qa-2](../../../../evidence/classic-tests-pinned-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
