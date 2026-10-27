# plan-source

capability · `01M5WSGCSYNYA9BFQX50REASAY` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A testscript launches a fake agent and finds RHIZA_SESSION in its environment; rhiza source add --raw-session on a session with a fake transcript gives raw_sha256 equal to the transcript file's sha256 and a session field naming the session, and the transcript's text appears nowhere under .rhiza/; an unknown session exits 4, and --raw-session with --raw or --raw-sha256 exits 2. Go guide tests find --raw-session "$RHIZA_SESSION" in rhiza guide plan. describe-golden is regenerated.

## statement

A planner records its interview as a Source in the plan's change set without storing the transcript: every agent Rhiza launches has RHIZA_SESSION set to its session ID, and rhiza source add --raw-session <session> hashes that session's transcript, found through its agent adapter, as the raw material and gives the Source a session field naming it. The plan text tells the planner to scrub its interview and add it this way before it exits.

## Links

### Out

- depends-on [plan-guide](plan-guide.md)
- decided-by [plan-workflow-shape](../../../../../decision/plan-workflow-shape.md)

### In

- depends-on [plan-decisions-from-source](plan-decisions-from-source.md)
- verifies [plan-source-verified](../../../../../evidence/plan-source-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
