# review-01m6jzdh-r3-f1

question · `01M6K16EYDVGEGPKMHZ6VWSPPB`

## statement

Plan review round 3 left finding f1 (contradiction on daemon-resumes-landed-plan-run): Its acceptance says the existing daemon-resumes-landed-past-land tests keep passing, but that leaf's test (TestResumableOnlyAcceptedOrLandedPastLand in internal/daemon/pick_test.go on branch cs/fix-a-run-interrupted-after, change set 01M6J965) adds a 'landed-plan' run and asserts Resumable leaves it out, and Resumable there skips landed runs with st.Workflow == plan on purpose. This leaf has to change that assertion, so the acceptance contradicts the code it relies on. It should say the landed-plan case in that test flips to resumed. The critic confirmed it: On branch cs/fix-a-run-interrupted-after (01M6J965, not on main), TestResumableOnlyAcceptedOrLandedPastLand in internal/daemon/pick_test.go adds add("landed-plan", "landed", "plan", true) and wants exactly [run-accepted run-landed-past]. Resumable in pick.go skips landed runs when st.Workflow == plan. So this leaf has to change that assertion, and an acceptance that says the daemon-resumes-landed-past-land tests 'keep passing' contradicts the code. It should say the landed-plan case flips to resumed. How should it be settled?

## Links

### Out

- blocks [daemon-resumes-landed-plan-run](../capability/rhiza/orchestration/daemon/daemon-resumes-landed-plan-run.md)

### In

- answers [decision-landed-plan-test-flips](../decision/decision-landed-plan-test-flips.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
