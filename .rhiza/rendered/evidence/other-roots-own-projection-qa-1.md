# other-roots-own-projection-qa-1

evidence · `01M7KB21D2KBB8Y390K7Q2B3Y5` · actor [claude-code-49a7d3f2](../actor/claude-code-49a7d3f2.md)

## result

PASS

## statement

QA: a root rhiza.toml registering project alpha (a copy of the scratch repo, with two overlapping change sets and a blocking question). rhiza inbox --json at the root tagged alpha's items and was byte-identical to the base build's and to a run with no index path; alpha's own index file was created under the test HOME's cache. question answer color-q into a change set (both with and without an index) wrote the Decision and cleared the blocker items as the base build does; a second answer is refused 'not an open question blocking this change set', an unknown question 'no question' and an unknown --run 'no run matches'. Scoped api tests (TestDaemonHomePrefersWorkingProject|TestPollGitBounded|TestGateAnswer|TestInbox|TestWebAnswerAction|TestWebApproveAction|TestSessionBrief|TestDriverPollIsCheap|TestEvidenceReadsPastSessionWrites) ok; scripts daemon-projects|inbox-gate-questions|inbox-projects ok. The IndexPath = "" grep prints nothing.

## Links

### Out

- verifies [other-roots-own-projection](../capability/rhiza/views/fast-reads/other-roots-own-projection.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
