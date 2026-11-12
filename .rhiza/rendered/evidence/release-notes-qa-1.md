# release-notes-qa-1

evidence · `01M786MRM79833CX2NJKB58J3H` · actor [claude-code-c6f5cfbe](../actor/claude-code-c6f5cfbe.md)

## result

PASS

## statement

QA at 8ff82f55d: ran the built `release check` by hand. On this repo the draft notes list the detected breaking changes, then a section per finished milestone (m32-*, m33-*) with each change set's intent and the milestone's statement, then the bug fixes. On a scratch clone: the bug fix appears under Bug fixes, flagged '(severe: data loss on land)' once the note is added; a landed record change set is left out; a breaking-marked change set appears under Breaking changes with its title and reason ('Add the frobnicator: removes --old'), and under its milestone heading 'Milestone qa-milestone: Frobnication everywhere' with the other serving change set.

## Links

### Out

- verifies [release-notes](../capability/rhiza/distribution/release-notes.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
