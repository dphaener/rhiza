# concepts-change-sets-voice

capability · `01M71G7Q2AXY67KQS8MP979WZW` · parent [documentation](../documentation.md)

## acceptance

go test -trimpath ./internal/doclinks/... passes (no broken relative link in the real README.md and docs/) and make self-check passes (rhiza check's docs-commands rule finds no rhiza command in README.md or docs/ that does not exist), and the owner approves docs/concepts/change-sets.md at the run's sign-off

## statement

docs/concepts/change-sets.md reads in the owner's voice and shows any model write as plain rhiza add, link or set commands

## Links

### Out

- decided-by [m31-polish-request](../../../decision/m31-polish-request.md)

### In

- verifies [concepts-change-sets-voice-verified](../../../evidence/concepts-change-sets-voice-verified.md)
- verifies [concepts-change-sets-voice-qa-1](../../../evidence/concepts-change-sets-voice-qa-1.md)
- verifies [pr-426-concepts-polish](../../../evidence/pr-426-concepts-polish.md)

## History

- `01M71G4G92NEJ33C04DSKC4RMK` Concepts pages in the owner's voice: docs/concepts/model.md, change-sets.md and why.md rewritten with humanize, every write shown as plain commands (no hand-written apply or Source JSON), signed off by the owner
