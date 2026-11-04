# headless-blocked-notice-in-agent-tab

decision · `01M6JBBD0R8CNZZ5FQ3KA1DNPZ`

## reasoning

Rex makes notifications only from bytes written to a terminal, for that terminal's main program's root record (memory rex-notifications-from-program-status; decision rex-notice-by-program-status), and no API sets another block's status. Since #210 daemonDetach unsets REX_BLOCK for drivers (internal/api/daemon.go:886), so runRex has no terminal to write to and the rich notice is lost. The owner ruled out watcher or driver tabs ("There's nothing actionable for me, the user, to do or see", 2026-11-03), so the only terminal left that the owner should look at is the agent tab itself, which already runs rhiza session run as a wrapper that follows its own block's rex events (followRex in internal/api/session_run.go) and so already knows when the session blocks and what its question is. Writing there puts the notice exactly where the action is and reuses rexNotice/rexRecord. Whether Rex shows a banner for it is unknown; Claude Code's own report only chimes, and the change set's intent says to fix ours and report what is not. Filing with an outside vendor is outward-facing and touches the network, which an autonomous run must not do, so the run drafts the report and the owner files it.

Plan review rounds 1 and 2 revised the first ruling, which Source mac-banner-plan records, on four points; the derived-from link stays as provenance for that first ruling. (1) No working report when the block ends: the agent's own next status report already replaces the notice, and a second report would fight Claude Code's state in Rex and in the session record. (2) Blocks the run's arbiter answers are skipped, as the driver's notice already does: the owner need not answer them, so a banner for them is noise. (3) The upstream draft is always written from VM captures with the Mac fields left blank, not only 'if it still only chimes': the builder runs on the VM and cannot see the Mac, so it cannot check for the banner, and the owner fills those fields and files the draft only if the banner is still missing. (4) The draft is pointed at by an open question in the record, not 'in the inbox': open questions are how the record reaches the owner, and a checker can verify that the question exists.

## statement

A headless run's blocked notice is written by the agent tab's own wrapper (rhiza session run), to its own terminal: when its session blocks on something the owner must answer, it writes the same root-record OSC 7501 blocked report rexNotice builds (kind, tab as title, a message naming the change set, step and question). Like the driver's notice it skips blocks the run's arbiter answers (a question in an autonomous run), so the owner is not bannered for questions they need not answer; a permission block still notifies. It writes no working report when the block ends, leaving the agent's own next report to replace the notice, and it keeps the echo of its own report out of the session record, so the record and Rex keep the agent's own state. The driver keeps writing its notice to its own block only when it has one (a --foreground run in Rex); no driver, daemon or watcher tab is opened for it. The Rex upstream report is drafted, not sent: after the notice is built, the builder, on the VM, writes the report from what the VM can capture (Rex server version, the exact OSC 7501 bytes from the worktree's built rhiza and from Claude Code, remote VM host, chime but no banner) into the repo with the Mac-side fields (Mac app version, whether a banner now shows, focus and notification-permission state) left for the owner, and raises an open question in the record naming the draft; that leaf is done when both exist. The owner checks the Mac, fills those fields and files it only if the banner is still missing, including whether Rex notifies on a blocked report that follows the agent's own blocked report.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-b5fe1d11-1](../question/question-b5fe1d11-1.md)
- derived-from [mac-banner-plan](../source/mac-banner-plan.md)

### In

- decided-by [agent-tab-blocked-notice](../capability/rhiza/orchestration/notifications/agent-tab-blocked-notice.md)
- decided-by [rex-banner-upstream-draft](../capability/rhiza/orchestration/notifications/rex-banner-upstream-draft.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
