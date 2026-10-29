# milestone-4

intent · `01M61R4Q8SA2SWHFE9D5TS6ZJD`

## milestone

4

## statement

Rhiza orchestrates: milestone 5 is driven by rhiza daemon, with no orchestrator session, the owner only accepting plans, answering the inbox and signing off.

## Links

### Out

- decided-by [answer-8cd68bc7-1](../decision/answer-8cd68bc7-1.md)
- decided-by [answer-8cd68bc7-2](../decision/answer-8cd68bc7-2.md)
- decided-by [answer-8cd68bc7-6](../decision/answer-8cd68bc7-6.md)
- decided-by [answer-8cd68bc7-8](../decision/answer-8cd68bc7-8.md)
- derived-from [milestone-4-conversation](../source/milestone-4-conversation.md)

### In

- serves change set `01M5ZNV3K796MKWEZM3TBQXGZ8` (Builder sessions' actor records are committed with the work they sign: after a build run no builder's actor record is left untracked in the build worktree, so every Actor: trailer names a committed actor.)
- serves change set `01M5ZNV40RB8VYBB0A3CPXCNXZ` (rhiza run findings numbers findings under the run's lock, so concurrent calls in one session never journal two findings with the same ID (leaf findings and plan review findings alike).)
- serves change set `01M5ZNV4E4AJR52EX631180DMS` (A change set has at most one live run: rhiza run start refuses a second run on a change set whose run is still going, naming that run.)
- serves change set `01M5ZNVHN2YNK6DDVYY8GK5CEQ` (TestGuideInstallConcurrent passes reliably: the flake that fails it on main too (seen once in make check during the build workflow's review) is found and fixed at its cause, without skipping or loosening the test.)
- serves change set `01M602TGXM184GVWJG4NE98P6Y` (Rex program status shows a review run's live fixer, verifier and change-set fixer sessions, its open questions and its code check, as rhiza run status already does)
- serves change set `01M602TH4C724WZB9MKBN5Z1PP` (A review check or critic session that writes no model entries still gets its actor's record entries committed, so no uncommitted record change set is left in the run worktree for a later fixer to sweep up)
- serves change set `01M60AWXJ785WZCDW1G1KA63BR` (Record the workflow step's limits as constraints: a workflow runs another at most once, and a workflow step takes no session, condition or max_visits keys)
- serves change set `01M61GW73H05YHE4DK0CXFD1KH` (No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.)
- serves change set `01M61GW7BC189FSAN4VEQ1A241` (A builder's or fixer's commit counts its Leaf: trailer even when Co-Authored-By sits in a separate final paragraph, and the build and review guides show all trailers in one final paragraph.)
- serves change set `01M61GW7K1PDG7HMFXMN3YNA7B` (Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.)
- serves change set `01M61GW7TJ6BH6YX1T7SQDWK1J` (A failed fixer does not use up a review fix cycle: the fixer is retried within its cycle, and the cycle's checks re-run only once its fixers pass.)
- serves change set `01M61GW821FWSX05A9FC231N65` (rhiza session message reaches a busy session without waiting for its prompt, so the owner or arbiter can steer a session mid-task.)
- serves change set `01M61GW89J8Y62X62RNJ8Z8HAD` (A plan run captures the interview's picks and answers as Decisions when the interview step says done, not only when its session exits.)
- serves change set `01M61GW8H09EPJ1YCR7P2X3E8J` (Rex status shows a session's AskUserQuestion as a question for the owner, not as 'Claude needs your permission'.)
- serves change set `01M61GW8RKSBM4BT1B8B3FR5JH` (A run's final printout and run status show each session's token use, never 'no token use recorded' for sessions that used tokens.)
- serves change set `01M61R3R8BS63ZGH0BVS7QQ5SS` (The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.)
- serves change set `01M61R3RG70WZD5EDH293ZGVWN` (Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.)
- serves change set `01M620GA5ZK6Z854D07NH0R8WX` (The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.)
- serves change set `01M625RC3X2VA17M5V2HVS3HBP` (ship-change's last step, ship, runs after close-out has removed the run's worktree, so it fails starting in a folder that is gone: steps after close-out start in the main checkout, and a ship-change run ends succeeded once its ship command passes.)
- serves change set `01M625W0AJZ6AMCCWGFXQW19FN` (A plan run that ends with confirmed review findings or parked questions still open does not report success: it stops for the owner with them listed, so nothing has to be fixed by hand after a succeeded run.)
- serves change set `01M625RCFY4RB9HC0B69Y08VFV` (A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.)
- serves change set `01M62MKZ6SA46DYC1K2R0W254S` (The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.)
- serves change set `01M62MKZFRPZ2C7QTYAF815XAG` (The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.)

## History

- `01M61H5VY8MD7E7SVS0HQQHQJ0` Plan milestone 4 with the owner: its intent, its done points and the change sets that serve it, planned by Rhiza's plan workflow.
