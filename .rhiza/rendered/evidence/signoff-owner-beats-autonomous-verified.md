# signoff-owner-beats-autonomous-verified

evidence · `01M73HJ2EZSSYF7BN65GKB0T62` · actor [claude-code-a169ab94](../actor/claude-code-a169ab94.md)

## result

PASS

## statement

At aee97577 I read the diff from 3502ad9d: signoff.go adds ownerDeclared, checked after an existing sign-off and before the Autonomous branch, which reads the recorded stage autonomy (stage signoff, source workflow, mode owner) and falls back to the started definition's [autonomy] when the journal has no record; a setting-sourced owner is not a declaration. go test -trimpath -count=1 ./internal/run/ -run TestSignoff passed all 9, including TestSignoffWorkflowOwnerBeatsAutonomous (waits with SignoffWait, no signoff entry, then approve succeeds), TestSignoffAutoAndSettingSignOffAutonomously (workflow auto, setting auto, setting owner and a ship-only declaration each sign off under rule autonomous), TestSignoffOwnerFallsBackToDefinition (no record) and TestSignoffWorkflowOwnerChangesGoBack (request-changes goes back to review, then waits again on visit 2). go test -trimpath ./internal/run/... passed. go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/(run-signoff-workflow-owner|review-signoff-autonomous|workflow-stage-autonomy|run-autonomy-source) passed all four. The new script runs run start --autonomous on a workflow with [autonomy] signoff = owner and fake rex/claude, sees the log line naming the workflow's declaration, state stopped at signoff with signoff_waiting true and no rule autonomous sign-off, then rhiza run approve succeeds with rule owner. The other three scripts are unchanged in the diff

## Links

### Out

- verifies [signoff-owner-beats-autonomous](../capability/rhiza/orchestration/accept-policy/signoff-owner-beats-autonomous.md)

## History

- `01M73FYE2PAG9E034MPBTTJ12D` Fix A workflow that declares 'autonomy signoff owner' is signed off autonomously anyway. The repo's docs workflow (#425) declares signoff = owner (rhiza workflow describe docs: 'autonomy signoff owner'), yet run 01M73CQ4D71CW038ZBNSQB1SJF (Concepts polish, PR #426) logged 'step signoff: signed off under rule autonomous' and went on to land; run 01M73CQ3TCG0TVWWR8ZV62VAE8 (README, PR #427) is headed the same way. Settings: rhiza.autonomy.signoff = auto (default, project); the daemon runs --autonomous. The milestone's Decisions require the owner to sign off each page, so a workflow's declared owner stage must win over the default/--autonomous (or the docs workflow is not honoured). Held by converting PRs #426 and #427 to draft so ship merge cannot merge them; the owner signs off on the PRs by hand.
