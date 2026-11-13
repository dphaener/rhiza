# review-role-skills-installed

capability · `01M79AH3C4ZPMBZ7TGZZCQG3P6` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/api -run GuideInstall passes with the nine new skill files written with their topic's description and expanded body and a re-run reporting them unchanged; go test -trimpath ./internal/guide/... asserts the review index names all nine topics with both their skill and `rhiza guide` routes and routes `scope` to check-compliance, stays under a third of the TopicCap, and that review.md's old phrases checked by review_test.go and budget_test.go are now found in the role topics; go test -trimpath ./cmd/rhiza -run 'TestScripts/(guide|guide-install)' passes with the new topics listed and the new skill paths in the install output.

## statement

`rhiza guide ls` lists the nine review role topics, guide install writes them as the skills rhiza-check-compliance, rhiza-check-quality, rhiza-check-qa, rhiza-check-model, rhiza-critic, rhiza-fixer, rhiza-ci-fix, rhiza-gate-fix and rhiza-conflict-fix beside the others, and the review topic (still installed as rhiza-review) is a short index naming the nine topics and which role loads which, a check named scope loading check-compliance, so workflows that still say /rhiza-review or `rhiza guide review` keep working.

## Links

### Out

- depends-on [check-compliance-topic](check-compliance-topic.md)
- depends-on [check-quality-topic](check-quality-topic.md)
- depends-on [check-qa-topic](check-qa-topic.md)
- depends-on [check-model-topic](check-model-topic.md)
- depends-on [critic-topic](critic-topic.md)
- depends-on [fixer-topic](fixer-topic.md)
- depends-on [ci-fix-topic](ci-fix-topic.md)
- depends-on [gate-fix-topic](gate-fix-topic.md)
- depends-on [conflict-fix-topic](conflict-fix-topic.md)
- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)

### In

- depends-on [review-prompts-one-line](review-prompts-one-line.md)
- depends-on [ship-fix-prompts-one-line](ship-fix-prompts-one-line.md)
- blocks [leaf-failed-01m79dzc-1](../../../../../../question/leaf-failed-01m79dzc-1.md)
- verifies [review-role-skills-installed-verified](../../../../../../evidence/review-role-skills-installed-verified.md)
- verifies [review-role-skills-installed-qa-1](../../../../../../evidence/review-role-skills-installed-qa-1.md)
- verifies [review-role-skills-installed-verified-2](../../../../../../evidence/review-role-skills-installed-verified-2.md)
- verifies [review-role-skills-installed-verified-5](../../../../../../evidence/review-role-skills-installed-verified-5.md)
- verifies [review-role-skills-installed-verified-8](../../../../../../evidence/review-role-skills-installed-verified-8.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
