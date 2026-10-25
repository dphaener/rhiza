# plan-task-4

source · `01M5SWRGQCHRVACT3HFXCR5RAN`

## quotes

```json
[]
```

## raw_sha256

d6b1324d76b6d22cb390255979e129dce05652983deb9e793a997e6e8491bbde

## summary

The shipit work plan for task 4: a pure internal/model package that projects folded change sets into three views, as-is (landed), to-be (as-is plus accepted) and what-if (as-is plus one change set), orders change sets causally, detects conflicts and never fails. Decided (D1-D13, agent wording accepted by the owner): scope is the pure package only, no CLI, index, rebase or field-level diff; AsIs, ToBe and WhatIf return a Model and problems; causal order comes from base edges and existence edges, ties by change set ID, with cycles reported and broken at the smallest ID; a conflict is a stale base that overlaps an earlier write (same field or handle, a move against a move, anything against a retire), while non-overlapping stale operations merge silently; a conflicting operation still applies, later wins, and is reported; operations that break model rules are skipped and reported; to-be and what-if use the same engine and rules; Overlaps reports field-level overlapping pairs among in-flight change sets; duplicate handles are reported, not skipped; nodes and links carry their history, version and creating change set; tests are TDD with rapid determinism and commutativity properties and the corpus projected in every view; guardrails keep other packages' APIs unchanged. Settled later in review (D14): a landed conflict counts as repaired once later operations written from the current version rewrite every aspect the conflicting operation wrote, and a conflict involving a retirement stays reported. Endpoint kinds were left to task 9's check.

## title

Shipit plan, milestone 1 task 4: projection

## Links

### In

- derived-from [capability-tree](../capability/genesis/model/capability-tree.md)
- derived-from [conflicts](../capability/genesis/change-sets/conflicts.md)
- derived-from [overlays](../capability/genesis/views/overlays.md)
- derived-from [t4-d1-pure-package](../decision/t4-d1-pure-package.md)
- derived-from [t4-d2-internal-model](../decision/t4-d2-internal-model.md)
- derived-from [t4-d3-causal-order](../decision/t4-d3-causal-order.md)
- derived-from [t4-d4-field-level-conflicts](../decision/t4-d4-field-level-conflicts.md)
- derived-from [t4-d5-conflicting-op-applies](../decision/t4-d5-conflicting-op-applies.md)
- derived-from [t4-d6-skip-broken-ops](../decision/t4-d6-skip-broken-ops.md)
- derived-from [t4-d7-one-engine-views](../decision/t4-d7-one-engine-views.md)
- derived-from [t4-d8-overlaps](../decision/t4-d8-overlaps.md)
- derived-from [t4-d9-duplicate-handles](../decision/t4-d9-duplicate-handles.md)
- derived-from [t4-d10-node-history](../decision/t4-d10-node-history.md)
- derived-from [t4-d11-tdd-properties-corpus](../decision/t4-d11-tdd-properties-corpus.md)
- derived-from [t4-d12-task-4-guardrails](../decision/t4-d12-task-4-guardrails.md)
- derived-from [t4-d13-task-4-acceptance](../decision/t4-d13-task-4-acceptance.md)
- derived-from [t4-d14-conflict-repair](../decision/t4-d14-conflict-repair.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
