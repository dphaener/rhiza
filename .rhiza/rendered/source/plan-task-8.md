# plan-task-8

source · `01M5SWRGVE8N122VSRM8K0ATGP`

## quotes

```json
[]
```

## raw_sha256

cb9af378db81a3304ba89641380a1d50539fcac6e7ca856e3cbe4bac252889fd

## summary

The shipit work plan for task 8: genesis slice, the composed context an agent reads for a node or a change set, in the ten agreed sections, with inheritance down the tree and a token budget enforced by fixed cuts, as a pure internal/slice package over a reader fed by the index or a what-if model. Decided (D1-D17, agent wording accepted by the owner): slice takes a ref for a node slice, or no ref for a change slice of the effective change set; as-is and to-be read the index and what-if is computed, with the 50 ms budget covering only as-is and to-be; the budget measures the bytes actually printed, in markdown or JSON, with tokens as ceil(bytes / 4); constraints apply only through constrains links, a constraint on a tree root applies below it marked with its root, and an unlinked constraint applies nowhere; short IDs are the shortest unique prefix of at least six characters; a change slice touches its operations' target nodes and the node ends of the links it creates or retires; items read the statement field, and evidence status is derived from verifies links; cuts go full, then brief, then count, from section 10 upward, never touching Focus, Constraints or Open questions, with a footer naming the command for what was cut; current decisions are live and not superseded, open questions are unanswered, blocking ones first; neighbours and in-flight sections; over-budget output prints everything required, notes the overage and exits 0; a typed JSON shape pinned by the describe golden; TDD with goldens and a property that index-built equals model-built. Revised by the orchestrator under the owner's authorization: D1's default view now follows the rule shared with task 7. Left open: a fast incremental what-if.

## title

Shipit plan, milestone 1 task 8: slice

## Links

### In

- derived-from [slice](../capability/rhiza/views/slice.md)
- derived-from [t8-d1-slice-shape](../decision/t8-d1-slice-shape.md)
- derived-from [t8-d2-budget-stored-views](../decision/t8-d2-budget-stored-views.md)
- derived-from [t8-d3-budget-printed-bytes](../decision/t8-d3-budget-printed-bytes.md)
- derived-from [t8-d4-token-estimate](../decision/t8-d4-token-estimate.md)
- derived-from [t8-d5-constraints-via-links](../decision/t8-d5-constraints-via-links.md)
- derived-from [t8-d6-short-ids](../decision/t8-d6-short-ids.md)
- derived-from [t8-d7-slice-package](../decision/t8-d7-slice-package.md)
- derived-from [t8-d8-change-slice-focus](../decision/t8-d8-change-slice-focus.md)
- derived-from [t8-d9-item-fields](../decision/t8-d9-item-fields.md)
- derived-from [t8-d10-cut-order](../decision/t8-d10-cut-order.md)
- derived-from [t8-d11-tests](../decision/t8-d11-tests.md)
- derived-from [t8-d12-guardrails](../decision/t8-d12-guardrails.md)
- derived-from [t8-d13-current-decisions](../decision/t8-d13-current-decisions.md)
- derived-from [t8-d14-neighbours-in-flight](../decision/t8-d14-neighbours-in-flight.md)
- derived-from [t8-d15-over-budget](../decision/t8-d15-over-budget.md)
- derived-from [t8-d16-json-shape](../decision/t8-d16-json-shape.md)
- derived-from [t8-d17-proposed-change](../decision/t8-d17-proposed-change.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
