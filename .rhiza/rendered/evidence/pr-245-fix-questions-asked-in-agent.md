# pr-245-fix-questions-asked-in-agent

evidence · `01M6K41N24YDP3QGHEZJE0HJBT` · actor [dphaener](../actor/dphaener.md)

## head

588dddd40f029dabd2137614dcef2b36b2a4d288

## pr

https://github.com/dphaener/rhiza/pull/245

## result

1 check passed: check SUCCESS

## statement

Pull request #245 (https://github.com/dphaener/rhiza/pull/245) passed CI at 588dddd

## status

verified

## Links

### Out

- verifies [agent-tab-blocked-notice](../capability/rhiza/orchestration/notifications/agent-tab-blocked-notice.md)
- verifies [rex-banner-upstream-draft](../capability/rhiza/orchestration/notifications/rex-banner-upstream-draft.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
