# plan-task-7

source · `01M5SWRGTD75VQX2D5JHRJF0Q4`

## quotes

```json
[]
```

## raw_sha256

9494cb7e8d26d6d93a0fe162851eccdad16894649d1b09eb74192ad271edf66d

## summary

The shipit work plan for task 7: the read commands show, tree, ls, why, log, roadmap and ready, each with --view and --json, served from the index. Decided (D1-D17, agent wording accepted by the owner): scope is the seven reads, with ls paged and find, change diff, events and paging elsewhere out; reads use the index only, with new queries where needed, and writes stay untouched; shared read plumbing in internal/api/read.go; why walks decided-by (including inherited from tree ancestors), derived-from and supersedes both ways, shows answered questions, then the change sets in the node's history, cycle-guarded; log lists the view's history of a node or link in apply order; a change set's links for roadmap and ready are the to-be links plus the link operations inside that change set; ready means accepted, every depends-on landed, and no open question blocking the change set, a node it touches or an ancestor, and it names the first ready one as next; roadmap groups proposed and accepted change sets by served intent in intent-tree order, then depends-on order, ties by ID, with a final no-intent group; tree, ls and show shapes; tests with each task; guardrails; finish at an open PR with green CI. Revised by the orchestrator under the owner's authorization: the default view (D2) became one rule shared with slice (draft or proposed change set gives what-if, accepted gives to-be, otherwise as-is), with roadmap and ready always defaulting to to-be (D17); and the read budget (D14) went from a 25 ms median to 50 ms, the milestone target, after CI measured about 25.5 ms, the 25 ms figures staying as benchmarks. Left open: speeding up roadmap and ready, and cursor paging on log, roadmap and why.

## title

Shipit plan, milestone 1 task 7: the reads

## Links

### In

- derived-from [why-walk](../capability/rhiza/lineage/why-walk.md)
- derived-from [change-history](../capability/rhiza/lineage/change-history.md)
- derived-from [node-show](../capability/rhiza/views/node-show.md)
- derived-from [tree-view](../capability/rhiza/views/tree-view.md)
- derived-from [node-list](../capability/rhiza/views/node-list.md)
- derived-from [roadmap](../capability/rhiza/views/roadmap.md)
- derived-from [ready-next](../capability/rhiza/views/ready-next.md)
- derived-from [paged-lists](../capability/rhiza/views/paged-lists.md)
- derived-from [t7-d1-reads-scope](../decision/t7-d1-reads-scope.md)
- derived-from [t7-d2-default-view](../decision/t7-d2-default-view.md)
- derived-from [t7-d3-reads-use-index](../decision/t7-d3-reads-use-index.md)
- derived-from [t7-d4-shared-read-plumbing](../decision/t7-d4-shared-read-plumbing.md)
- derived-from [t7-d5-why-walk](../decision/t7-d5-why-walk.md)
- derived-from [t7-d6-log-history](../decision/t7-d6-log-history.md)
- derived-from [t7-d7-change-set-links](../decision/t7-d7-change-set-links.md)
- derived-from [t7-d8-ready-rule](../decision/t7-d8-ready-rule.md)
- derived-from [t7-d9-roadmap-grouping](../decision/t7-d9-roadmap-grouping.md)
- derived-from [t7-d10-tree-shape](../decision/t7-d10-tree-shape.md)
- derived-from [t7-d11-ls-filters-paging](../decision/t7-d11-ls-filters-paging.md)
- derived-from [t7-d12-show-shape](../decision/t7-d12-show-shape.md)
- derived-from [t7-d13-tests](../decision/t7-d13-tests.md)
- derived-from [t7-d14-read-budget](../decision/t7-d14-read-budget.md)
- derived-from [t7-d15-guardrails](../decision/t7-d15-guardrails.md)
- derived-from [t7-d16-finish](../decision/t7-d16-finish.md)
- derived-from [t7-d17-plan-reads-to-be](../decision/t7-d17-plan-reads-to-be.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
