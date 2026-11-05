# rex-banner-upstream-draft-verified

evidence · `01M6JTVFX02PHEBYD74P6XPR88` · actor [claude-code-ab5c7c59](../actor/claude-code-ab5c7c59.md)

## result

PASS

## statement

Verified at 48276c56. docs/rex-upstream/notification-banner.md is committed. Its Rex server version matches rex version run on the VM (0.0.1793645657-gd018d2, commit d018d2de). Its OSC 7501 bytes match the passing session-run-blocked-notice testscript (go test ./cmd/rhiza -run TestScripts/session-run-blocked-notice -count=1: ok), and end with ST (ESC \), which rexRecord in internal/api/run.go writes (attempt-1 finding f2 fixed). It holds a rex events --json program_status_changed record for Claude Code 2.1.316 blocked on a permission prompt on the VM (finding f1 fixed), plus the VM host setup, the symptom (chime, no banner) and blank OWNER Mac fields: app version, whether a banner shows, focus and notification permission. rhiza show rex-banner-report-ready --json gives an open question written in a record-kind change set (apply --record). Its statement names the doc and asks the owner to check the Mac, fill the fields and file with Rex if the banner is still missing. make self-check: ok. The commit changes no Go files. make check's go test hit only timing and timeout failures (TestReadBudget, slice TestBudget, a plan-workflow timeout) under VM load averages of 13 to 17, none in code this leaf touches.

## Links

### Out

- verifies [rex-banner-upstream-draft](../capability/rhiza/orchestration/notifications/rex-banner-upstream-draft.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
