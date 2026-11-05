# daemon-resumes-landed-plan-run-verified

evidence · `01M6MJ36YQX57PE01PE3NKTQD1` · actor [claude-code-826b7528](../actor/claude-code-826b7528.md)

## result

PASS

## statement

At 904ebb48 read the diff: daemon.Resumable (internal/daemon/pick.go) now skips a landed change set only when the run's workflow is not plan and its journal never entered land, so a plan run on a landed record resumes, while rejected/withdrawn/superseded/proposed/draft states still fall to default and are skipped. Ran go test -count=1 ./internal/daemon -run Resumable: TestResumableOnlyAcceptedOrLandedPastLand passes and returns exactly [run-accepted run-landed-past run-landed-plan], excluding withdrawn-plan, rejected-plan, superseded-plan and landed-never (a non-plan run on a landed change set that never reached land); TestResumable and TestResumableLeavesAStoreLockedRun pass. Ran the new testscript daemon-resume-landed-plan (fake rex): a landed plan run with a dead process is resumed by rhiza daemon --once, journals exactly one daemon-resumed entry and one detached run resume, and a second --once resumes nothing. All daemon* testscripts, including daemon-resume and daemon-resume-landed (accepted-only, once, landed-past-land), pass.

## Links

### Out

- verifies [daemon-resumes-landed-plan-run](../capability/rhiza/orchestration/daemon/daemon-resumes-landed-plan-run.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
