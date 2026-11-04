# agent-tab-blocked-notice-verified-2

evidence · `01M6JXES9C7ZCCEPT94867PH9P` · actor [claude-code-3458be59](../actor/claude-code-3458be59.md)

## result

PASS

## statement

Verifier claude-code-3458be59 at bd4bc087 (attempt 2, after review cycle 1). Read the diff: blockedNotice.notify now remembers state=blocked, the kind only when rexRecord writes it (permission/question/auth), and rexField(rexNoticeMsg(n), rexMsgMax), the same msg rexRecord encodes, so a multi-line or over-long question's echo matches and is skipped (quality f1). Ran go test ./cmd/rhiza -count=1 -run TestScripts/(headless-run-blocked-notice|session-run-blocked-notice|run-blocked|headless-run-agent-tabs)$: all PASS. headless-run-blocked-notice drives rhiza run start --foreground with REX_BLOCK empty under a fake rex and claude: each agent pane holds exactly one state=blocked:kind=permission:app=rhiza:title= report and the driver's output holds no 7501 or state= report (compliance f1). session-run-blocked-notice covers, with the agent's blocked report, the wrapper's echo and the agent's working report: one root-record blocked report with app=rhiza and kind/title/msg naming change set, step and tab, no working report; the record's Rex status is the agent's blocked report and message while blocked and working: Back to work after; the same for a question in a non-autonomous run and a new multi-line question case; an autonomous question writes no report while a permission block does. run-blocked still passes for --foreground with REX_BLOCK. Full go test ./cmd/rhiza -count=1: ok. go vet ./...: ok; gofmt -l: clean. go test ./internal/api fails only TestReadBudget (timing medians ~50-60ms under load average ~9); the diff does not touch the read path. Gaps: the session record is checked through session show's '- rex:' text lines, not --json (compliance f2, minor); the headless script's events stream gives no echo, the echo being exercised in session-run-blocked-notice instead.

## Links

### Out

- verifies [agent-tab-blocked-notice](../capability/rhiza/orchestration/notifications/agent-tab-blocked-notice.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
