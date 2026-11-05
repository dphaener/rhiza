# inbox-lists-owner-run-proposals-verified

evidence · `01M6PW6JTRJG651GZ0NP549ZW7` · actor [claude-code-bba5fe8b](../actor/claude-code-bba5fe8b.md)

## result

PASS

## statement

Verifier claude-code-bba5fe8b at fb81edc3: read the diff (internal/api/inbox.go drops the st.Autonomous condition so any plan run's Gated change sets still proposed list as accept items, unless a newer run retried the change set) and the new cmd/rhiza/testdata/script/inbox-owner-run-accept.txtar, which sets rhiza.autonomy.accept owner, gates a non-autonomous plan run (asserts "autonomous": false), shows rhiza inbox --json listing kind accept with detail 'waiting for acceptance' and next 'rhiza accept --rule owner --change <id>', then after rhiza accept --rule owner shows no accept item, and passes fold-ok. Ran go test -trimpath ./cmd/rhiza -count=1 -run TestScripts/(inbox-.*|owner-run-gate-accept|plan-workflow-autonomous.*|run-gate|autonomous-gate-rule|tui-inbox-app|autonomy-accept-rule|stage-autonomy): ok. stage-autonomy still lists the autonomous plan's accept items under accept owner. go test -trimpath ./internal/api -run Inbox: ok. (The whole internal/api package hit go test's 10m timeout once while the machine ran several other suites at the same time; the inbox tests pass when run on their own.)

## Links

### Out

- verifies [inbox-lists-owner-run-proposals](../capability/rhiza/orchestration/inbox/inbox-lists-owner-run-proposals.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
