# concepts-change-sets-voice-verified

evidence · `01M73DJT7KNB4D1A8JCA265TFH` · actor [claude-code-4b98a4f3](../actor/claude-code-4b98a4f3.md)

## result

PASS

## statement

At 6204d340: go test -trimpath ./internal/doclinks/... passes. make self-check reports no docs-commands finding (no unknown rhiza command in README.md or docs/); its only failures are stale .rhiza/rendered actor pages and index.md from the run's agent actor records, the same systemic render staleness the passed concepts-model-voice leaf had, not from the docs change. Read the diff of docs/concepts/change-sets.md: it now shows model writes as plain rhiza add capability/decision, rhiza link ... decided-by ... and rhiza set ... --field commands, whose syntax matches rhiza add/link/set --help, and the prose keeps the owner's plain voice. Owner approval is left to the run's sign-off.

## Links

### Out

- verifies [concepts-change-sets-voice](../capability/rhiza/documentation/concepts-change-sets-voice.md)

## History

- `01M71G4G92NEJ33C04DSKC4RMK` Concepts pages in the owner's voice: docs/concepts/model.md, change-sets.md and why.md rewritten with humanize, every write shown as plain commands (no hand-written apply or Source JSON), signed off by the owner
