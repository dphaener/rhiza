# done-captures-answers

capability · `01M6506ZFX7SJEH3W2TYSQWWJR` · parent [interview-capture](../interview-capture.md)

## acceptance

A testscript whose fake agent leaves a transcript with one picked and one typed answer and says done: the change set holds both answer Decisions when session done returns, a second session done proposes nothing, a later session run exit adds nothing already proposed, and a session with no transcript says done with a warning and keeps its mark. A Go test with a capture hook shows the done mark is not yet set when the capture runs

## statement

rhiza session done proposes the session's answered questions as Decisions before it writes the done mark, so a step that waits on done finds them; a dry run captures nothing, and a capture that fails is a warning while the mark stands

## Links

### Out

- relates-to [interview-on-session-exit](../../../../decision/interview-on-session-exit.md) — adds a capture at done; the capture on exit stays
- decided-by [capture-before-done-mark](../../../../decision/capture-before-done-mark.md)
- decided-by [done-capture-failure-warns](../../../../decision/done-capture-failure-warns.md)

### In

- verifies [done-captures-answers-verified](../../../../evidence/done-captures-answers-verified.md)
- verifies [done-captures-answers-qa-1](../../../../evidence/done-captures-answers-qa-1.md)
- verifies [done-captures-answers-verified-2](../../../../evidence/done-captures-answers-verified-2.md)
- verifies [pr-101-plan-captures-on-done](../../../../evidence/pr-101-plan-captures-on-done.md)

## History

- `01M61GW89J8Y62X62RNJ8Z8HAD` A plan run captures the interview's picks and answers as Decisions when the interview step says done, not only when its session exits.
