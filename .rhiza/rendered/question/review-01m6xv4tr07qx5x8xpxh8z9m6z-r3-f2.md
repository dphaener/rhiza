# review-01m6xv4tr07qx5x8xpxh8z9m6z-r3-f2

question · `01M6XW7W29EAN8MHAK7CJJSX60`

## statement

Plan review round 3 left finding f2 (contradiction on milestone-interview-starts-at-once): The plan rewrites this landed decision's statement to say a milestone interview waits at maxRuns, but leaves its reasoning ('it costs little compute and need not compete for a slot') and wording unchanged, so the decision now contradicts its own reasoning and attributes the new rule to the earlier ruling. The supersession belongs in plan-runs-share-max-runs (which already says it supersedes that part); this decision's statement should keep its original text, or its reasoning must be updated to match and point at plan-runs-share-max-runs. The critic confirmed it: In the what-if view milestone-interview-starts-at-once's statement now says the interview waits at rhiza.daemon.maxRuns, while its unchanged reasoning still argues it 'costs little compute and need not compete for a slot'. The decision's reasoning contradicts its own statement. Of the two fixes the finding offers, reverting the statement would bring back round 1's f1, so the reasoning should be updated to say plan-runs-share-max-runs superseded that part. How should it be settled?

## Links

### In

- answers [milestone-interview-reasoning-follows-supersession](../decision/milestone-interview-reasoning-follows-supersession.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
