# review-role-skills-installed-verified-2

evidence · `01M79TEM7RXXFFT21M225F67H2` · actor [claude-code-f1af48e5](../actor/claude-code-f1af48e5.md)

## result

PASS

## statement

Verifier attempt 2 at 848c4e08e (review cycle 2 fix commit, .rhiza records only: node.update on the eight review-rule nodes command-budget-guide, done-own-command-guide, guides-scoped-tests, unstable-ci-guidance, delay-repro-guidance, guides-stdin-findings, guides-one-final-paragraph, role-prompts-scoped-tests, plus actor records and rendered files). Read the ops: each statement now puts the rule in the shared blocks expanded into the nine review role topics and says the review index does not hold it; each acceptance names tests that exist (TestReviewRoleText, TestBuildReviewCommandBudget, TestGuidesFeedFindingsFromStdin, TestReviewFixerTrailersInOneFinalParagraph, TestCriticRecordsNoFindingAndNoCommit, TestSharedPhrasesInExactlyOneFile and the nine per-topic Text tests) and no longer TestReviewText, which grep confirms is gone from internal/ and cmd/. Re-ran every Checks: command, all exit 0: rhiza slice --change 01M76XRETZFTD2TQ5GS3WYAERA --json; go test -trimpath ./internal/guide/...; go test -trimpath ./internal/workflow -run 'TestBuiltinRolePromptsScopeTests|TestBuiltin'; go test -trimpath ./cmd/rhiza -run TestClaudeMDScopedTests; rhiza ship generate; rhiza render (tree stayed clean apart from this session's actor log); make self-check. Re-ran the leaf acceptance: go test -trimpath ./internal/api -run GuideInstall ok; go test -trimpath ./internal/guide/... ok; go test -trimpath ./cmd/rhiza -run 'TestScripts/(guide|guide-install)' ok with both scripts PASS under -v. No code, stubs or TODOs in the diff.

## Links

### Out

- verifies [review-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-role-skills-installed.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
