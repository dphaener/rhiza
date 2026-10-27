# plan-decisions-from-source

capability · `01M5WSGCY2WM4E3XGQ2T8R1ET8` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A testscript adds a Source with --raw-session for a fake session whose transcript holds one picked and one typed answer, then runs rhiza session interview on it: both proposed Decisions link derived-from the Source, the picked one with its wording and the typed one with its quote; with no such Source both are proposed with no derived-from link, as now. Go guide tests find the derived-from and typed-words rules in rhiza guide plan.

## statement

Every Decision a plan settles links derived-from the plan's Source: the Decisions an interview's capture proposes link derived-from the Source in the session's change set whose session field names that session, and the plan text tells the planner that each Decision it writes from a typed answer quotes the owner's typed words and links derived-from that Source.

## Links

### Out

- depends-on [plan-source](plan-source.md)
- decided-by [plan-workflow-shape](../../../../../decision/plan-workflow-shape.md)

### In

- verifies [plan-decisions-from-source-verified](../../../../../evidence/plan-decisions-from-source-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
