# release-marks

capability · `01M76B7DJ5HD25ETA8QRYQP9GA` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./tools/release/... reads breaking: and severe: notes from a fixture .rhiza log and ignores other notes; make self-check passes with the updated .rhiza/workflows guidance

## statement

In this repo a change set is marked breaking or a bug fix severe by a note on the change set starting breaking: or severe: with the reason; the plan-bug and review guidance in this repo tell planners and reviewers when to write them, and tools/release reads them

## Links

### Out

- decided-by [m33-batched-patches](../../../decision/m33-batched-patches.md)
- decided-by [m33-breaking-is-major](../../../decision/m33-breaking-is-major.md)

### In

- verifies [release-marks-verified](../../../evidence/release-marks-verified.md)
- verifies [release-marks-qa-1](../../../evidence/release-marks-qa-1.md)
- verifies [release-marks-verified-2](../../../evidence/release-marks-verified-2.md)
- verifies [release-marks-qa-2](../../../evidence/release-marks-qa-2.md)
- verifies [pr-492-release-marks](../../../evidence/pr-492-release-marks.md)

## History

- `01M76B077D23CJ28C690Y166VY` Release marks: a change set can be marked breaking at plan or review, a bug request carries a severity (severe or not), and a bug request records the rhiza build it was filed against
