# release-marks-verified-2

evidence · `01M782ZCVPQC4C528PJCBBHDXW` · actor [claude-code-62c4a362](../actor/claude-code-62c4a362.md)

## result

PASS

## statement

Verifier, cycle 2 attempt 2 at 0874dbd08. Read the fix diff: ReadMarks now skips ops that an entry's body.drop retracts, matching changeset.revision.apply, where op IDs are unique and a drop names that change set's ops. go test -trimpath -count=1 -v ./tools/release/... passes, including TestReadMarks: breaking:/severe: notes are read, plain/colon-less/retired notes and retired links are ignored, and a dropped mark is ignored. go vet and gofmt are clean. End to end with scratch builds of cmd/rhiza and tools/release in a fresh git repo: change new, note add 'breaking: drops flag' and 'severe: kept' made release marks print both. An unlanded `rhiza retire <note> --reason` wrote body.drop of the note and link create ops. After that, release marks printed only severe: kept. make self-check in a clean detached worktree passes at cbb3a6a44 (render after build, plan-bug/bugfix guidance included). At 658c6af56 and 0874dbd08 it fails only on stale .rhiza/rendered/index.md (plus actor pages in the run worktree). Those come from the review run's chore(rhiza) writes; the fix commit touches only tools/release.

## Links

### Out

- verifies [release-marks](../capability/rhiza/distribution/release-marks.md)

## History

- `01M76B077D23CJ28C690Y166VY` Release marks: a change set can be marked breaking at plan or review, a bug request carries a severity (severe or not), and a bug request records the rhiza build it was filed against
