# review-01m6jb1w-r3-f1

question · `01M6JCC9CSH2PS08H1JPAZHZKD`

## statement

Plan review round 3 left finding f1 (contradiction on headless-blocked-notice-in-agent-tab): The decision is derived-from Source mac-banner-plan, but the Source's summary still records the original ruling that review rounds 1 and 2 changed. The summary says the wrapper writes 'a working report when the block ends', that the Mac is checked and the report is drafted only 'if it still only chimes', and that the draft is 'pointed at in the inbox'. The decision and its leaves now say: no working report, skip arbitrated questions, always draft from VM captures, and point the owner at it with an open record question. The provenance therefore contradicts the decision it is supposed to back. Fix: note in the Source (or in the decision's reasoning) that the plan review revised the ruling on these points and why. Do not leave the Source claiming the arbiter decided what the plan no longer does. The critic confirmed it: Source mac-banner-plan, the decision's only derived-from Source, has a summary that says the arbiter decided three things: the wrapper writes 'a working report when the block ends', a draft is written only 'if it still only chimes', and the draft is 'pointed at in the inbox'. The decision's statement now says the opposite on all three: no working report, the draft is always written from VM captures with the Mac fields left blank, and an open record question points to it. It also adds a skip for arbitrated questions that the summary never mentions. The decision's reasoning does not say the plan review changed the ruling. So the provenance contradicts the decision it backs. The Source or the reasoning must record that review rounds 1 and 2 revised these points, and why. How should it be settled?

## Links

### In

- answers [plan-review-revisions-in-decision-reasoning](../decision/plan-review-revisions-in-decision-reasoning.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
