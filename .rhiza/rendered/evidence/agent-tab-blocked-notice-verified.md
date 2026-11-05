# agent-tab-blocked-notice-verified

evidence · `01M6JP26776AAN80DJ9FHBM6SM` · actor [claude-code-afebdd89](../actor/claude-code-afebdd89.md)

## result

PASS

## statement

Verified at eb5352e6. Read the diff (internal/api/session_run.go blockedNotice.echoes/notify in followRex; internal/api/run.go rexNoticeMsg split out of rexNotice). Ran go test ./cmd/rhiza -run 'TestScripts/(session-run-blocked-notice|run-blocked)$' -count=1: both PASS, and the new script passed 3 more runs in a row. The script runs session run under a fake rex whose events give the agent's blocked report, the echo of the wrapper's report, then working. It checks exactly one root-record state=blocked:kind=permission:app=rhiza report whose title is the tab and whose msg names the change set, the build step, the tab and the question, and no state=working. session show gives the agent's blocked report while blocked and working after. Also checked: kind=question in a non-autonomous run, no report for an autonomous question, a report for an autonomous permission block. Mutation checks in scratch copies: dropping the echo skip fails the script at line 36 (the record takes the echoed message), and dropping the autonomous-question filter fails it at line 64. The driver writing no report rests on code this commit does not change (cmd/rhiza/ctx.go sets ctx.Rex only when REX_BLOCK is set). The script drives session run with --run directly, not a full headless run driver. go vet, fmt-check and self-check OK. go test ./... -timeout 40m: all cmd/rhiza scripts pass. Only the timing tests internal/api TestReadBudget and internal/slice TestBudget fail, at load average about 16-23 from concurrent agents; the diff touches neither read path. A default make check hit the 10m package timeout because plan-workflow takes about 400s here, and the parent commit 46cc1d58 is just as slow (431s).

## Links

### Out

- verifies [agent-tab-blocked-notice](../capability/rhiza/orchestration/notifications/agent-tab-blocked-notice.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
