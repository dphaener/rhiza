# conduct-handback-topic-verified

evidence · `01M7MKXZPK0FYKBEFTDYSZJ266` · actor [claude-code-9aa90c1e](../actor/claude-code-9aa90c1e.md)

## result

PASS

## statement

Verified commit 71b326d26 for conduct-handback-topic. Re-ran its Checks: go test -p 2 -trimpath ./internal/guide -run 'TestConductHandback|TestConduct|TestList' (TestConductHandback, TestConductReadsHandback, TestConductHandbackCatchesDrift, TestListOrder pass), go test -count=1 ./internal/guide/... ok, go test -count=1 ./cmd/rhiza -run TestScripts/guide ok (RHIZA_/SHIP_/REX_ unset), go vet ./internal/guide, rhiza ship generate and rhiza render all exit 0. Diff read: conduct-handback.md names handback.md's parts as internal/run/brief.go RenderHandback writes them (Failing step, Outcome, Failing commands with output in <path>, Output tail, Cause, Failed CI runs, Conflicted files), routes for gate (its four commands from ship-change.toml), signoff, push, wait-ci (no role, no ci-fix) and merge (ci-fix, rhiza ship rebase --no-push, conflict-fix), land/render/session done, handback_cap 3 and the post-merge steps that never hand back. The test parses the built-in ship-change and the drift test covers each failure case (unrouted hand-back step, unnamed gate command, routes for main-ci/close-out/ship/deploy-check, ci-fix on wait-ci). conduct.md Read first names rhiza guide conduct-handback; built binary reports conduct 723 and conduct-handback 629 tokens, under TopicCap 1000. Only the leaf's files plus the builder's actor record changed; no TODOs or stubs.

## Links

### Out

- verifies [conduct-handback-topic](../capability/rhiza/guide/agent-guide/conduct-handback-topic.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
