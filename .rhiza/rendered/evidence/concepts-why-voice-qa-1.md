# concepts-why-voice-qa-1

evidence · `01M73EE88NJJA4YB0P291K2ZNV` · actor [claude-code-f2a92615](../actor/claude-code-f2a92615.md)

## result

PASS

## statement

QA at 5074bbfb with a binary built from the branch, in a fresh repo: ran docs/concepts/why.md's decision block as written (rhiza add decision with a multi-line --field statement and a single-quoted --field quote containing double quotes, then rhiza link export-csv decided-by); rhiza show printed the quote verbatim with its quotes and the decides link. rhiza source add <handle> <file> --raw <file> added a source node as the page describes, and the rhiza-scrub skill the page points to exists. go test -trimpath ./internal/doclinks/... passes; rhiza check on a clean clone at the branch head prints ok. Owner approval is left to the run's sign-off.

## Links

### Out

- verifies [concepts-why-voice](../capability/rhiza/documentation/concepts-why-voice.md)

## History

- `01M71G4G92NEJ33C04DSKC4RMK` Concepts pages in the owner's voice: docs/concepts/model.md, change-sets.md and why.md rewritten with humanize, every write shown as plain commands (no hand-written apply or Source JSON), signed off by the owner
