# t8-d16-json-shape

decision · `01M5SWSH46GYVEFFY59CV5XJ3E`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Sections of typed items (Recommended): {view, change, budget, tokens, focus:[...], context:[...], ..., omitted:[{section, count, next}]}. Each item has full id, short, handle, kind, plus tier fields (statement, note, from, status); focus items carry every field and the proposed change. Shape pinned in describe-golden as a contract.

## statement

JSON: `{view, change, ref, budget, tokens, over_budget, focus, context, constraints, decisions, questions, anchors, evidence, children, neighbours, in_flight, omitted}`. Items are typed: id, short, handle, kind, statement, plus note, from, status and relation where they apply. Focus items carry every field and the proposed change. `omitted` is `[{section, count, tier, next}]`. Pinned by describe-golden.

## wording

agent of milestone-1 task 8 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-8](../source/plan-task-8.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [slice](../capability/rhiza/views/slice.md)

## History

- `01M5SWSE0Z4QGWQGN2FC9WMVAW` Task 8's plan decisions (PR #9, merge a951e38), accepted by the owner
