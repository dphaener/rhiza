# milestone-number-reserved-at-run-start

decision · `01M777A77AV16MSTDPMCADCPM6`

## reasoning

The roadmap shows milestone-25 is the live catch-all for daemon and run reliability: its open change set 01M76RH0 (pre-push self-check) was filed from the same CI failure (duplicate handle milestone-32 on PR #456), and the question notes the earlier milestone-request bug served it too; no other live milestone is about planning runs. On the how: the change set's own cause is that the planner counts milestones in its branch's model, which cannot see other in-flight milestone branches. Run journals live in the git common dir (run.Dir(common) in run.go/request.go), shared by every worktree of the repo, so they are the one place every concurrent run on this machine can see; reserving at run start there removes the race at its source with no new storage. Counting every past reservation, not only live runs, keeps a failed or abandoned plan's number from being reused by a later plan that then collides when the first is retried (as happened with milestone-34). Renumbering on conflict at close-out would rewrite handles, fields and done-point names after review, is larger and riskier, and still lets CI see the clash first. On milestone-check: a reservation the planner may ignore is only advisory; milestone-check already runs between the interview and review and already reads the run's milestone state (run_milestone.go), so checking the new milestone intent against the reserved number is a small addition that turns the fix into a guarantee and fails fast and cheaply, before review, PR and CI.

## statement

(1) Intent: this bug change set serves @milestone-25 (the backlog is empty: the daemon builds and ships the queued work on its own), as the earlier milestone-request bug fix and the pre-push check change set 01M76RH0 do. (2) How: the milestone number is reserved when a milestone plan run starts, not chosen by the planner from its own branch. run start for a milestone plan (a milestone request, or --new --milestone) draws N = 1 + the highest of every landed milestone intent's milestone field and every number already reserved in any run journal in the repo's shared run store (git common dir), ended or failed runs included, so a number is never handed out twice and gaps are fine. N is stored in the run's journal, shown by rhiza run status, and given to the planner as RHIZA_MILESTONE_NUMBER; the plan guide (rhiza-plan / plan-milestone) tells the planner to name the milestone intent milestone-N with field milestone=N and done points prefixed mN-. A run that plans an existing milestone (--milestone <ref>) reserves nothing. No renumbering on conflict at close-out. (3) Yes: run.milestone-check also fails a plan whose new milestone intent does not use the run's reserved number (handle milestone-N, field milestone=N, done points prefixed mN-), naming the intent and the expected number in its Fix line, so the clash is caught before review and push rather than in CI.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-d85dd023-1](../question/question-d85dd023-1.md)
- derived-from [plan-milestone-number-clash](../source/plan-milestone-number-clash.md)

### In

- decided-by [milestone-number-reserved](../capability/milestone-number-reserved.md)
- decided-by [planner-gets-milestone-number](../capability/planner-gets-milestone-number.md)
- decided-by [milestone-check-enforces-number](../capability/milestone-check-enforces-number.md)

## History

- `01M7774MAPJ7VCPN1Q9M3RHMV4` Fix Milestone plans running at the same time pick the same milestone number and handle, so the second to merge fails CI self-check with a duplicate handle. Evidence: the agent docs plan (run 01M76EC9, record 01M76ECB) named its intent milestone-32, already held by the cleanup milestone 01M7640B, and failed pr-green-infra on PR #456. I renamed it to milestone-34 by hand (record cs 01M76V0W). The efficient build cycle plan (run 01M76STX, record 01M76SV0) had meanwhile also chosen milestone-34 and failed CI on PR #466 ('duplicate handle: node 01M76NWA... has handle milestone-34, already held by node 01M76XMNA9...'); renamed to milestone-35 by hand (record cs 01M774VK) and resume --retry. Cause: the planner takes the next number from its own branch's model, which can't see other in-flight milestone branches. Fix ideas: the daemon assigns the milestone number when it opens the record change set (counting open milestone records too), or the plan's close-out/rebase renumbers on conflict. The pre-push self-check (01M76R1E) would catch it locally only after merging main.
