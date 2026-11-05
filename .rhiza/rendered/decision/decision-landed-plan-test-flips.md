# decision-landed-plan-test-flips

decision · `01M6K19WHWHD4FSQ54D82GD905`

## reasoning

The finding and the critic agree, and the leaf's own statement says it: the daemon resumes a plan run on a landed change set. 01M6J965's Resumable skips that case on purpose and its test asserts the skip, so this leaf cannot leave that assertion passing; saying 'keep passing' contradicts the leaf's purpose. decision-plan-landed-resumable already settled the behaviour (one rule with run resume and the inbox), so the only fix is the acceptance wording: name the one assertion that flips and keep every other case of those tests as a regression guard. Narrowing the flip to that single case keeps the acceptance checkable and protects the non-plan landed-without-land and withdrawn-plan cases.

## statement

daemon-resumes-landed-plan-run's acceptance is corrected, not the code: the landed-plan case in TestResumableOnlyAcceptedOrLandedPastLand (internal/daemon/pick_test.go, from 01M6J965) flips from left out to resumed, so that test then expects run-accepted, run-landed-past and the landed plan run; Resumable drops its skip of landed runs with Workflow == plan; the rest of that test and the daemon-resumes-accepted-only, daemon-resume-once and daemon-resumes-landed-past-land tests keep passing unchanged

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6jzdh-r3-f1](../question/review-01m6jzdh-r3-f1.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
