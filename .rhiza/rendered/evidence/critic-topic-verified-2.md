# critic-topic-verified-2

evidence · `01M79V107W7NBD40YDTNEZCDPA` · actor [claude-code-f13f3f8c](../actor/claude-code-f13f3f8c.md)

## result

PASS

## statement

At 27d008d91 (attempt 2, review fix for compliance f1) re-ran the Checks block: go test -trimpath ./internal/workflow/ -run TestBuiltinRolePromptsScopeTests exits 0, and 1 on an archived copy with the critic test-scope marker removed (guide topic critic lacks the scoped-tests rule); go test -trimpath ./internal/guide/... 0; go test -trimpath ./internal/api/... 0; rhiza ship generate 0 and rhiza render 0 with no diff. Acceptance: TestCriticText checks the six sections in order, the session-done, test-scope, command-budget, never-list and severity blocks expanded, the isolation rule with the own-check-last-round exception, confirmed/downgraded (lower severity)/dropped, the never add/raise/rewrite rule and the rhiza run verdict --check heredoc; TestCriticRecordsNoFindingAndNoCommit bans rhiza run findings and commit text. rhiza guide critic now carries the scoped-tests rule. Diff touches only critic.md, its test, the workflow scope test, the actor record and rendered files; no stubs or TODOs.

## Links

### Out

- verifies [critic-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/critic-topic.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
