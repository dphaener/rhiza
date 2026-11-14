# conductor-checks-gate-pick

decision · `01M7AHQJAMY18Z3M41T5BJN6D8`

## reasoning

Read the question, the change set's slice (leaves review-three-checks, critic-on-major, conduct-plans-leaves, conduct-skill and Decisions conductor-three-checks, conductor-plans-leaves, conductor-spawned-sessions) and the code. (1) The question says classic must stay unchanged and guide.go and guide_install.go map check-compliance/check-quality to their skills, so renaming or merging those topics would change classic; new topics beside them keep classic intact and give the conductor's checks rules that match their names. (2) run_review.go's reviewCheck validates --check against the current workflow step's declared Checks; declaring the conductor step's checks the same way reuses that path, keeps one source of truth, and meets TestRunFindingsCheckNames (plan, code, qa accepted, scope refused naming the three). (3) run.review-check already is the daemon-run gate command that judges the latest round against critic verdicts for plan reviews; extending it keeps the gate on the journal, so the conductor cannot pass its own gate (conduct-skill's never-do list), and matches critic-on-major and TestGateOpenBlocker. (4) daemon.go's daemonCandidates skips !hasLeaves; conduct-plans-leaves names exactly the accepted, milestone-serving case and says umbrellas stay skipped, and plan.go already marks umbrella change sets not ready. Limiting the change to the conductor default keeps classic dispatch unchanged, and keeping backlog leafless sets skipped stays inside what the leaf asks for.

## statement

(1) Role topics: add two new guide topics and skills, check-plan (the compliance and model checks' rules merged: each leaf's acceptance met as written, under-built, scope creep, unaccounted files, and code that contradicts the slice) and check-code (the quality check's rules under the new name); the plan and code checks load them, qa keeps check-qa, and check-compliance, check-model and check-quality stay as they are for the classic ship-change, review and docs-review, with the review index topic naming which runs load which. (2) Check names: the conductor default's conduct step declares its checks (plan, code, qa) in the workflow definition the way a checks step does, and rhiza run findings --check and rhiza run verdict --check accept a run in that step and validate the name against the step's declared checks, refusing any other (scope included) with a message naming them; nothing hard-codes the three names. (3) Gate: the conductor default's gate step runs rhiza run review-check, extended to conductor runs: it reads the run's journal and, for each declared check, counts its latest round's blocker and major findings that no critic verdict downgraded or dropped as open; any open finding fails the step (exit non-zero), a round with only minors passes with no critic. The count comes from the journal the daemon keeps, never from the conductor's word. (4) Daemon pick: an accepted, ready change set that serves a milestone and has no leaves is picked when its dispatch resolves to the conductor default; umbrella change sets stay skipped (they are already not ready), backlog change sets with no leaves and change sets dispatched to a classic workflow keep the no-leaves skip.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-4efb76f8-1](../question/question-4efb76f8-1.md)
- derived-from [conduct-skill-plan-2026-11-13](../source/conduct-skill-plan-2026-11-13.md)

### In

- decided-by [conduct-skill](../capability/rhiza/guide/agent-guide/conduct-skill.md)
- decided-by [check-plan-code-topics](../capability/rhiza/guide/agent-guide/check-plan-code-topics.md)
- decided-by [review-three-checks](../capability/rhiza/workflows/built-in-workflows/review-three-checks.md)
- decided-by [critic-on-major](../capability/rhiza/workflows/built-in-workflows/critic-on-major.md)
- decided-by [conduct-plans-leaves](../capability/rhiza/orchestration/daemon/conduct-plans-leaves.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
