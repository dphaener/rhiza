# review-code-check-in-rex-verified

evidence · `01M66CWEDR2F0MRX4XCBTXENX1` · actor [claude-code-7f63c65c](../actor/claude-code-7f63c65c.md)

## result

PASS

## statement

At 12a2ea4, go test ./internal/api -run TestProgramStatusReviewCodeCheck -v -count=1 passes: it runs the builtin review shape (command step code running rhiza run code-check, on.pass to the checks step) over a fake runner whose command exits 1, with Rex reporting wired as run.start wires it into a buffer, and logs exactly two root reports: state=working msg="step code", then state=error msg="run failed: step \"code\" ended with outcome \"fail\""; the test asserts both. The diff only adds that test. make check (with RHIZA_ROOT/CHANGE/RUN unset) exits 0.

## Links

### Out

- verifies [review-code-check-in-rex](../capability/rhiza/workflows/workflow-runs/run-program-status/review-code-check-in-rex.md)

## History

- `01M602TGXM184GVWJG4NE98P6Y` Rex program status shows a review run's live fixer, verifier and change-set fixer sessions, its open questions and its code check, as rhiza run status already does
