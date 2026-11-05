# agent-tab-blocked-notice

capability · `01M6JBG2WB3SDJJE4R0YFTWTV1` · parent [notifications](../notifications.md)

## acceptance

A testscript drives a headless run (driver with no REX_BLOCK) whose agent session, under a fake rex whose events stream gives the agent's blocked permission report, then the echo of the wrapper's own blocked report, then the agent's working report: the wrapper's terminal output holds exactly one root-record blocked report with app=rhiza, kind=permission, title the tab and msg naming the change set, step and block, and no working report; rhiza session show --json for the session gives, after the block, Rex status working (the agent's report) and, while blocked, the agent's own blocked report and message, never the wrapper's echoed message; the same for a question block in a non-autonomous run (kind=question); in an autonomous run a question block (answered by the arbiter) makes the wrapper write no blocked report while a permission block still does; the driver writes no report; a --foreground run with REX_BLOCK set still writes its notice to its own block as run-blocked.txtar checks

## statement

When a run's agent session blocks on something the owner must answer, its own tab's wrapper (rhiza session run) writes rhiza's root-record OSC 7501 blocked report to its terminal (rexNotice's bytes: the kind, the tab as title and a message naming the change set, the step and the question), so a headless run's block reaches Rex with no driver, daemon or watcher tab. Like the driver's notice (internal/run/engine.go notifies only when !r.arbitrates(b)), it writes nothing for a block the run's arbiter answers: in an autonomous run (the session record's Autonomous, already read in internal/api/session_run.go) a question block writes no notice, while a permission block still does. It writes no working report when the block ends: the agent's own next report on the same root record (working, idle or done) replaces the blocked notice, so the wrapper never overwrites the agent's later state. followRex keeps the wrapper's own report out of the session record: rexStatus has no app filter, so the wrapper remembers the report it wrote and skips the matching echo (same state, kind and message) from rex events, and the record's Rex status stays the agent's own report

## Links

### Out

- decided-by [headless-blocked-notice-in-agent-tab](../../../../decision/headless-blocked-notice-in-agent-tab.md)
- decided-by [rex-notice-by-program-status](../../../../decision/rex-notice-by-program-status.md)

### In

- depends-on [rex-banner-upstream-draft](rex-banner-upstream-draft.md)
- verifies [agent-tab-blocked-notice-verified](../../../../evidence/agent-tab-blocked-notice-verified.md)
- verifies [agent-tab-blocked-notice-qa-1](../../../../evidence/agent-tab-blocked-notice-qa-1.md)
- verifies [agent-tab-blocked-notice-verified-2](../../../../evidence/agent-tab-blocked-notice-verified-2.md)
- verifies [pr-245-fix-questions-asked-in-agent](../../../../evidence/pr-245-fix-questions-asked-in-agent.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
