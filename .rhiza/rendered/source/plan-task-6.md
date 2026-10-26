# plan-task-6

source · `01M5SWRGSC1YSW7SY99DV753TK`

## quotes

```json
[]
```

## raw_sha256

7e50f7919ce365dc614bea84f2118127313f28121f3688d4a886dd8f8c1b679b

## summary

The shipit work plan for task 6: transport-independent operations with reflected JSON schemas in a new internal/api package, a thin generic CLI frame (describe, api call, --json, exit codes, next), and every milestone-1 write command over an interim model. Decided (D1-D15, agent wording accepted by the owner): writes depend on a narrow model interface, with a minimal unexported interim implementation over FoldLog until real projection is wired in; scope is the frame, init, change new/use/show/ls, add, set, mv, retire, link, unlink, apply, --record, --dry-run and the lifecycle verbs, with status, reads, --view, render and check out; operation names mirror command paths with dots and one registry serves the CLI and api call; schemas are reflected from Go types by an in-house generator, pinned by a describe golden; the actor comes from --actor or GENESIS_ACTOR, else the actor node matching git user.email, else one recorded by the first write in its own landed record change set, and the for defaults to the actor; the current change set is kept in .genesis/local/, gitignored; a reference is a full ID, a unique prefix of four or more characters, or a live handle, and handles are unique across live nodes; a write to a node the change set already touches revises the existing operation; apply reads a list of write requests and writes all or nothing; errors with --json go to stdout as {error: {code, message, next}}; accept and reject take --rule (default owner) and --reason; guardrails keep the core packages additive only and add no dependencies; TDD with unit tests, one testscript per command and an end-to-end fold. Left as a follow-up at sign-off: a write lock so concurrent writes cannot duplicate handles or actors.

## title

Shipit plan, milestone 1 task 6: the CLI frame and writes

## Links

### In

- derived-from [repo-init](../capability/rhiza/model/repo-init.md)
- derived-from [current-change](../capability/rhiza/change-sets/current-change.md)
- derived-from [record-shortcut](../capability/rhiza/change-sets/record-shortcut.md)
- derived-from [model-writes](../capability/rhiza/change-sets/model-writes.md)
- derived-from [apply-files](../capability/rhiza/change-sets/apply-files.md)
- derived-from [actors](../capability/rhiza/lineage/actors.md)
- derived-from [command-schemas](../capability/rhiza/guide/command-schemas.md)
- derived-from [t6-d1-narrow-model-seam](../decision/t6-d1-narrow-model-seam.md)
- derived-from [t6-d2-frame-writes-init-scope](../decision/t6-d2-frame-writes-init-scope.md)
- derived-from [t6-d3-api-registry](../decision/t6-d3-api-registry.md)
- derived-from [t6-d4-reflected-schemas](../decision/t6-d4-reflected-schemas.md)
- derived-from [t6-d5-actor-lookup](../decision/t6-d5-actor-lookup.md)
- derived-from [t6-d6-current-change-file](../decision/t6-d6-current-change-file.md)
- derived-from [t6-d7-tdd-testscript](../decision/t6-d7-tdd-testscript.md)
- derived-from [t6-d8-ref-resolution](../decision/t6-d8-ref-resolution.md)
- derived-from [t6-d9-revise-existing-op](../decision/t6-d9-revise-existing-op.md)
- derived-from [t6-d10-apply-file-format](../decision/t6-d10-apply-file-format.md)
- derived-from [t6-d11-interim-model-unexported](../decision/t6-d11-interim-model-unexported.md)
- derived-from [t6-d12-error-shape](../decision/t6-d12-error-shape.md)
- derived-from [t6-d13-lifecycle-verbs](../decision/t6-d13-lifecycle-verbs.md)
- derived-from [t6-d14-guardrails](../decision/t6-d14-guardrails.md)
- derived-from [t6-d15-done](../decision/t6-d15-done.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
