# concepts-model-voice-qa-1

evidence · `01M73EE825YEX6HWDEG9RRY929` · actor [claude-code-f2a92615](../actor/claude-code-f2a92615.md)

## result

PASS

## statement

QA at 5074bbfb with a binary built from the branch, in a fresh git repo after rhiza init and rhiza change new: ran docs/concepts/model.md's commands as written (rhiza add capability export-csv --field statement=..., rhiza set export-csv --field acceptance=..., rhiza add decision csv-streaming --field statement=..., rhiza link export-csv decided-by csv-streaming); each wrote an entry and rhiza show export-csv listed both fields and the decided-by link. go test -trimpath ./internal/doclinks/... passes, and rhiza check on a clean clone at the branch head prints ok with no docs-commands finding. Owner approval is left to the run's sign-off.

## Links

### Out

- verifies [concepts-model-voice](../capability/rhiza/documentation/concepts-model-voice.md)

## History

- `01M71G4G92NEJ33C04DSKC4RMK` Concepts pages in the owner's voice: docs/concepts/model.md, change-sets.md and why.md rewritten with humanize, every write shown as plain commands (no hand-written apply or Source JSON), signed off by the owner
