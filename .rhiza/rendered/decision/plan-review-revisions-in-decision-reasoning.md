# plan-review-revisions-in-decision-reasoning

decision · `01M6JCEGEJFXHC7H5T6K1DS71Y`

## reasoning

The finding offers two fixes: note the revision in the Source, or in the decision's reasoning. A Source is a scrubbed record of one conversation (rhiza-scrub keeps what was said and decided at the time). Its summary is accurate for the plan session it records. Rewriting it to cover later review rounds would make it record things that session never said. The decision is the node that review rounds 1 and 2 changed, so its reasoning is where the change and its reasons belong, and the derived-from link stays true as provenance for the first ruling. I compared the Source's summary with the decision's statement. The summary says: a working report when the block ends; draft only if it still only chimes; pointed at in the inbox. It does not mention skipping arbitrated questions. The statement says the opposite on all three and adds the skip. The four reasons I give come from the statement and its existing reasoning: the agent's own next report replaces the notice; the owner need not answer arbitrated questions; the builder runs on the VM, while the Mac-side fields are left for the owner; and an open question names the draft, as the record's way of reaching the owner. This is a model-only fix inside the change set's intent, and it changes no leaf or acceptance check.

## statement

Settle the contradiction in the decision, not the Source: leave Source mac-banner-plan unchanged, since it truthfully records what the autonomous plan session's arbiter first decided, and add a closing paragraph to the reasoning of headless-blocked-notice-in-agent-tab. The paragraph says plan review rounds 1 and 2 revised that first ruling, which the Source records, on four points, and gives the reason for each. (1) No working report when the block ends: the agent's own next status report already replaces the notice, and a second report would fight Claude Code's state in Rex and in the session record. (2) Blocks the run's arbiter answers are skipped, as the driver's notice already does: the owner need not answer them, so a banner for them is noise. (3) The upstream draft is always written from VM captures with the Mac fields left blank, not only 'if it still only chimes': the builder runs on the VM and cannot see the Mac, so it cannot check for the banner, and the owner fills those fields and files the draft only if the banner is still missing. (4) The draft is pointed at by an open question in the record, not 'in the inbox': open questions are how the record reaches the owner, and a checker can verify that the question exists. The derived-from link to the Source stays.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6jb1w-r3-f1](../question/review-01m6jb1w-r3-f1.md)

## History

- `01M6JB1XQVPJ3J50YKX8MSBVRX` Fix Questions asked in agent tabs on the VM chime on the owner's Mac but show no notification banner. Owner, 2026-11-03, after a test question in a VM agent tab: "I have seen no notification on my mac. But I have heard chimes". Find why Rex on the Mac plays the sound for a VM block's OSC 7501 blocked status but posts no macOS notification (Rex's notification rules for remote hosts or focused windows, the status payload Claude Code or rhiza's blocked report writes, or macOS permissions for Rex) and fix what is ours; file a report with Rex upstream for what is not. This blocks the VM move's leaf rhiza-on-vm (change set 01M6FENG43982SZ, run 01M6H8HYW4ZT failed at build, waiting in the inbox), whose acceptance needs a Rex notification on the Mac; the owner chose: "Keep it open, bug first"
