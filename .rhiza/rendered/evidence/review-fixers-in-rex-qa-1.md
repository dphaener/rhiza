# review-fixers-in-rex-qa-1

evidence · `01M66DH2A15QFYP3RHS77D32D8` · actor [claude-code-31db92cd](../actor/claude-code-31db92cd.md)

## result

PASS

## statement

Built ./cmd/rhiza at 5dd508d and ran rhiza run start --workflow review --permission-mode manual inside a real Rex block (rex run --wait). It ran on a scratch repo with a built two-leaf change set: a fake claude played the agents and a fake rex opened their tabs. Real Rex got the OSC 7501 reports, read back with rex events --json. Quality's cycle 1 had two confirmed majors, f1 on leaf alpha and f2 on the change set. Rex showed, in order: leaf/1 working 'alpha: fixer working (attempt 2)', a clear of leaf/1, leaf/1 'alpha: verifier working (attempt 2)', then fix 'review: fixer working (attempt 2)', and clears of leaf/1 and fix. The root record stayed 'step review · 4 of 5 checks passed' throughout. Quality cycle 2 then passed, and the run stopped at signoff.

## Links

### Out

- verifies [review-fixers-in-rex](../capability/rhiza/workflows/workflow-runs/run-program-status/review-fixers-in-rex.md)

## History

- `01M602TGXM184GVWJG4NE98P6Y` Rex program status shows a review run's live fixer, verifier and change-set fixer sessions, its open questions and its code check, as rhiza run status already does
