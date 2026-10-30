# plan-captures-on-done-interview

source · `01M6505J7GJEFAA0DMFD4994D2`

## quotes

```json
[]
```

## raw_sha256

96655f23f3d9e7e07263f4fadbc9403875f08d299e4921db3a0bfaf6f890a295

## session

966f91a7-f1fb-4e11-a81b-01a18a012bca

## summary

Autonomous plan run, no owner to ask. The change set's intent: a plan run captures the interview's picks and answers as Decisions when the interview step says done, not only when its session exits. Code read: session done (internal/api/session_done.go) writes only the done mark; answers are captured only on exit (captureAfterRun in internal/api/session_run.go) and by session interview; captureInterview skips handles already proposed, so a second capture adds nothing. Token use already reads the transcript at done (recordTokenUse), the precedent for capturing there. Choices settled here, each marked as assumed in its Decision's wording: capture at done before the done mark is written, so a step waiting on done sees the Decisions; a capture failure is a warning and the mark stands; session done reports the Decisions it proposed.

## title

Plan captures on done: autonomous planning of the change set

## Links

### In

- derived-from [capture-before-done-mark](../decision/capture-before-done-mark.md)
- derived-from [done-capture-failure-warns](../decision/done-capture-failure-warns.md)
- derived-from [done-reports-decisions](../decision/done-reports-decisions.md)

## History

- `01M61GW89J8Y62X62RNJ8Z8HAD` A plan run captures the interview's picks and answers as Decisions when the interview step says done, not only when its session exits.
