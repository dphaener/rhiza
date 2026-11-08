# milestone-interview-reasoning-follows-supersession

decision · `01M6XWC3THKT0H06REGBKXX7JT`

## reasoning

The finding offers two fixes: revert the statement to its original text, or update the reasoning. Reverting brings back round 1's f1 (the landed decision's statement would again say the interview starts without a free slot, contradicting plan-runs-share-max-runs and the change set's intent that nothing may exceed the cap), and the critic said so. The statement already names plan-runs-share-max-runs as the superseding decision, so the misattribution concern is mostly met; what remains is the stale reasoning sentence, which now argues against its own statement. plan-runs-share-max-runs's reasoning already gives the replacement argument (a plan run starts agents, was in the reported burst, and the owner's 'must not exceed the cap' outranks an assumed exception), so pointing the reasoning at it makes the two decisions agree without changing any behaviour. Wording is unchanged because neither decision quotes the owner.

## statement

Keep milestone-interview-starts-at-once's rewritten statement (the interview waits while running is at rhiza.daemon.maxRuns, with its parenthetical naming plan-runs-share-max-runs) and update that decision's reasoning to match: drop the claim that the interview 'costs little compute and need not compete for a slot', and say instead that plan-runs-share-max-runs superseded the start without a free slot because a plan run still starts agents and the owner's words 'must not exceed the cap' outrank the earlier assumption, while first claim on the next free slot keeps the point that an owner interview is not queued behind build picks. The rest of the reasoning (one owner, waiting status at launch, Rex notifies from OSC 7501, usage reserve does not hold it, pause and disk-low still hold it) stays. Its wording stays '(assumed, would have escalated)': both the original and the supersession are arbiter assumptions, so no owner ruling is misattributed.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6xv4tr07qx5x8xpxh8z9m6z-r3-f2](../question/review-01m6xv4tr07qx5x8xpxh8z9m6z-r3-f2.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
