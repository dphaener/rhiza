# decision-paging-log-roadmap-why

decision · `01M6H3FRAAG0RAS8HVZ9YRJ0P1`

## reasoning

The change set's intent is to page log, roadmap and why 'as ls already does', so I copied ls and feed: limit default 100, max 1000, opaque base64url cursor, next_cursor absent at the end, a next: line in text, usage errors for bad limits and cursors. For the escape hatch the repo already has one precedent, run log's --all ('ignoring --limit and --cursor'), and run log's Lua pane runs 'run log --all', so --all is the shared way to read everything. roadmap's groups (internal/api/plan.go RoadmapGroup) are few but their change sets are what grows, so the page unit is the change set and groups are kept as framing, repeated when split; a change set can serve more than one milestone and so appear in several groups, which is why the cursor must name the group as well as the ID. why (internal/api/lineage.go WhyResponse) needs the full walk to know walk order and quote cutting anyway, so paging after the walk keeps its output identical to today's when concatenated; one combined sequence keeps one limit and one cursor instead of three. The callers that need everything (healthRoadmap in internal/api/health.go, internal/tui/roadmapapp and the modelapp roadmap tab) would otherwise read only the first 100 change sets and report a wrong roadmap, so they ask for all explicitly.

## statement

log, roadmap and why page as ls and feed do, and share one shape. Each request gains limit (1 to 1000, default 100; out of range is a usage error), cursor (opaque base64url, a usage error when it is not one) and all (flag --all: return everything, ignoring limit and cursor, as run log --all does). Each response gains next_cursor, absent when no more remain, and the text output ends with a 'next: <command> --cursor <c>' line as ls prints. log pages its change_sets in apply order; target and view come on every page. roadmap pages over its change sets taken in group order (milestones in intent-tree order, then backlog, then umbrellas; within a group depends-on order): limit counts change sets, a page holds only the groups that have change sets on it, a group split across pages repeats its kind and intent on the next page, and the cursor names the group and the last change set returned so a change set that sits in two groups is not confused. why computes its whole walk, then pages the three lists as one sequence: decisions, then sources, then change sets, each in its own order; limit counts entries across all three, node and view come on every page, and the cursor names the list and the offset within it. The in-repo callers that need the whole roadmap ask for it with all: health's roadmap check calls runRoadmap with RoadmapRequest{All: true}, and the roadmap TUI app and the model app's roadmap tab pass all, so none of them silently truncates at 100.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-af1aebc5-1](../question/question-af1aebc5-1.md)
- derived-from [plan-paging-log-roadmap-why](../source/plan-paging-log-roadmap-why.md)

### In

- decided-by [log-paged](../capability/rhiza/views/paged-lists/log-paged.md)
- decided-by [roadmap-paged](../capability/rhiza/views/paged-lists/roadmap-paged.md)
- decided-by [why-paged](../capability/rhiza/views/paged-lists/why-paged.md)

## History

- `01M6H3A65RDP761HKE82STZSQZ` log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan)
