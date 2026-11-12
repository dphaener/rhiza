# release-marks-qa-2

evidence · `01M7833NKT0EYJMEQG6HBM203V` · actor [claude-code-bb06b38b](../actor/claude-code-bb06b38b.md)

## result

PASS

## statement

QA cycle 2 at 29f2cf3c0 (fix 0874dbd08), with scratch builds of cmd/rhiza and tools/release in a fresh git repo. Real `rhiza note add <cs>` notes 'breaking: drops --old', '  SEVERE:  loses data ', 'plain remark', 'breaking changes ahead' and 'breaking: to retire' made `tools/release marks` print breaking [drops --old, to retire] and severe [loses data]; the plain and colon-less notes were ignored. `rhiza retire <note> --reason wrong` while unlanded wrote an entry with body.drop of the note and link create ops and empty ops; after it, marks no longer listed 'to retire' (the cycle 1 major is fixed). A second change set's 'severe: crash on start' was listed under its own ID. `rhiza set <note> --field statement=just a remark now` on an unlanded mark removed its breaking mark. `tools/release marks` on this repo prints {}. go test -trimpath -count=1 ./tools/release/... passes; workflow check plan-bug and bugfix are valid (expected bypass warnings only). Still open, minor: a node.update to a landed note is not read, and an empty 'severe:' reason counts.

## Links

### Out

- verifies [release-marks](../capability/rhiza/distribution/release-marks.md)

## History

- `01M76B077D23CJ28C690Y166VY` Release marks: a change set can be marked breaking at plan or review, a bug request carries a severity (severe or not), and a bug request records the rhiza build it was filed against
