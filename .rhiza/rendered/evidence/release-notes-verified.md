# release-notes-verified

evidence · `01M7862MTES4PRC702NMF45Q51` · actor [claude-code-c4be694e](../actor/claude-code-c4be694e.md)

## result

PASS

## statement

Verifier at 7f7fabf17 (attempt 2): go test -trimpath -count=1 ./tools/release/... -run TestNotes PASS; go vet and gofmt clean on tools/release. The test builds a fixture log (bug, milestone-served, breaking-marked, severe-marked and two kind=record landings, plus one before the tag) and compares Notes() with testdata/notes.golden: breaking changes with reasons and detected differences first, the milestone group with its statement, other changes, bug fixes with the severe one flagged; record change sets left out, output deterministic. Attempt 1's finding is fixed: ReadHistory now reads the open body's kind and Notes skips kind=record (still calling out a breaking mark on one). A scratch build of tools/release run as `check` in this checkout drafted v0.2.0 notes since v0.1.0 in 53 lines, grouped by milestone with bug fixes listed, and no 'Record agent actor' or 'apply' lines.

## Links

### Out

- verifies [release-notes](../capability/rhiza/distribution/release-notes.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
