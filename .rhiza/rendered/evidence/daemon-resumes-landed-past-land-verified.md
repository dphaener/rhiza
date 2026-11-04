# daemon-resumes-landed-past-land-verified

evidence · `01M6JNG113KKPAYK0FRDW746ND` · actor [claude-code-35d02452](../actor/claude-code-35d02452.md)

## result

PASS

## statement

At 499e0789 I read the diff and ran go test ./internal/daemon -run Resumable -count=1: PASS. TestResumableOnlyAcceptedOrLandedPastLand returns run-accepted and run-landed-past (ship-change, entered land, change set landed), and leaves out a landed run that never reached land, a landed plan run, and runs on rejected, withdrawn, proposed and draft change sets. go test ./cmd/rhiza -run TestScripts/daemon-resume -count=1 passes both daemon-resume and the new daemon-resume-landed.txtar. The builder added the scenario as a sibling file, not inside daemon-resume.txtar. It takes a ship-change run that stopped after land, makes it interrupted, lands its change set and runs daemon --once: Resumed run, exactly one daemon-resumed journal entry, one detached run resume --foreground. A second --once resumes nothing. The resumed run then prints succeeded past its stopped wait step, with no 'no longer accepted'. The same txtar on the parent commit c092ae36 fails at 'Resumed run', so it covers the fix. make check under load average 22-30 failed only arbiter-agent-tabs (a timeout, which passes when run alone), TestReadBudget and slice TestBudget (timing budgets). TestBudget fails the same way on the parent commit, and the diff touches neither package's read path, so these failures are environmental.

## Links

### Out

- verifies [daemon-resumes-landed-past-land](../capability/rhiza/orchestration/daemon/daemon-resumes-landed-past-land.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
