# review-01m7ah6rhkmzrahgjkjfm4vtnt-r3-f1

question · `01M7AMAMP79YB822F2TPS54Z7B`

## statement

Plan review round 3 left finding f1 (missing-acceptance on run-brief-files): The statement names five kinds of session directory, but TestSessionDir checks only three: an agent step, a leaves step's builder, and a checks step's check and critic. It never checks a checks step's leaf fixer or verifier (<role>-<leaf>-<cycle>-<attempt>) or its change set fixer (fixer-<cycle>-<attempt>). Those are the names most likely to collide, since a leaf fixer and the change set fixer both have the role fixer in one cycle. Add them to the 'all names different' case. The critic confirmed it: Real: run-brief-files' statement and the decision brief-files-and-invoke-placement give five session-directory forms, but TestSessionDir's acceptance names only an agent step on two visits, a builder, a check, two critics and one check across two cycles. It never checks <role>-<leaf>-<cycle>-<attempt> for a checks step's leaf fixer or verifier, or fixer-<cycle>-<attempt> for its change set fixer. Those forms were added in the round-2 revision to stop collisions, and nothing tests them, so the 'all names different' case does not cover the stated scheme. Major, as the missing-acceptance kind is graded. How should it be settled?

## Links

### Out

- blocks [run-brief-files](../capability/rhiza/workflows/workflow-runs/run-brief-files.md)

### In

- answers [session-dir-acceptance-covers-fixers](../decision/session-dir-acceptance-covers-fixers.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
