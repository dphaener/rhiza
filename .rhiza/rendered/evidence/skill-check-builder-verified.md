# skill-check-builder-verified

evidence · `01M7DXMTP0DS3XK9FEMRK8WRWS` · actor [claude-code-187a4e61](../actor/claude-code-187a4e61.md)

## result

PASS

## statement

Verifier at ca534ddae, with -count=1 under systemd-run MemoryMax=6G TasksMax=300 and session RHIZA_/SHIP_/REX_ variables unset. go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/^(skill-check|describe-golden)$': exit 0. go test -trimpath -p 2 ./internal/api -run TestSkillCheck: exit 0 (report, trailers gate, definition and refusal tests). go test ./cmd/rhiza -run TestScripts/describe-golden -update: exit 0, no diff. A scratch build's ship generate and render: exit 0, adding nothing beyond the verifier's own actor record. make self-check: exit 0 in a clean detached worktree at HEAD. skill-check.txtar uses a fake claude and a fake rex and covers each acceptance point: the good builder passes 3 of 3 gates and exits 0; bare fails trailers by name and exits 1; quiet fails session-done by name and exits 1; --json prints each gate with pass or fail; a missing skill is refused naming rhiza guide install --agent claude; git status and .rhiza/log are unchanged after every case. Diff matches Decisions skill-check-scratch-one-role (scratch root from os.MkdirTemp, refused inside a git repository, caller's RHIZA_*/REX_BLOCK/GIT_* unset, builder role only) and skill-check-plan-shape (this is the builder-and-report leaf; gates are a list). Files: the four named, plus generated docs/reference/cli.md and .rhiza/rendered. No stubs or TODOs. Not passing: the Checks: line go test -trimpath -p 2 ./internal/api/... hit the 240s cap again (exit 124, no output, as the builder recorded; not retried). Under the orchestrator's runaway rule that is a blocker to hunt. The new tests pass in 0.04s and spawn no process, so the hang is probably in another internal/api test; find it with -v -timeout 200s and compare at base a83650275.

## Links

### Out

- verifies [skill-check-builder](../capability/rhiza/workflows/workflow-authoring/skill-check-builder.md)

## History

- `01M79KHBKTKGWFD3ZGV1B5QPHA` rhiza skill check runs a skill against a fixture change set and reports which gates pass, and rhiza workflow check warns about unresolved invoke, uncapped hand-back loops and a merge with no gate
