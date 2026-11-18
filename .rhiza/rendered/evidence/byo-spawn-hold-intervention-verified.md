# byo-spawn-hold-intervention-verified

evidence · `01M7RDADF60FE8Q8Z91XTWEJKR` · actor [claude-code-ef6b593a](../actor/claude-code-ef6b593a.md)

## result

PASS

## statement

At 38f0f597e, env -u RHIZA_*/SHIP_*/REX_* go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/byo-workflow-hold-intervention$' -count=1 passed (exit 0), as did TestDocsExampleBYOWorkflowCopies|TestProcessScripts|TestProcessListConverted|TestScriptOrder|TestScriptEnv; rhiza ship generate and rhiza render exit 0. The script holds myteam-flow's work step at RHIZA_TEST_LOAD1/5=3.5 (one spawn-held for step work in rhiza inbox --json and text, zero session journal entries, no claude call), releases at low load via daemon --once (the session starts, inbox clear), asserts intervention_count 0 after the session's own message 'Own note' (failure case), then 1 message intervention naming actor test, step work and the text after an outside rhiza session message; the run then succeeds and lands. The embedded myteam-flow.toml is byte-identical to docs/examples/byo/workflow.toml; the commit changes no docs/ or internal/ file and no Go code names myteam or work; processScripts entry added with its reason.

## Links

### Out

- verifies [byo-spawn-hold-intervention](../capability/rhiza/workflows/workflow-runs/byo-spawn-hold-intervention.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
