# session-dir-acceptance-covers-fixers

decision · `01M7AMCVPFJR2PHRF2MFV6CE9K`

## reasoning

The question comes from plan review round 3, finding f1, which the critic confirmed. run-brief-files' statement and the decision brief-files-and-invoke-placement name five session-directory forms. The round-2 revision added the two checks-step fixer and verifier forms to stop directories colliding, yet the acceptance tests only three forms. A leaf fixer and the change set fixer both carry the role fixer in one cycle, so they are the likeliest collision and need a test most. Extending the acceptance is the fix the finding asks for. It needs no change to the statement or the scheme, and it stays inside brief_test.go, the file the leaf already touches, so it is the narrowest answer that fits the change set's intent.

## statement

run-brief-files' TestSessionDir acceptance covers all five session-directory forms its statement names. To the 'all names different' case add, within one cycle of a checks step: a leaf fixer and a leaf verifier of the same leaf at one attempt (<role>-<leaf>-<cycle>-<attempt>), the same leaf fixer at attempt 2, and the change set fixer at attempts 1 and 2 (fixer-<cycle>-<attempt>). The leaf fixer and the change set fixer of one cycle and attempt must get distinct directories, and none of these may equal a check's or critic's directory of that cycle. The leaf's statement and the rest of its acceptance stay as they are.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7ah6rhkmzrahgjkjfm4vtnt-r3-f1](../question/review-01m7ah6rhkmzrahgjkjfm4vtnt-r3-f1.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
