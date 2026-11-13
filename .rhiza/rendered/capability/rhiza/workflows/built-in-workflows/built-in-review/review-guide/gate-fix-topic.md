# gate-fix-topic

capability · `01M79AH32SZK0CK9G0Z05VJX4H` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a gate-fix test that finds the six section headings in order (Task, Expected outcome, Read first, Must do, Must not do, Context), the shared done, test-scope, unstable-CI, command-budget and never-list blocks expanded, rhiza.ship.localGate and rhiza.ship.prePushCheck, the `Full log:` rule and never rerun `rhiza ship gate`, the flake quarantine with a bug request, the Checks: block and the Review: trailer; a negative case asserts it names no Go, make or repo-specific command.

## statement

`rhiza guide gate-fix` gives the ship's gate fixer its prompt on the six sections: read the failed local gate or pre-push check's kept output and its full log in `rhiza run status --json` instead of rerunning `rhiza ship gate`, reproduce the failing test with its own scoped command until it fails again or plainly passes, fix it or quarantine it with a bug request named at the skip, commit once it and the touched packages pass, with the pre-handoff checklist and a `Review: <run>` trailer, then session done as its own command.

## Links

### Out

- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)
- decided-by [afk-prompts-shipped-default](../../../../../../decision/afk-prompts-shipped-default.md)

### In

- depends-on [review-role-skills-installed](review-role-skills-installed.md)
- verifies [gate-fix-topic-verified](../../../../../../evidence/gate-fix-topic-verified.md)
- verifies [gate-fix-topic-qa-1](../../../../../../evidence/gate-fix-topic-qa-1.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
