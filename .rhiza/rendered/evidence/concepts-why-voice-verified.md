# concepts-why-voice-verified

evidence · `01M73DZK66NR3WTEGQN1HM000V` · actor [claude-code-0616b85a](../actor/claude-code-0616b85a.md)

## result

PASS

## statement

At 90b5eedf: go test -trimpath ./internal/doclinks/... ok. make self-check reports no docs-command (or docs) finding; check.go runs docsCommands unconditionally. Its only fails are outside this leaf: a fold error in the run's chore evidence entry 01M73DBN12V9QWEYY25C1JVHGG (from 06285a98), whose after names an entry not committed in the leaf checkout, and stale .rhiza/rendered for new actor records, which the run renders at close-out. grep -nE '\{"writes|"quotes"' docs/concepts/why.md finds nothing. The diff replaces the apply JSON with rhiza add decision --field ... and rhiza link export-csv decided-by csv-streaming. The Source section asks the agent to use the rhiza-scrub skill, which writes the file and runs rhiza source add, and names --raw, --raw-sha256 and --raw-session in prose. The prose is plain and in the owner's voice. The commit adds only the builder's actor record under .rhiza. Owner approval is left to the run's sign-off.

## Links

### Out

- verifies [concepts-why-voice](../capability/rhiza/documentation/concepts-why-voice.md)

## History

- `01M71G4G92NEJ33C04DSKC4RMK` Concepts pages in the owner's voice: docs/concepts/model.md, change-sets.md and why.md rewritten with humanize, every write shown as plain commands (no hand-written apply or Source JSON), signed off by the owner
