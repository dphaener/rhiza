# conduct-skill

capability · `01M79KPD8SQX6BBRG3N6WRDD29` · parent [agent-guide](../agent-guide.md)

## acceptance

go test -trimpath ./internal/guide -run TestConduct shows the conduct topic has the role sections in order, names each phase with its command (rhiza apply for missing leaves, rhiza run spawn with the builder, verifier, check-plan, check-code, check-qa, critic and fixer roles, rhiza run findings, rhiza run review-check, rhiza land, rhiza render, rhiza session done), carries the shared never-list, command-budget and session-done blocks and the never-do list (merge, push, edit .rhiza/log, pass its own gate), names rhiza guide contract, and stays under the guide size cap; go test -trimpath ./internal/api -run TestGuideInstall shows rhiza guide install --agent claude writing skills/rhiza-conduct/SKILL.md with the conduct topic. Touches internal/guide/text/conduct.md, internal/guide/guide.go, internal/guide/conduct_test.go, internal/api/guide_install.go, internal/api/guide_install_test.go.

## statement

Rhiza ships a conduct guide topic and rhiza-conduct skill for the conductor of a change set: read the slice and the code; plan the leaves with rhiza apply when there are none, under the autonomous rules (questions go to the arbiter); spawn a builder and then a verifier per leaf with a tight brief; spawn the plan, code and qa checks (loading check-plan, check-code and check-qa); spawn a critic only for a check whose latest round has a blocker or major finding, and a fixer for confirmed findings; check rhiza run review-check passes; land and render on the branch; then rhiza session done. It never merges, pushes, edits .rhiza/log or passes its own gate.

## Links

### Out

- depends-on [run-spawn](../../workflows/workflow-runs/run-spawn.md)
- depends-on [guide-contract](guide-contract.md)
- depends-on [check-plan-code-topics](check-plan-code-topics.md)
- decided-by [conductor-spawned-sessions](../../../../decision/conductor-spawned-sessions.md)
- decided-by [conductor-plans-leaves](../../../../decision/conductor-plans-leaves.md)
- decided-by [conductor-checks-gate-pick](../../../../decision/conductor-checks-gate-pick.md)

### In

- verifies [conduct-skill-verified](../../../../evidence/conduct-skill-verified.md)
- verifies [conduct-skill-qa-1](../../../../evidence/conduct-skill-qa-1.md)
- verifies [conduct-skill-verified-2](../../../../evidence/conduct-skill-verified-2.md)
- verifies [conduct-skill-qa-2](../../../../evidence/conduct-skill-qa-2.md)
- verifies [pr-575-a-rhiza-conduct-skill-conducts](../../../../evidence/pr-575-a-rhiza-conduct-skill-conducts.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
