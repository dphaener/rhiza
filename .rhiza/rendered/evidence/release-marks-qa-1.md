# release-marks-qa-1

evidence · `01M782JBRMHVB70GW3714CCJMW` · actor [claude-code-ba6c57ba](../actor/claude-code-ba6c57ba.md)

## result

FAIL

## statement

QA at cbb3a6a44 with built rhiza and tools/release binaries in a scratch repo: real `rhiza note add <cs>` notes 'breaking: drops --old', '  SEVERE:  loses data ', 'breaking changes ahead' and 'severe: crash on start' on a second change set gave `tools/release marks` = breaking [drops --old], severe [loses data] and severe [crash on start]; the plain-text note was ignored; editing an unlanded note with `rhiza set` updated its mark. Gaps: a note retracted with `rhiza retire` before landing (written as body.drop, not node.retire) still counted as a mark, and 'severe:' with an empty reason counted with reason ''. `tools/release marks` on this repo prints {}; go test ./tools/release passes; workflow check plan-bug and bugfix pass.

## Links

### Out

- verifies [release-marks](../capability/rhiza/distribution/release-marks.md)

## History

- `01M76B077D23CJ28C690Y166VY` Release marks: a change set can be marked breaking at plan or review, a bug request carries a severity (severe or not), and a bug request records the rhiza build it was filed against
