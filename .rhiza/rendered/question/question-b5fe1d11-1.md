# question-b5fe1d11-1

question · `01M6JB9067EQTYHPG8V1N343KW`

## session

b5fe1d11-4e23-4674-ae5d-01a1a4b11ff3

## statement

Since #210 run drivers are headless (no REX_BLOCK), so rhiza's own blocked-run notice (title = tab, msg naming change set, step and question) is never written; the only status Rex sees is Claude Code's native OSC 7501 report in the agent tab, which chimes but shows no banner. Where should a headless run's blocked notice be written? / The missing banner (chime heard, no macOS banner) is Rex's presentation choice, which its docs leave undocumented. How should the Rex upstream report be handled?

## Links

### In

- answers [headless-blocked-notice-in-agent-tab](../decision/headless-blocked-notice-in-agent-tab.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
