# release-gates-qa-1

evidence · `01M786MRHSHC6HNGB9S0VCYP93` · actor [claude-code-c6f5cfbe](../actor/claude-code-c6f5cfbe.md)

## result

PASS

## statement

QA at 8ff82f55d: built tools/release (go build -trimpath) and ran `release check` by hand, with fake gh scripts on PATH, on this repo and on a scratch clone whose log I extended. CI gate: success passed; failure, in_progress (pending) and [] (no run) each FAIL with a reason; gh exiting 4 aborts with exit 1. Soak: a missing -runs file fails with 'not in the daemon run history'; a 30h/48h-old first run fails because a bug request came after; the earliest of duplicate lines counts; RHIZA_RELEASE_RUNS works as the default; a malformed line errors with file:line. Patch: a v0.2.1 tag 1 minute old fails 'last patch was 1m0s ago'. Severe: a 'severe: data loss on land' note on the bug fix gives soak and patch 'skipped' (ready even with HEAD missing from the history), and a red CI still blocks. Milestone: a blocking question on a serving change set fails the gate (lists which-frob) and passes once an answers link lands. Minor: -commit is not resolved to a SHA, and the in-flight check can never fire.

## Links

### Out

- verifies [release-gates](../capability/rhiza/distribution/release-gates.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
