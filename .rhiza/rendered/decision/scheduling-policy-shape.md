# scheduling-policy-shape

decision · `01M6VYH21ZJVTTS56VT3D87NYH`

## reasoning

Read the question, the change set's intent, decision workflows-are-the-rules in PR #353 (owner: "we should be able to have bug fixes without a milestone ... even no milestones at all", "rhiza ships with sensible defaults but should be flexible"), internal/api/run_milestone.go (the check fails any plan change set serving no milestone), internal/api/daemon.go:1150-1162 (backlog groups always skipped), lookupWorkflow in internal/api/run.go (personal, then repo .rhiza/workflows, then built-in) and the root rhiza.toml (only [projects] and [daemon], machine-level). (1) This run's plan cannot pass milestone-check without a milestone; milestones 26-29 are the live ones and 29 (one root daemon building several projects, the owner's VM priority) is the one whose daemon pick this change rewrites, so it fits best; creating a milestone needs the owner. (2) Scheduling is a pick across change sets, not a step in one run, so it cannot live in a single workflow file; it must be per repo and versioned like workflows, hence a checked-in .rhiza file beside .rhiza/workflows rather than machine config. (3) The intent fixes the default (after milestone work, oldest accepted first); 'skip' keeps today's behaviour available and 'before' covers the other order the intent names, with nothing more until someone asks. Failing loudly on a bad value follows the model's habit of check problems over silent fallbacks. (4) The intent relaxes fix 01M6JDMJ so a plan links a milestone only when one fits or --milestone names one; the check then guards exactly the two cases where a missing milestone is still a mistake, and stays a workflow step the repo can override, as the principle requires.

## statement

(1) Change set 01M6VY7 serves @milestone-29 directly: the root daemon that milestone builds picks work for every project, and this change is what lets a project's milestoneless work (Nebari's bug fixes, small features) be picked; no new milestone is made for it. (2) The scheduling rule is repo policy checked into the repo, in a new `.rhiza/policy.toml` under a `[schedule]` table, so it travels with the repo and every daemon (per-repo or root) reads the same rule; not in the machine's root rhiza.toml, and not in repo.json, which stays identity only. A missing file or key means the shipped default. (3) `[schedule] backlog` takes `"after"` (default: accepted change sets that serve no milestone are built after every milestone's ready work), `"before"` (built ahead of milestone work) or `"skip"` (today's behaviour: skipped as 'serves no milestone'); within the backlog the order is oldest accepted first. An unknown value is a `rhiza check` problem, and the daemon then skips backlog work with the bad value in the skip detail rather than guess. (4) The plan workflow keeps its milestone-check step, but it fails only when it still means something: a change set of the plan serving no milestone passes, unless the run was started with --milestone (then each change set must serve that milestone, directly or through an intent under it) or the repo's backlog is `"skip"` (then the work would never be built, and the step fails with today's message). A repo that wants no check drops the step in its own plan workflow.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-c611d540-1](../question/question-c611d540-1.md)
- supersedes [request-plans-serve-a-milestone](request-plans-serve-a-milestone.md)
- derived-from [scheduling-as-policy-interview](../source/scheduling-as-policy-interview.md)

### In

- decided-by [schedule-policy-file](../capability/rhiza/orchestration/daemon/schedule-policy-file.md)
- decided-by [daemon-picks-backlog](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-picks-backlog.md)
- decided-by [plan-step-milestone-check](../capability/rhiza/orchestration/daemon/plan-requests/plan-step-milestone-check.md)
- decided-by [planner-serves-milestone](../capability/rhiza/orchestration/daemon/plan-requests/planner-serves-milestone.md)
- decided-by [milestone-check-follows-policy](../capability/rhiza/orchestration/daemon/plan-requests/plan-step-milestone-check/milestone-check-follows-policy.md)
- decided-by [plan-guide-milestone-optional](../capability/rhiza/orchestration/daemon/plan-requests/planner-serves-milestone/plan-guide-milestone-optional.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
