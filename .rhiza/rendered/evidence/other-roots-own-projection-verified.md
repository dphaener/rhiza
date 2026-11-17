# other-roots-own-projection-verified

evidence · `01M7K619MFNY8V5NCCSTTYK7A4` · actor [claude-code-972e5ce6](../actor/claude-code-972e5ce6.md)

## result

PASS

## statement

Verified commit 78c179d91. The acceptance grep (IndexPath = "" outside tests in internal/api) prints nothing. go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestDaemonHomePrefersWorkingProject|TestPollGitBounded|TestGateAnswer|TestInboxGate|TestWebAnswerAction|TestWebApproveAction|TestSessionBrief|TestDriverPollIsCheap|TestEvidenceReadsPastSessionWrites|TestOtherRoot)' passes, TestOtherRoot{DaemonProject,InboxGateRead,GateAnswer,Unreadable} included. Each runs with an own-index case and a throwaway case. In the own-index case the worktree's index file exists under the test HOME's .cache, and the answers are as expected: ls lists other-plan-intent; inbox statements; question.answer --run writes to the worktree, not the caller; gateAnswered goes false then true. When the root is unreadable the result is unchanged: empty statements, or the no Rhiza repository error. TestPropertyProjectionParity passes. The scripts daemon-projects, inbox-gate-questions and inbox-projects pass, with RHIZA_/SHIP_/REX_ unset. go vet ./internal/api is clean. rhiza ship generate and rhiza render, run with this commit's build in a detached worktree, changed nothing. In the diff, the one rootIndexPath helper (workspace.go) feeds projectCtx, gateAnswerCtx, inboxGateRead, run_brief, run_model_cache and openRootWorkspace. gateAnswered no longer clears IndexPath. There is no log-fold fallback and no TODOs, and the root each op reads is unchanged.

## Links

### Out

- verifies [other-roots-own-projection](../capability/rhiza/views/fast-reads/other-roots-own-projection.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
