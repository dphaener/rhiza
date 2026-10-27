# plan-guide-install

capability · `01M5WSGCA14V3HZZFPJ4A1WXVB` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A testscript runs rhiza guide install with its own HOME: the plugin holds skills/rhiza-plan/SKILL.md whose frontmatter names rhiza-plan with the plan text's description and whose body is rhiza guide plan's text; after that file is deleted the status says the install is not current, and a second install restores it.

## statement

rhiza guide install also installs the plan text as the rhiza-plan skill, and its status reports the install as not current when that skill is missing or out of date.

## Links

### Out

- depends-on [plan-guide](plan-guide.md)
- decided-by [plan-workflow-shape](../../../../../decision/plan-workflow-shape.md)

### In

- depends-on [plan-workflow-file](plan-workflow-file.md)
- verifies [plan-guide-install-verified](../../../../../evidence/plan-guide-install-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
