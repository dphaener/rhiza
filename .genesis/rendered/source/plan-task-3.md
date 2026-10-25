# plan-task-3

source · `01M5SWRGPFD3NEBM2VR9ARZXZE`

## quotes

```json
[]
```

## raw_sha256

b60de4f8c0a964c5d1eb3d1fded6f5319da91db66328af52704b51345426e35e

## summary

The shipit work plan for task 3: typed entry bodies, the six operations, folding a change set's entries in causal order into its current operations and lifecycle state, and fractional sibling keys, as pure packages (internal/key, op, changeset). Decided (D1-D16, agent wording accepted by the owner): pure domain packages only, no CLI, projection or index; each operation keeps one stable ID, a later ops entry replaces it by repeating the ID and drops operations through a drop list, and a replacement keeps its position, op type and target, a dropped ID never returning; ops entries are allowed in draft, proposed and accepted without resetting the state; a record change set lands from draft, proposed or accepted, an intention change set only from accepted; folding orders entries causally with ties by entry ID, skips and reports illegal entries and concurrent lifecycle entries, and never fails; bodies are strict like the envelope and the corpus is amended to match; node fields are free snake_case names to any JSON value, per-kind requirements deferred, the handle one slug segment; nine node kinds, with change sets not nodes but allowed as link endpoints; accept and reject bodies are {rule, reason}, the reason may be empty; tests are TDD with table tests, fuzzed body decoding, rapid properties and a clean corpus fold; guardrails keep the envelope and store APIs, add no dependencies and no filesystem or git access; no CONTEXT.md glossary. Per-kind field schemas were left for later.

## title

Shipit plan, milestone 1 task 3: change sets

## Links

### In

- derived-from [node-kinds](../capability/genesis/model/node-kinds.md)
- derived-from [typed-links](../capability/genesis/model/typed-links.md)
- derived-from [handles](../capability/genesis/model/handles.md)
- derived-from [six-operations](../capability/genesis/model/six-operations.md)
- derived-from [change-lifecycle](../capability/genesis/change-sets/change-lifecycle.md)
- derived-from [revisions](../capability/genesis/change-sets/revisions.md)
- derived-from [t3-d1-pure-packages](../decision/t3-d1-pure-packages.md)
- derived-from [t3-d2-stable-op-ids](../decision/t3-d2-stable-op-ids.md)
- derived-from [t3-d3-revise-keeps-state](../decision/t3-d3-revise-keeps-state.md)
- derived-from [t3-d4-record-landing](../decision/t3-d4-record-landing.md)
- derived-from [t3-d5-deterministic-fold](../decision/t3-d5-deterministic-fold.md)
- derived-from [t3-d6-strict-bodies](../decision/t3-d6-strict-bodies.md)
- derived-from [t3-d7-free-fields-slug-handle](../decision/t3-d7-free-fields-slug-handle.md)
- derived-from [t3-d8-nine-node-kinds](../decision/t3-d8-nine-node-kinds.md)
- derived-from [t3-d9-fractional-keys](../decision/t3-d9-fractional-keys.md)
- derived-from [t3-d10-tdd-fuzz-properties](../decision/t3-d10-tdd-fuzz-properties.md)
- derived-from [t3-d11-task-3-guardrails](../decision/t3-d11-task-3-guardrails.md)
- derived-from [t3-d12-no-context-md](../decision/t3-d12-no-context-md.md)
- derived-from [t3-d13-three-packages](../decision/t3-d13-three-packages.md)
- derived-from [t3-d14-accept-reject-bodies](../decision/t3-d14-accept-reject-bodies.md)
- derived-from [t3-d15-replacement-rules](../decision/t3-d15-replacement-rules.md)
- derived-from [t3-d16-task-3-acceptance](../decision/t3-d16-task-3-acceptance.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
