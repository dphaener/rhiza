# plan-review-journal

capability · `01M5WSGD28YB7NQC8BQR8MCVFM` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A testscript with RHIZA_RUN set records two findings and a verdict confirming one and dropping the other: the journal holds both entries, and run status --json and its markdown show the round with one confirmed finding. A verdict naming a third finding, a finding of an unknown kind, or one naming an ID that is neither a node nor a link exits 2 with nothing written; with no RHIZA_RUN both commands exit 2. describe-golden is regenerated. Run: go test -trimpath ./cmd/rhiza -run 'TestScripts/(describe-golden|<the plan-review-journal script>)'.

## statement

A plan review round is journaled in its run: rhiza run findings <file> records the reviewer's findings, each naming a node or a link (a model link, or a link between change sets such as depends-on), one kind (missing-acceptance, missing-provenance, unreal-depends-on, outside-intent, contradiction) and its reason; rhiza run verdict <file> records the critic's verdict on each, confirmed or dropped with a reason, and refuses a finding the round does not hold, so a critic can't add one; rhiza run status shows each round's findings and verdicts.

## Links

### Out

- decided-by [plan-workflow-shape](../../../../../decision/plan-workflow-shape.md)
- decided-by [link-findings-allowed](../../../../../decision/link-findings-allowed.md)

### In

- depends-on [plan-review-guide](plan-review-guide.md)
- depends-on [plan-review-rounds](plan-review-rounds.md)
- verifies [plan-review-journal-verified](../../../../../evidence/plan-review-journal-verified.md)
- verifies [pr-40-the-plan-workflow](../../../../../evidence/pr-40-the-plan-workflow.md)
- verifies [plan-review-journal-qa-1](../../../../../evidence/plan-review-journal-qa-1.md)
- verifies [pr-452-fix-a-plan-review-round-whose](../../../../../evidence/pr-452-fix-a-plan-review-round-whose.md)

## History

- `01M5WPTXDS9R3BJVS8PNJQZ834` The plan workflow: Rhiza interviews the owner one topic at a time and writes a milestone or change set straight into the model, with a Source, Decisions quoting the owner's typed words, a capped plan review, and acceptance as the gate.
- `01M76F4G0JQ649RR29M9QPFX68` Fix A plan review round whose findings fail to record still passes the review step, so review-check judges the previous round again and the plan run fails at revise's max_visits with its fixes in. Evidence 2026-11-11, plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels, change set 01M765854D7DGJDBW5CYGJ4SR1): round 2 (22:34Z) had f1 release-provenance and f2 install-script; revise visit 2 fixed both (ops 22:37-22:38Z in 01M76AZ94M8HWY6GCD06H8ES46 and the install-script change set); the round-3 reviewer (session 1f769929-ee04-4c59-961f-01a1ccc326ee) found that round 1's confirmed f6 (an unneeded depends-on link) was never retired, but 'rhiza run findings' rejected its finding because it named the link's ID ('it accepts only node IDs'), and the reviewer had chained 'rhiza session done' with ';', so the session ended and the review step passed ('said done and until held') with no round 3 journaled. check visit 3 then read LatestRound (internal/api/run_review.go:991), which was still round 2, reported its already-fixed findings, and revise went over max_visits 2: 'run 01M76583 failed: step "revise" went over max_visits 2'. Fix: (1) the review step's until (and session done for a reviewer) requires that this visit journaled a new findings round, empty for a clean review, else the step does not pass; (2) review-check refuses to judge when the latest round predates the latest revise; (3) run findings accepts a link (or change-set link) as a finding's target, since unreal-depends-on findings are about links; (4) the review prompt says not to chain session done after a command that can fail. Unstuck by rhiza run resume --retry 01M76583.
