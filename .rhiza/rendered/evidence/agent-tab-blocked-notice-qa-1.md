# agent-tab-blocked-notice-qa-1

evidence · `01M6JV7HGQ78TF5KY6612HW6Z8` · actor [claude-code-eaefca50](../actor/claude-code-eaefca50.md)

## result

PASS

## statement

QA at e2cbef27 against real Rex on the VM (rex 0.0.1793645657-gd018d2), with the branch's built rhiza. In a scratch repo and a scratch Rex session, rhiza session start --run run-1 made the session record (tab opening stubbed); then rex window new ran the built `rhiza session run <sid>` in a real Rex tab with a fake claude on PATH that writes OSC 7501 working, then blocked (kind permission or question), then working on a resume file. rex events --json for the tab's block showed, for a permission block: claude-code working, claude-code blocked permission 'approve Bash: touch x', then app=rhiza blocked kind=permission with title 'repo: Agent runs · run run-1 · build' and msg 'Agent runs: build waits on you in tab repo: Agent runs · run run-1 · build: approve Bash: touch x', then claude-code working; no rhiza working report. rhiza session show gave 'rex: blocked permission: approve Bash: touch x' while blocked (the agent's message, not the echoed one) and 'rex: working: Back to work' after. Non-autonomous question block: rhiza blocked kind=question record written. Autonomous record (autonomous=true) question block: no app=rhiza record at all. Autonomous permission block: app=rhiza blocked kind=permission written. Wrapper exit 0 every time. go test ./cmd/rhiza -run 'TestScripts/(session-run-blocked-notice|run-blocked)$' -count=1: ok. Not exercised by hand: a full headless run driver, or the Mac-side banner.

## Links

### Out

- verifies [agent-tab-blocked-notice](../capability/rhiza/orchestration/notifications/agent-tab-blocked-notice.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
