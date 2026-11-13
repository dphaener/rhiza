# severity-shared-block

capability · `01M79AH1SVENNRB0MBMN3TRZCZ` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with tests that the expanded check-*, critic, fixer, plan-reviewer and plan-critic topics each contain the severities block's text exactly once, that the block has the four severities, the confirmed/downgraded/dropped critic outcomes are not in it but in the critic topic, that its mapping names compliance, quality, qa and model and no `scope` check, and that it names no Go, make or Rhiza-repo command.

## statement

The severity definitions (blocker, major, minor, nit) and the per-check mapping live once in internal/guide/text/shared/severities.md (the block plan role topics landed under that name), with the mapping naming compliance, quality, qa and model (no scope) and plan review's kinds (uncovered-intent and oversized-leaf major too), and every review role topic that grades or fixes findings includes it with `<!-- shared: severities nested -->`, as the plan reviewer and plan critic topics do.

## Links

### Out

- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)

### In

- depends-on [check-compliance-topic](check-compliance-topic.md)
- depends-on [check-quality-topic](check-quality-topic.md)
- depends-on [check-qa-topic](check-qa-topic.md)
- depends-on [check-model-topic](check-model-topic.md)
- depends-on [critic-topic](critic-topic.md)
- depends-on [fixer-topic](fixer-topic.md)
- verifies [severity-shared-block-verified](../../../../../../evidence/severity-shared-block-verified.md)
- verifies [severity-shared-block-qa-1](../../../../../../evidence/severity-shared-block-qa-1.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
