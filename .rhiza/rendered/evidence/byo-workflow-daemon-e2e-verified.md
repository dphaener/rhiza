# byo-workflow-daemon-e2e-verified

evidence · `01M7RCP18C7F2HZVYBPHZHEV2G` · actor [claude-code-d48f4d38](../actor/claude-code-d48f4d38.md)

## result

PASS

## statement

At 9b99a6597 (capped scope): env -u RHIZA_*/SHIP_*/REX_* go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/byo-workflow-daemon$' -count=1 -v: PASS (script ran). It asserts that daemon --once logs run start --workflow myteam-flow --pick rule ready.feature (not ship-change), that the run is succeeded with every step in order (gate visit 1 fail then work, gate visit 2 pass, signoff, push, wait-ci, merge, close-out, no visit 3), that intervention_count is 0, that the fake docker logs 'compose up -d --wait' then 'compose down' and the stop-services cleanup passes, that make.log shows 'lint fail' then the second gate's 'lint pass'/'test pass', and that handback.md names make lint. go test -run TestDocsExample: PASS. In a scratch worktree with the archived copy's handback_cap edited, TestDocsExampleBYOWorkflowCopies fails, naming byo-workflow-daemon.txtar (exit 1). The archived copy diffs byte-identical to docs/examples/byo/workflow.toml. docs/examples/byo and internal/ are untouched by the commit, and no non-test Go file names myteam. gofmt -l cmd/rhiza is clean. The processScripts entry has its reason.

## Links

### Out

- verifies [byo-workflow-daemon-e2e](../capability/rhiza/workflows/workflow-runs/byo-workflow-daemon-e2e.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
