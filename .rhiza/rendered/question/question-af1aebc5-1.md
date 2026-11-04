# question-af1aebc5-1

question · `01M6H3DJABRP3Z865V3CGPVNPB`

## session

af1aebc5-0e9d-4d99-8461-01a1a2353954

## statement

How should `rhiza roadmap` page? Its output is groups (milestone, backlog, umbrella) of change sets. / How should `rhiza why` page? It returns three lists: decisions, sources and change sets, in walk order. / What limits and escape hatch should log, roadmap and why share? / In-repo callers want the whole roadmap: rhiza health (roadmap check), the roadmap TUI app and the model app's roadmap tab. How should they read it once roadmap pages?

## Links

### In

- answers [decision-paging-log-roadmap-why](../decision/decision-paging-log-roadmap-why.md)

## History

- `01M6H3A65RDP761HKE82STZSQZ` log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan)
