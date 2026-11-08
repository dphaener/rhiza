# plan-ci-flake-reruns

decision · `01M6WJ4SM7KZVYHEW98G3KJG66`

## reasoning

Intent: the evidence is a Nebari plan run (01M6VY5APJP9, PR #78) blocked by a Rails flake, and the closest precedent, the ci-triage RSpec fix (01M6TB6N, decision-ci-triage-runners), served m27-nebari-on-rhiza ('Nebari's daemon builds and ships a Nebari change set on its own with Rhiza's workflows'); a plan run that dies on an unrelated flake is the same blocker. Rerun: internal/api/ship.go shipRerunFailed stops after one rerun per head (st.RerunHead == resp.Commit && len(st.Reruns) > 0), and plan.toml's pr-green and rebased-green just run 'rhiza ship check', so the run fails on the second red; the question notes the first rerun already fired and failed and a third attempt passed, so the plan needs a second rerun per head, not the first. A budgeted flag on ship check keeps it deterministic, bounded and visible in the workflow file, leaves ship-change and ci-rerun-once-per-head unchanged by default, and avoids an agent step for a plan whose diff is only .rhiza/. Naming: ship check already parses failing tests from --log-failed (parseFailedTests, shipReadFailures), so reusing it names the test; the note says only that the diff is .rhiza-only, never that the failure is unrelated, because a repo whose CI reads .rhiza/ (Rhiza's make check runs self-check) can go red on a plan's own content, so the note stays true whichever workflow or repo calls it. Failing: the change set's intent says 'rerun the failed jobs once (budgeted) and only then fail', and merging on red would bypass the repo's required checks and could hide a real break on main; the reruns absorb ordinary flakes and a persistent red is worth the owner's eye, now with a message that names the failing tests and says the PR changes only .rhiza/.

## statement

Change set 01M6WHT serves intent m27-nebari-on-rhiza (01M6J303FW4PYT5E5CE39QNHJ8). (1) Rerun budget: rhiza ship check takes a --reruns <n> flag (default 1, today's ci-rerun-once-per-head behaviour; n >= 1) that reruns a head's failed Actions jobs up to n times per head, each rerun recorded in the ship record (a per-head count beside RerunHead/Reruns) before gh run rerun --failed is called, so a crash costs a missed rerun, never an extra one; plan.toml passes --reruns 2 to every ship check in wait-pr, pr-green, wait-rebased and rebased-green, giving a plan PR three CI attempts per head; ship-change keeps the default, since its red CI goes to ci-triage and ci-fix. No agent step and no change to max_visits. (2) Naming: when a check is still red, ship check keeps today's detail form ('failed: <checks>; failing tests: <tests>', the tests parsed from the failed runs' logs as it does for ci-triage) and the human line's 'PR #<n> at <short head>: ' prefix; when reruns were spent on the head it adds ' after <n> reruns' after the failed checks, and when the PR's diff from its merge base touches only .rhiza/ paths it appends that the PR changes only .rhiza/, a fact and not a claim about the cause (a check that reads .rhiza/, such as Rhiza's own self-check, can fail on the plan's content), for example: 'PR #78 at 18230d2: failed: test (4/5) after 2 reruns; failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]; the PR changes only .rhiza/'. A head with no rerun spent and a diff touching other paths keep today's detail unchanged. (3) Yes: a .rhiza-only PR whose checks stay red after the reruns still fails the step and ends the run failed; Rhiza never merges past a red check. The run's inbox item carries that message, and the owner judges from it and the failing tests whether the plan is at fault; and the owner resumes with gh run rerun --failed and rhiza run resume --retry as today.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-5d744c51-1](../question/question-5d744c51-1.md)
- derived-from [plan-ci-flake-reruns-interview](../source/plan-ci-flake-reruns-interview.md)

### In

- decided-by [ship-check-rerun-budget](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-rerun-budget.md)
- decided-by [ship-check-failure-names-reruns-and-rhiza-only-diff](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-failure-names-reruns-and-rhiza-only-diff.md)
- decided-by [plan-run-reruns-ci-flakes](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-reruns-ci-flakes.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
