# test-scope-block-names-every-role

decision · `01M79WRVM4D2403DAVSXE78CK4`

## reasoning

Finding f1 is accurate: the block, written for build role topics in #521, opens 'Builders and verifiers run only the tests...', yet this change set includes it in all nine review role topics (check-compliance, check-quality, check-qa, check-model, critic, fixer, ci-fix, gate-fix, conflict-fix), so a check, critic or fixer reads a rule that addresses other roles. guides-scoped-tests and role-prompts-scoped-tests, both updated in this change set, say each review role topic tells checks, critics and fixers to run only scoped tests through this block, so the model is right and the text is wrong. The fix is one sentence in one shared file plus one asserted phrase in shared_test.go; workflow_test.go's topicRule phrases ('only the tests covering what the change touches', 'never the repo's full suite') and review_test.go's ('scoped commands', 'never the repo's full suite') survive the rewording. Changing the leaf or the nodes to accept a builders-only sentence would contradict the change set's intent that review roles get the same rule. The fix was never committed across the 3 cycles (the file is unchanged since b635a1c8d), so a fixer should make this edit and keep the line breaks the tests match (or adjust those strings).

## statement

Fix the code, not the leaf: reword the shared test-scope block (internal/guide/text/shared/test-scope.md) so it is role-neutral or names every role that carries it (builders, verifiers, checks, critics and fixers), e.g. 'Run only the tests covering what the change touches, with the repo's scoped commands, never the repo's full suite: the PR's CI runs it before merge.' Keep the phrases the tests look for ('only the tests covering what the change touches', 'scoped commands', 'never the repo's full suite', 'the PR's CI runs it before merge'), update the opening phrase TestSharedPhrasesInExactlyOneFile asserts in internal/guide/shared_test.go to match, and name no Go, make or testscript command. Leaf review-role-skills-installed and the guides-scoped-tests and role-prompts-scoped-tests nodes stay as they are.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [leaf-failed-01m79dzc-1](../question/leaf-failed-01m79dzc-1.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
