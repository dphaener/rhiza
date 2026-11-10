# concepts-change-sets-voice-qa-1

evidence · `01M73EE861J1TP3SCHWM8YGAVJ` · actor [claude-code-f2a92615](../actor/claude-code-f2a92615.md)

## result

PASS

## statement

QA at 5074bbfb with a binary built from the branch, in a fresh repo: ran docs/concepts/change-sets.md's new block (add capability, add decision, link decided-by, set acceptance) in the current change set opened by rhiza change new; all four wrote into it. Checked the surrounding claims by hand: --change <id> wrote into the named change set rather than the current one, RHIZA_CHANGE did the same, a plain write went to the current one, and rhiza change use switched it back. go test -trimpath ./internal/doclinks/... passes; rhiza check on a clean clone at the branch head prints ok. Owner approval is left to the run's sign-off.

## Links

### Out

- verifies [concepts-change-sets-voice](../capability/rhiza/documentation/concepts-change-sets-voice.md)

## History

- `01M71G4G92NEJ33C04DSKC4RMK` Concepts pages in the owner's voice: docs/concepts/model.md, change-sets.md and why.md rewritten with humanize, every write shown as plain commands (no hand-written apply or Source JSON), signed off by the owner
