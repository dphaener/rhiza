# severity-contract

capability · `01M76XZGFGB9C8J6GA5NF3RBQE` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide -run 'Severity|Critic|Plan' asserts the expanded plan-review topic and the review critic topic contain the severity block byte for byte, that both critic texts name confirmed, downgraded and dropped, and that the block names no Go, make or Rhiza-repo command.

## statement

The review role topics and the plan-review topic hold one shared severity block (internal/guide/text/shared/severity.md, included by marker): blocker (behaviour broken, a build or test fails, or a leaf is missing), major (an acceptance point unmet or a real defect), minor (style, naming, docs or a non-blocking follow-up) and nit, with a mapping per check; and the critic topics say confirm, downgrade or drop, never add or raise.

## Links

### Out

- decided-by [efficient-cycle-milestone-shape](../../../../../../decision/efficient-cycle-milestone-shape.md)
- decided-by [review-role-topics-shape](../../../../../../decision/review-role-topics-shape.md)

### In

- verifies [severity-contract-verified](../../../../../../evidence/severity-contract-verified.md)
- verifies [severity-contract-qa-1](../../../../../../evidence/severity-contract-qa-1.md)
- verifies [severity-contract-verified-2](../../../../../../evidence/severity-contract-verified-2.md)
- verifies [severity-contract-verified-5](../../../../../../evidence/severity-contract-verified-5.md)
- verifies [pr-510-critic-downgrade](../../../../../../evidence/pr-510-critic-downgrade.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
