# concepts-why-voice

capability · `01M71G7Q3G8SDGPGR4W6ZEQBSV` · parent [documentation](../documentation.md)

## acceptance

go test -trimpath ./internal/doclinks/... passes (no broken relative link in the real README.md and docs/) and make self-check passes (rhiza check's docs-commands rule finds no rhiza command in README.md or docs/ that does not exist), grep finds neither {"writes" nor "quotes" in docs/concepts/why.md, and the owner approves the page at the run's sign-off

## statement

docs/concepts/why.md reads in the owner's voice and writes the decision and its link with plain rhiza add and rhiza link commands instead of an apply JSON file; the Source example shows no hand-written JSON either: it has the agent scrub the conversation with the rhiza scrub skill, which writes the file and runs rhiza source add, and the page names rhiza source add and its --raw, --raw-sha256 and --raw-session flags in prose

## Links

### Out

- decided-by [m31-polish-request](../../../decision/m31-polish-request.md)

### In

- verifies [concepts-why-voice-verified](../../../evidence/concepts-why-voice-verified.md)
- verifies [concepts-why-voice-qa-1](../../../evidence/concepts-why-voice-qa-1.md)
- verifies [pr-426-concepts-polish](../../../evidence/pr-426-concepts-polish.md)

## History

- `01M71G4G92NEJ33C04DSKC4RMK` Concepts pages in the owner's voice: docs/concepts/model.md, change-sets.md and why.md rewritten with humanize, every write shown as plain commands (no hand-written apply or Source JSON), signed off by the owner
