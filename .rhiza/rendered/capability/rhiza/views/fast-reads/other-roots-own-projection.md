# other-roots-own-projection

capability · `01M7GFZT8RBC63C71V0CQB20XZ` · parent [fast-reads](../fast-reads.md)

## acceptance

grep -nE 'IndexPath = (""|[A-Za-z.]+, "")' internal/api/*.go | grep -v _test.go prints nothing. go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestDaemonHomePrefersWorkingProject|TestPollGitBounded|TestGateAnswer|TestInboxGate|TestWebAnswerAction|TestWebApproveAction|TestSessionBrief|TestDriverPollIsCheap|TestEvidenceReadsPastSessionWrites)' passes. A test per site (daemon project, gate answer, inbox gate read, web answer) shows the read answers from the other root's index (its index file exists afterwards under the user cache folder of the test's HOME) with the same answer as before, and one with no user cache folder answers through the throwaway index. Failure case: another root whose log does not read gives the same error or empty answer as today. go test -trimpath -count=1 -p 2 ./cmd/rhiza -run 'TestScripts/(daemon-projects|inbox-gate-questions|inbox-projects)' passes. Touches internal/api/daemon_projects.go, question.go, inbox.go (inboxGateRead), web_actions.go, run_brief.go, run_model_cache.go and their tests. Must not change which root an op reads.

## statement

An op that reads another root opens that root's own projection through one helper, never clearing IndexPath to fold that root's log: daemon_projects.go's projectCtx, question.go's gate answer (:163 and :181), inbox.go's inboxGateRead (:1414), web_actions.go's gateAnswered (:108), run_brief.go:93 and run_model_cache.go:125 take the root's default index path, or the throwaway index when it has none

## Links

### Out

- depends-on [throwaway-index-fallback](throwaway-index-fallback.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)

### In

- depends-on [interim-model-deleted](interim-model-deleted.md)
- verifies [other-roots-own-projection-verified](../../../../evidence/other-roots-own-projection-verified.md)
- verifies [other-roots-own-projection-qa-1](../../../../evidence/other-roots-own-projection-qa-1.md)
- verifies [other-roots-own-projection-qa-2](../../../../evidence/other-roots-own-projection-qa-2.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
